package app

import (
	"go-bin/internal/config"
	"go-bin/internal/handler"
	"go-bin/internal/repository"
	"go-bin/internal/routes"
	"go-bin/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config config.Config
	router *gin.Engine
}

func NewServer(cfg config.Config, db *gorm.DB) *Server {
	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	secretRepo := repository.NewSecretRepository(db)
	secretService := service.NewSecretService(secretRepo, cfg.SecretKey)
	secretHandler := handler.NewSecretHandler(secretService)

	routes.User(router, userHandler, secretHandler)

	return &Server{
		config: cfg,
		router: router,
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
