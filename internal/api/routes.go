package api

import (
	"github.com/abdelmounim-dev/go-tshirt/internal/api/handlers"
	"github.com/abdelmounim-dev/go-tshirt/internal/config"
	"github.com/abdelmounim-dev/go-tshirt/internal/db"
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	database, err := db.NewSQLite(cfg.DBPath)
	if err != nil {
		panic(err)
	}

	// Auto-migrate models
	err = database.AutoMigrate(&models.Product{}, &models.ProductVariant{})
	if err != nil {
		panic(err)
	}

	productHandler := handlers.NewProductHandler(database)

	api := r.Group("/api")
	productHandler.Register(api)

	return r
}
