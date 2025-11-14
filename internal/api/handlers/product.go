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

type ProductHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *ProductHandler) Register(r *gin.RouterGroup) {
	// Product routes
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("", h.GetAll)
		productRoutes.GET("/:id", h.GetByID)
		productRoutes.POST("", h.Create)
		productRoutes.PUT("/:id", h.Update)
		productRoutes.DELETE("/:id", h.Delete)

		// Product Variant routes
		variantRoutes := productRoutes.Group("/:id/variants")
		{
			variantRoutes.GET("", h.GetAllVariants)
			variantRoutes.GET("/:variant_id", h.GetVariantByID)
			variantRoutes.POST("", h.CreateVariant)
			variantRoutes.PUT("/:variant_id", h.UpdateVariant)
			variantRoutes.DELETE("/:variant_id", h.DeleteVariant)
		}
	}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	var products []models.Product
	if err := h.db.Preload("Variants").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.db.Preload("Variants").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingProduct models.Product
	if err := h.db.First(&existingProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.ID = existingProduct.ID // Ensure the ID from the URL is used
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	result := h.db.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetAllVariants(c *gin.Context) {
	productId := c.Param("id")
	var variants []models.ProductVariant
	if err := h.db.Where("product_id = ?", productId).Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variants)
}

func (h *ProductHandler) GetVariantByID(c *gin.Context) {
	id := c.Param("variant_id")
	productId := c.Param("id")
	var variant models.ProductVariant
	if err := h.db.Where("id = ? AND product_id = ?", id, productId).First(&variant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) CreateVariant(c *gin.Context) {
	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var variant models.ProductVariant
	if err := c.ShouldBindJSON(&variant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant.ProductID = uint(productId)

	if err := h.validate.Struct(variant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, variant)
}

func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	id := c.Param("variant_id")
	productId := c.Param("id")
	var variant models.ProductVariant
	if err := c.ShouldBindJSON(&variant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(variant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingVariant models.ProductVariant
	if err := h.db.Where("id = ? AND product_id = ?", id, productId).First(&existingVariant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	variant.ID = existingVariant.ID // Ensure the ID from the URL is used
	variant.ProductID = existingVariant.ProductID
	if err := h.db.Save(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variant)
}

func (h *ProductHandler) DeleteVariant(c *gin.Context) {
	id := c.Param("variant_id")
	productId := c.Param("id")
	result := h.db.Where("id = ? AND product_id = ?", id, productId).Delete(&models.ProductVariant{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
