package dto

import (
	"time"

	"go-bin/internal/entity"
)

type CreateSecretRequest struct {
	Content        string `json:"content" binding:"required,max=50000"`
	ContentType    string `json:"contentType" binding:"omitempty,oneof=text env markdown"`
	Password       string `json:"password" binding:"omitempty,min=6,max=128"`
	ExpiresInHours int    `json:"expiresInHours" binding:"omitempty,min=0,max=720"`
	OneTime        bool   `json:"oneTime"`
}

type CreateSecretResponse struct {
	Token       string     `json:"token"`
	URL         string     `json:"url"`
	HasPassword bool       `json:"hasPassword"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type SecretMetaResponse struct {
	Token        string     `json:"token"`
	ContentType  string     `json:"contentType"`
	HasPassword  bool       `json:"hasPassword"`
	RequiresAuth bool       `json:"requiresAuth"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	IsConsumed   bool       `json:"isConsumed"`
}

type SecretContentResponse struct {
	Token       string     `json:"token"`
	Content     string     `json:"content"`
	ContentType string     `json:"contentType"`
	HasPassword bool       `json:"hasPassword"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type UnlockSecretRequest struct {
	Password string `json:"password"`
}

func ToCreateSecretResponse(secret entity.Secret) CreateSecretResponse {
	return CreateSecretResponse{
		Token:       secret.Token,
		URL:         "/api/v1/secrets/" + secret.Token,
		HasPassword: secret.HasPassword,
		ExpiresAt:   secret.ExpiresAt,
	}
}

func ToSecretMetaResponse(secret entity.Secret) SecretMetaResponse {
	return SecretMetaResponse{
		Token:        secret.Token,
		ContentType:  secret.ContentType,
		HasPassword:  secret.HasPassword,
		RequiresAuth: secret.HasPassword,
		ExpiresAt:    secret.ExpiresAt,
		IsConsumed:   secret.IsConsumed,
	}
}

func ToSecretContentResponse(secret entity.Secret, content string) SecretContentResponse {
	return SecretContentResponse{
		Token:       secret.Token,
		Content:     content,
		ContentType: secret.ContentType,
		HasPassword: secret.HasPassword,
		ExpiresAt:   secret.ExpiresAt,
	}
}
