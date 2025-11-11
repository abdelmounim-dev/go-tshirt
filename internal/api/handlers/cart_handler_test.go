package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

	testCases := []struct {
		name             string
		setup            func(db *gorm.DB) (map[string]interface{}, uint)
		expectedCode     int
		expectedStock    uint
		expectedQuantity uint
	}{
		{
			name: "should add an item to the cart successfully",
			setup: func(db *gorm.DB) (map[string]interface{}, uint) {
				product := models.Product{Name: "T-shirt", Price: 20}
				db.Create(&product)
				variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
				db.Create(&variant)
				cart := models.Cart{}
				db.Create(&cart)
				return map[string]interface{}{
					"product_variant_id": variant.ID,
					"quantity":           1,
					"cart_id":            cart.ID,
				}, variant.ID
			},
			expectedCode:     http.StatusCreated,
			expectedStock:    9,
			expectedQuantity: 1,
		},
		{
			name: "should update the quantity of an existing item",
			setup: func(db *gorm.DB) (map[string]interface{}, uint) {
				product := models.Product{Name: "T-shirt", Price: 20}
				db.Create(&product)
				variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
				db.Create(&variant)
				cart := models.Cart{}
				db.Create(&cart)
				db.Create(&models.CartItem{CartID: cart.ID, ProductVariantID: variant.ID, Quantity: 1})
				return map[string]interface{}{
					"product_variant_id": variant.ID,
					"quantity":           2,
					"cart_id":            cart.ID,
				}, variant.ID
			},
			expectedCode:     http.StatusCreated,
			expectedStock:    7,
			expectedQuantity: 3,
		},
		{
			name: "should return an error for insufficient stock",
			setup: func(db *gorm.DB) (map[string]interface{}, uint) {
				product := models.Product{Name: "T-shirt", Price: 20}
				db.Create(&product)
				variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 5}
				db.Create(&variant)
				cart := models.Cart{}
				db.Create(&cart)
				return map[string]interface{}{
					"product_variant_id": variant.ID,
					"quantity":           10,
					"cart_id":            cart.ID,
				}, variant.ID
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "should return an error for non-existent variant",
			setup: func(db *gorm.DB) (map[string]interface{}, uint) {
				cart := models.Cart{}
				db.Create(&cart)
				return map[string]interface{}{
					"product_variant_id": 999,
					"quantity":           1,
					"cart_id":            cart.ID,
				}, 0
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			err := db.AutoMigrate(&models.Cart{}, &models.CartItem{}, &models.Product{}, &models.ProductVariant{})
			assert.NoError(t, err)

			item, variantID := tc.setup(db)

			handler := NewCartHandler(db)
			router := gin.Default()
			api := router.Group("/api")
			handler.Register(api)

			body, _ := json.Marshal(item)
			req, _ := http.NewRequest(http.MethodPost, "/api/cart/items", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusCreated {
				var createdItem models.CartItem
				json.Unmarshal(rec.Body.Bytes(), &createdItem)
				assert.Equal(t, item["product_variant_id"].(uint), createdItem.ProductVariantID)
				assert.Equal(t, tc.expectedQuantity, createdItem.Quantity)

				var updatedVariant models.ProductVariant
				db.First(&updatedVariant, variantID)
				assert.Equal(t, tc.expectedStock, updatedVariant.Stock)
			}
		})
	}
}

func TestCartHandler_GetCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should get the cart with items successfully", func(t *testing.T) {
		db := setupTestDB(t)
		err := db.AutoMigrate(&models.Cart{}, &models.CartItem{}, &models.Product{}, &models.ProductVariant{})
		assert.NoError(t, err)

		// Create a product and variant and add it to the cart
		product := models.Product{Name: "T-shirt", Price: 20}
		db.Create(&product)
		variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
		db.Create(&variant)
		cart := models.Cart{}
		db.Create(&cart)
		cartItem := models.CartItem{CartID: cart.ID, ProductVariantID: variant.ID, Quantity: 2}
		db.Create(&cartItem)

		handler := NewCartHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		req, _ := http.NewRequest(http.MethodGet, "/api/cart", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var fetchedCart models.Cart
		json.Unmarshal(rec.Body.Bytes(), &fetchedCart)
		assert.Equal(t, cart.ID, fetchedCart.ID)
		assert.Len(t, fetchedCart.Items, 1)
		assert.Equal(t, cartItem.ID, fetchedCart.Items[0].ID)
		assert.Equal(t, variant.Color, fetchedCart.Items[0].ProductVariant.Color)
	})
}

func TestCartHandler_RemoveItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should remove an item from the cart successfully", func(t *testing.T) {
		db := setupTestDB(t)
		err := db.AutoMigrate(&models.Cart{}, &models.CartItem{}, &models.Product{}, &models.ProductVariant{})
		assert.NoError(t, err)

		// Create a product and variant and add it to the cart
		product := models.Product{Name: "T-shirt", Price: 20}
		db.Create(&product)
		variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
		db.Create(&variant)
		cart := models.Cart{}
		db.Create(&cart)
		cartItem := models.CartItem{CartID: cart.ID, ProductVariantID: variant.ID, Quantity: 1}
		db.Create(&cartItem)

		handler := NewCartHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		req, _ := http.NewRequest(http.MethodDelete, "/api/cart/items/"+strconv.Itoa(int(cartItem.ID)), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify the item was deleted
		var deletedItem models.CartItem
		err = db.First(&deletedItem, cartItem.ID).Error
		assert.Error(t, err) // Should not find the item
	})
}