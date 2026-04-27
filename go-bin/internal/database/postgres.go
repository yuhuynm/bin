package database

import (
	"go-bin/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.Secret{}); err != nil {
		return nil, err
	}

	return db, nil
}
