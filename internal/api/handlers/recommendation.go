package handlers

import (
	"net/http"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecommendationHandler struct {
	db *gorm.DB
}

func NewRecommendationHandler(db *gorm.DB) *RecommendationHandler {
	return &RecommendationHandler{
		db: db,
	}
}

func (h *RecommendationHandler) Register(r *gin.RouterGroup) {
	r.GET("/recommendations", h.GetRecommendations)
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	color := c.Query("color")
	if color == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "color query parameter is required"})
		return
	}

	var products []models.Product
	if err := h.db.Joins("JOIN product_variants ON product_variants.product_id = products.id").Where("product_variants.color = ?", color).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
