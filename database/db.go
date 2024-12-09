package database

import (
	"log"
	"product-management/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to the database
func Connect() {
	dsn := "host=localhost user=postgres password=productmanage dbname=product_management port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Automatically migrate models to the database
	err = DB.AutoMigrate(&models.Product{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}

	log.Println("Database connection successful and models migrated")
}