package models

import "time"

// Cart represents a shopping cart
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CartID           uint           `json:"cart_id"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID" validate:"omitempty"`
	Quantity         uint           `json:"quantity" validate:"required,gte=1"`
}
