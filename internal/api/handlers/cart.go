package handlers

import (
	"errors"
	"net/http"
	"strconv"

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
	cartRoutes := r.Group("/cart")
	{
		cartRoutes.POST("", h.CreateCart)
		cartRoutes.GET("/:cart_id", h.GetCart)
		cartRoutes.POST("/:cart_id/items", h.AddItem)
		cartRoutes.DELETE("/:cart_id/items/:item_id", h.RemoveItem)
	}
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
	cartID := c.Param("cart_id")

	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.CartID = 0 // Reset to avoid any passed in value
	if _, err := strconv.Atoi(cartID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}
	
	cid, _ := strconv.Atoi(cartID)
	item.CartID = uint(cid)


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
	cartID := c.Param("cart_id")

	var cart models.Cart
	if err := h.db.Preload("Items.ProductVariant").First(&cart, cartID).Error; err != nil {
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
	cartID := c.Param("cart_id")
	itemID := c.Param("item_id")

	var cartItem models.CartItem
	if err := h.db.Where("cart_id = ? AND id = ?", cartID, itemID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Delete(&cartItem)
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
