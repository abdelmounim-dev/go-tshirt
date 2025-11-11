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

	// Check if product exists
	var product models.Product
	if err := h.db.First(&product, item.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	// Implementation to follow
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	// Implementation to follow
}
