package db

import (
	"fmt"
	"log"
	"os"

	"github.com/AronditFire/todo-app/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=postgres port=5432 sslmode=disable", os.Getenv("DB_PASSWORD"))
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&entity.Task{}, &entity.User{}); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	return db, err
}

func CloseConnection(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}
