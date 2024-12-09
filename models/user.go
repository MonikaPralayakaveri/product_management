package models

import "gorm.io/gorm"

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:255;not null"`
	Email string `gorm:"size:255;unique;not null"`
	gorm.Model  // Embedding to automatically include fields like CreatedAt, UpdatedAt, etc.
}