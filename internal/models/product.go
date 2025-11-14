package models

import "time"

type Product struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description"`
	Price       float64          `json:"price" validate:"required,gt=0"`
	ImageURL    string           `json:"image_url"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID" validate:"dive"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ProductVariant struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ProductID uint   `json:"product_id"`
	Color     string `json:"color" validate:"required"`
	Size      string `json:"size" validate:"required"`
	Stock     uint   `json:"stock" validate:"required,gte=0"`
}
