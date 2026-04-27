package routes

import (
	"go-bin/internal/handler"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine, userHandler *handler.UserHandler, secretHandler *handler.SecretHandler) {
	api := router.Group("/api/v1")
	users := api.Group("/users")
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	secrets := api.Group("/secrets")
	secrets.POST("", secretHandler.Create)
	secrets.GET("/:token", secretHandler.GetByToken)
	secrets.POST("/:token/unlock", secretHandler.Unlock)
}
