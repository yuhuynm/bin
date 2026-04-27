package handler

import (
	"errors"
	"net/http"

	"go-bin/internal/dto"
	"go-bin/internal/service"

	"github.com/gin-gonic/gin"
)

type SecretHandler struct {
	service service.SecretService
}

func NewSecretHandler(service service.SecretService) *SecretHandler {
	return &SecretHandler{service: service}
}

func (h *SecretHandler) Create(c *gin.Context) {
	var input dto.CreateSecretRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secret, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ToCreateSecretResponse(*secret))
}

func (h *SecretHandler) GetByToken(c *gin.Context) {
	token := c.Param("token")

	secret, content, err := h.service.GetByToken(token)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSecretNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrSecretUnavailable):
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if secret.HasPassword {
		c.JSON(http.StatusOK, dto.ToSecretMetaResponse(*secret))
		return
	}

	c.JSON(http.StatusOK, dto.ToSecretContentResponse(*secret, content))
}

func (h *SecretHandler) Unlock(c *gin.Context) {
	token := c.Param("token")

	var input dto.UnlockSecretRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secret, content, err := h.service.Unlock(token, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSecretNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrSecretUnavailable):
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrSecretInvalidPassword), errors.Is(err, service.ErrSecretPasswordRequired):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToSecretContentResponse(*secret, content))
}
