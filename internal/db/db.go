package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string, migrate interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	if err := db.AutoMigrate(&migrate); err != nil {
		log.Fatalf("Could not migrate db table: %v", err)
	}

	return db, nil
}
