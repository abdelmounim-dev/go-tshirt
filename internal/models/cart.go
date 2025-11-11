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
	ID         uint    `json:"id" gorm:"primaryKey"`
	CartID     uint    `json:"cart_id"`
	ProductID  uint    `json:"product_id"`
	Product    Product `json:"product,omitempty" gorm:"foreignKey:ProductID" validate:"omitempty"`
	Quantity   uint    `json:"quantity" validate:"required,gte=1"`
}
