package repository

import (
	"go-bin/internal/entity"

	"gorm.io/gorm"
)

type SecretRepository interface {
	Create(secret *entity.Secret) error
	FindByToken(token string) (*entity.Secret, error)
	Update(secret *entity.Secret) error
}

type secretRepository struct {
	db *gorm.DB
}

func NewSecretRepository(db *gorm.DB) SecretRepository {
	return &secretRepository{db: db}
}

func (r *secretRepository) Create(secret *entity.Secret) error {
	return r.db.Create(secret).Error
}

func (r *secretRepository) FindByToken(token string) (*entity.Secret, error) {
	var secret entity.Secret
	err := r.db.Where("token = ?", token).First(&secret).Error
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func (r *secretRepository) Update(secret *entity.Secret) error {
	return r.db.Save(secret).Error
}
