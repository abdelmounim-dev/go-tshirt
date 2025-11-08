package models

import "time"

type Product struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	ImageURL    string           `json:"image_url"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ProductVariant struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	ProductID uint    `json:"product_id"`
	Color     string  `json:"color"`
	Size      string  `json:"size"`
	Stock     uint    `json:"stock"`
}
