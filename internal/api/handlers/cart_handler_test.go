package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCartHandler_CreateCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create a new cart successfully", func(t *testing.T) {
		db := setupTestDB(t)
		err := db.AutoMigrate(&models.Cart{}, &models.CartItem{})
		assert.NoError(t, err)

		handler := NewCartHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		req, _ := http.NewRequest(http.MethodPost, "/api/cart", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestCartHandler_AddItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should add an item to the cart successfully", func(t *testing.T) {
		db := setupTestDB(t)
		err := db.AutoMigrate(&models.Cart{}, &models.CartItem{}, &models.Product{})
		assert.NoError(t, err)

		// Create a product to add to the cart
		product := models.Product{Name: "T-shirt", Price: 20}
		db.Create(&product)

		// Create a cart
		cart := models.Cart{}
		db.Create(&cart)

		handler := NewCartHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		item := map[string]interface{}{
			"product_id": product.ID,
			"quantity":   1,
			"cart_id":    cart.ID,
		}
		body, _ := json.Marshal(item)
		req, _ := http.NewRequest(http.MethodPost, "/api/cart/items", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdItem models.CartItem
		json.Unmarshal(rec.Body.Bytes(), &createdItem)
		assert.Equal(t, item["product_id"].(uint), createdItem.ProductID)
		assert.Equal(t, uint(item["quantity"].(int)), createdItem.Quantity)
	})
}