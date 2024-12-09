package models

import "time"

type Product struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"not null" json:"user_id"`
	ProductName        string         `gorm:"size:255;not null" json:"product_name"`
	ProductDescription string         `gorm:"type:text" json:"product_description"`
	ProductPrice       float64        `json:"product_price"`
	ProductImages      string         `gorm:"type:text" json:"product_images"`   // Comma-separated image URLs
	CompressedImages   string         `gorm:"type:text" json:"compressed_images"` // Comma-separated compressed image URLs
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
}