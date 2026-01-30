package gorm

import (
	"log"

	"github.com/RodrigoGonzalez78/internal/infrastructure/persistence/gorm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Base de datos conectada correctamente")
	return &Database{DB: db}, nil
}

func (d *Database) Migrate() error {
	return d.DB.AutoMigrate(&models.UserModel{}, &models.ImageModel{})
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
