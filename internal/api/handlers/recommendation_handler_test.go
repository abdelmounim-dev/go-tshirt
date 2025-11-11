package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationHandler_GetRecommendations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return recommended products based on color", func(t *testing.T) {
		db := setupTestDB(t)
		err := db.AutoMigrate(&models.Product{}, &models.ProductVariant{})
		assert.NoError(t, err)

		// Create products
		p1 := models.Product{Name: "T-shirt", Price: 20, Variants: []models.ProductVariant{{Color: "Black", Size: "M", Stock: 10}}}
		p2 := models.Product{Name: "Polo", Price: 30, Variants: []models.ProductVariant{{Color: "Black", Size: "L", Stock: 5}}}
		p3 := models.Product{Name: "Hoodie", Price: 50, Variants: []models.ProductVariant{{Color: "White", Size: "M", Stock: 15}}}
		db.Create(&p1)
		db.Create(&p2)
		db.Create(&p3)

		handler := NewRecommendationHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		req, _ := http.NewRequest(http.MethodGet, "/api/recommendations?color=Black", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var recommendations []models.Product
		json.Unmarshal(rec.Body.Bytes(), &recommendations)
		assert.Len(t, recommendations, 2)
		assert.Equal(t, p1.Name, recommendations[0].Name)
		assert.Equal(t, p2.Name, recommendations[1].Name)
	})
}
