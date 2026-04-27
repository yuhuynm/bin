package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"time"

	"go-bin/internal/dto"
	"go-bin/internal/entity"
	"go-bin/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrSecretNotFound         = errors.New("secret not found")
	ErrSecretUnavailable      = errors.New("secret unavailable")
	ErrSecretPasswordRequired = errors.New("password required")
	ErrSecretInvalidPassword  = errors.New("invalid password")
	ErrSecretEncryptionFailed = errors.New("Secret encryption failed")
)

type SecretService interface {
	Create(input dto.CreateSecretRequest) (*entity.Secret, error)
	GetByToken(token string) (*entity.Secret, string, error)
	Unlock(token, password string) (*entity.Secret, string, error)
}

type secretService struct {
	repo      repository.SecretRepository
	secretKey []byte
}

func NewSecretService(repo repository.SecretRepository, secretKey string) SecretService {
	key := normalizeSecretKey(secretKey)
	return &secretService{repo: repo, secretKey: key}
}

func (s *secretService) Create(input dto.CreateSecretRequest) (*entity.Secret, error) {
	contentType := input.ContentType
	if contentType == "" {
		contentType = "text"
	}

	encryptedContent, err := s.encrypt(input.Content)
	if err != nil {
		return nil, err
	}

	secret := &entity.Secret{
		Token:            generateToken(32),
		EncryptedContent: encryptedContent,
		ContentType:      contentType,
		HasPassword:      strings.TrimSpace(input.Password) != "",
		MaxViews:         0,
	}

	if secret.HasPassword {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		hashValue := string(passwordHash)
		secret.PasswordHash = &hashValue
	}

	if input.ExpiresInHours > 0 {
		expiresAt := time.Now().Add(time.Duration(input.ExpiresInHours) * time.Hour)
		secret.ExpiresAt = &expiresAt
	}

	if input.OneTime {
		secret.MaxViews = 1
	}

	if err := s.repo.Create(secret); err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *secretService) GetByToken(token string) (*entity.Secret, string, error) {
	secret, err := s.repo.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrSecretNotFound
		}
		return nil, "", err
	}

	if err := validateSecretAvailability(secret); err != nil {
		return nil, "", err
	}

	if secret.HasPassword {
		return secret, "", nil
	}

	content, err := s.decrypt(secret.EncryptedContent)
	if err != nil {
		return nil, "", err
	}

	if err := s.markViewed(secret); err != nil {
		return nil, "", err
	}

	return secret, content, nil
}

func (s *secretService) Unlock(token, password string) (*entity.Secret, string, error) {
	secret, err := s.repo.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrSecretNotFound
		}
		return nil, "", err
	}

	if err := validateSecretAvailability(secret); err != nil {
		return nil, "", err
	}

	if !secret.HasPassword || secret.PasswordHash == nil {
		return nil, "", ErrSecretPasswordRequired
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*secret.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrSecretInvalidPassword
	}

	content, err := s.decrypt(secret.EncryptedContent)
	if err != nil {
		return nil, "", err
	}

	if err := s.markViewed(secret); err != nil {
		return nil, "", err
	}

	return secret, content, nil
}

func (s *secretService) markViewed(secret *entity.Secret) error {
	secret.ViewCount++

	if secret.MaxViews > 0 && secret.ViewCount >= secret.MaxViews {
		secret.IsConsumed = true
		now := time.Now()
		secret.ConsumedAt = &now
	}

	return s.repo.Update(secret)
}

func validateSecretAvailability(secret *entity.Secret) error {

	if secret.IsConsumed {
		return ErrSecretUnavailable
	}

	if secret.ExpiresAt != nil && time.Now().After(*secret.ExpiresAt) {
		return ErrSecretUnavailable
	}

	if secret.MaxViews > 0 && secret.ViewCount >= secret.MaxViews {
		return ErrSecretUnavailable
	}

	return nil
}

func (s *secretService) encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (s *secretService) decrypt(cipherText string) (string, error) {
	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", errors.New("invalid encrypted content")
	}

	nonce, encrypted := decoded[:nonceSize], decoded[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func generateToken(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "fallbackToken1234567890abcdef"
	}

	for index, value := range bytes {
		bytes[index] = alphabet[int(value)%len(alphabet)]
	}

	return string(bytes)
}

func normalizeSecretKey(secretKey string) []byte {
	key := []byte(secretKey)
	switch {
	case len(key) >= 32:
		return key[:32]
	default:
		padded := make([]byte, 32)
		copy(padded, key)
		return padded
	}
}
