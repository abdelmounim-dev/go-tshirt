package api

import (
	"github.com/abdelmounim-dev/go-tshirt/internal/api/handlers"
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Auto-migrate models
	db.AutoMigrate(&models.Product{}, &models.ProductVariant{}, &models.Cart{}, &models.CartItem{})

	// Setup routes
	api := r.Group("/api")
	{
		productHandler := handlers.NewProductHandler(db)
		productHandler.Register(api)

		cartHandler := handlers.NewCartHandler(db)
		cartHandler.Register(api)
	}

	return r
}
