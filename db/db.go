package db

import (
	shared "auth-ms/shared/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Database interface {
	Init() *gorm.DB
}

type DatabaseImpl struct {
	instance *gorm.DB
}

var Instance *gorm.DB

func NewDatabase() Database {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_PATH")))
	if err != nil {
		shared.Log().Fatal("Database error: %v", err)
	}
	return &DatabaseImpl{instance: db}
}

func (d *DatabaseImpl) Init() *gorm.DB {
	Instance = d.instance
	return d.instance
}
