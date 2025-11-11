package handlers

import (
	"errors"
	"net/http"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CartHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *CartHandler) Register(r *gin.RouterGroup) {
	r.POST("/cart", h.CreateCart)
	r.POST("/cart/items", h.AddItem)
	r.GET("/cart", h.GetCart)
	r.DELETE("/cart/items/:id", h.RemoveItem)
}

func (h *CartHandler) CreateCart(c *gin.Context) {
	var cart models.Cart
	if err := h.db.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cart)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For simplicity, assume a single cart with ID 1
	item.CartID = 1

	// Check if product variant exists and has enough stock
	var variant models.ProductVariant
	if err := h.db.First(&variant, item.ProductVariantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if variant.Stock < item.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Check if the item already exists in the cart
	var existingItem models.CartItem
	err := h.db.Where("cart_id = ? AND product_variant_id = ?", item.CartID, item.ProductVariantID).First(&existingItem).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := h.db.Begin()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new cart item
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Update quantity of existing item
		existingItem.Quantity += item.Quantity
		if err := tx.Save(&existingItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		item = existingItem
	}

	// Decrement stock
	variant.Stock -= item.Quantity
	if err := tx.Save(&variant).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	var cart models.Cart
	// For simplicity, assume a single cart with ID 1
	if err := h.db.Preload("Items.ProductVariant").First(&cart, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	id := c.Param("id")
	result := h.db.Delete(&models.CartItem{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
