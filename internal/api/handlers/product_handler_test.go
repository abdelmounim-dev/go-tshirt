package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Product{}, &models.ProductVariant{})
	assert.NoError(t, err)
	return db
}

func TestProductHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		product      models.Product
		expectedCode int
		expectedBody string
	}{
		{
			name:         "should return 400 when name is empty",
			product:      models.Product{Name: "", Price: 10},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name:         "should return 400 when price is zero",
			product:      models.Product{Name: "T-shirt", Price: 0},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'required' tag"}`,
		},
		{
			name:         "should return 400 when variant color is empty",
			product:      models.Product{Name: "T-shirt", Price: 10, Variants: []models.ProductVariant{{Color: "", Size: "M", Stock: 10}}},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Variants[0].Color' Error:Field validation for 'Color' failed on the 'required' tag"}`,
		},
		{
			name:         "should return 400 when variant size is empty",
			product:      models.Product{Name: "T-shirt", Price: 10, Variants: []models.ProductVariant{{Color: "Black", Size: "", Stock: 10}}},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Variants[0].Size' Error:Field validation for 'Size' failed on the 'required' tag"}`,
		},
		{
			name:         "should return 201 when product is created successfully",
			product:      models.Product{Name: "T-shirt", Price: 10},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			handler := NewProductHandler(db)
			router := gin.Default()
			api := router.Group("/api")
			handler.Register(api)

			body, _ := json.Marshal(tc.product)
			req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusCreated {
				var createdProduct models.Product
				err := json.Unmarshal(rec.Body.Bytes(), &createdProduct)
				assert.NoError(t, err)
				assert.Equal(t, tc.product.Name, createdProduct.Name)
				assert.Equal(t, tc.product.Price, createdProduct.Price)
				assert.NotZero(t, createdProduct.ID)

				var dbProduct models.Product
				err = db.First(&dbProduct, createdProduct.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tc.product.Name, dbProduct.Name)
			} else {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestProductHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		productID    string
		product      models.Product
		setupDB      func(*gorm.DB) uint
		expectedCode int
		expectedBody string
	}{
		{
			name:      "should return 400 when name is empty",
			productID: "1",
			product:   models.Product{Name: "", Price: 10},
			setupDB: func(db *gorm.DB) uint {
				p := models.Product{Name: "Old Name", Price: 50}
				db.Create(&p)
				return p.ID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
		{
			name:      "should return 400 when price is zero",
			productID: "1",
			product:   models.Product{Name: "Updated T-shirt", Price: 0},
			setupDB: func(db *gorm.DB) uint {
				p := models.Product{Name: "Old Name", Price: 50}
				db.Create(&p)
				return p.ID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'required' tag"}`,
		},
		{
			name:      "should return 400 when variant color is empty",
			productID: "1",
			product:   models.Product{Name: "T-shirt", Price: 10, Variants: []models.ProductVariant{{Color: "", Size: "M", Stock: 10}}},
			setupDB: func(db *gorm.DB) uint {
				p := models.Product{Name: "Old Name", Price: 50}
				db.Create(&p)
				return p.ID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Variants[0].Color' Error:Field validation for 'Color' failed on the 'required' tag"}`,
		},
		{
			name:      "should return 400 when variant size is empty",
			productID: "1",
			product:   models.Product{Name: "T-shirt", Price: 10, Variants: []models.ProductVariant{{Color: "Black", Size: "", Stock: 10}}},
			setupDB: func(db *gorm.DB) uint {
				p := models.Product{Name: "Old Name", Price: 50}
				db.Create(&p)
				return p.ID
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Product.Variants[0].Size' Error:Field validation for 'Size' failed on the 'required' tag"}`,
		},
		{
			name:         "should return 404 when product not found",
			productID:    "999",
			product:      models.Product{Name: "Updated T-shirt", Price: 20},
			setupDB:      func(db *gorm.DB) uint { return 0 }, // No product in DB
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"Product not found"}`,
		},
		{
			name:      "should return 200 when product is updated successfully",
			productID: "1",
			product:   models.Product{Name: "Updated T-shirt", Price: 20},
			setupDB: func(db *gorm.DB) uint {
				p := models.Product{Name: "Old Name", Price: 50}
				db.Create(&p)
				return p.ID
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			productID := tc.setupDB(db)
			handler := NewProductHandler(db)
			router := gin.Default()
			api := router.Group("/api")
			handler.Register(api)

			idStr := tc.productID
			if productID != 0 {
				idStr = strconv.Itoa(int(productID))
			}

			body, _ := json.Marshal(tc.product)
			req, _ := http.NewRequest(http.MethodPut, "/api/products/"+idStr, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				var updatedProduct models.Product
				err := json.Unmarshal(rec.Body.Bytes(), &updatedProduct)
				assert.NoError(t, err)
				assert.Equal(t, tc.product.Name, updatedProduct.Name)
				assert.Equal(t, tc.product.Price, updatedProduct.Price)

				var dbProduct models.Product
				err = db.First(&dbProduct, idStr).Error
				assert.NoError(t, err)
				assert.Equal(t, tc.product.Name, dbProduct.Name)
			} else {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestProductHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 with products and variants", func(t *testing.T) {
		db := setupTestDB(t)
		handler := NewProductHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		product1 := models.Product{Name: "T-shirt Black", Price: 25.00}
		db.Create(&product1)
		db.Create(&models.ProductVariant{ProductID: product1.ID, Color: "Black", Size: "M", Stock: 10})

		req, _ := http.NewRequest(http.MethodGet, "/api/products", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var products []models.Product
		err := json.Unmarshal(rec.Body.Bytes(), &products)
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, product1.Name, products[0].Name)
		assert.Len(t, products[0].Variants, 1)
		assert.Equal(t, "Black", products[0].Variants[0].Color)
	})

	t.Run("should return empty array if no products", func(t *testing.T) {
		db := setupTestDB(t)
		handler := NewProductHandler(db)
		router := gin.Default()
		api := router.Group("/api")
		handler.Register(api)

		req, _ := http.NewRequest(http.MethodGet, "/api/products", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[]`, rec.Body.String())
	})
}

func TestProductHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		productID    string
		setupDB      func(*gorm.DB) uint // Function to set up initial DB state and return ID
		expectedCode int
		expectedBody string
	}{
		{
			name:      "should return 200 with product and variants",
			productID: "1",
			setupDB: func(db *gorm.DB) uint {
				product := models.Product{Name: "T-shirt Black", Price: 25.00}
				db.Create(&product)
				db.Create(&models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10})
				return product.ID
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should return 404 when product not found",
			productID:    "999",
			setupDB:      func(db *gorm.DB) uint { return 0 }, // No product in DB
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"Product not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			productID := tc.setupDB(db)
			handler := NewProductHandler(db)
			router := gin.Default()
			api := router.Group("/api")
			handler.Register(api)

			idStr := tc.productID
			if productID != 0 {
				idStr = strconv.Itoa(int(productID))
			}

			req, _ := http.NewRequest(http.MethodGet, "/api/products/"+idStr, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedCode == http.StatusOK {
				var product models.Product
				err := json.Unmarshal(rec.Body.Bytes(), &product)
				assert.NoError(t, err)
				assert.Equal(t, "T-shirt Black", product.Name)
				assert.Equal(t, 25.00, product.Price)
				assert.Len(t, product.Variants, 1)
			} else {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		productID    string
		setupDB      func(*gorm.DB) uint // Function to set up initial DB state
		expectedCode int
		expectedBody string
	}{
		{
			name:      "should return 204 when product is deleted successfully",
			productID: "1",
			setupDB: func(db *gorm.DB) uint {
				p := models.Product{Name: "Product to Delete", Price: 10}
				db.Create(&p)
				return p.ID
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "should return 404 when product not found for deletion",
			productID:    "999",
			setupDB:      func(db *gorm.DB) uint { return 0 }, // No product in DB
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"Product not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			productID := tc.setupDB(db)
			handler := NewProductHandler(db)
			router := gin.Default()
			api := router.Group("/api")
			handler.Register(api)

			idStr := tc.productID
			if productID != 0 {
				idStr = strconv.Itoa(int(productID))
			}

			req, _ := http.NewRequest(http.MethodDelete, "/api/products/"+idStr, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedCode != http.StatusNoContent {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}

			if tc.expectedCode == http.StatusNoContent {
				var product models.Product
				err := db.First(&product, idStr).Error
				assert.Error(t, err) // Should not find the product
				assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
			}
		})
	}
}

func TestProductHandler_GetAllVariants(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	handler := NewProductHandler(db)
	router := gin.Default()
	api := router.Group("/api")
	handler.Register(api)

	product := models.Product{Name: "T-shirt", Price: 20}
	db.Create(&product)
	variant1 := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
	variant2 := models.ProductVariant{ProductID: product.ID, Color: "White", Size: "L", Stock: 5}
	db.Create(&variant1)
	db.Create(&variant2)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+strconv.Itoa(int(product.ID))+"/variants", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var variants []models.ProductVariant
	json.Unmarshal(rec.Body.Bytes(), &variants)
	assert.Len(t, variants, 2)
}

func TestProductHandler_GetVariantByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	handler := NewProductHandler(db)
	router := gin.Default()
	api := router.Group("/api")
	handler.Register(api)

	product := models.Product{Name: "T-shirt", Price: 20}
	db.Create(&product)
	variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
	db.Create(&variant)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+strconv.Itoa(int(product.ID))+"/variants/"+strconv.Itoa(int(variant.ID)), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var fetchedVariant models.ProductVariant
	json.Unmarshal(rec.Body.Bytes(), &fetchedVariant)
	assert.Equal(t, variant.ID, fetchedVariant.ID)
	assert.Equal(t, variant.Color, fetchedVariant.Color)
}

func TestProductHandler_CreateVariant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	handler := NewProductHandler(db)
	router := gin.Default()
	api := router.Group("/api")
	handler.Register(api)

	product := models.Product{Name: "T-shirt", Price: 20}
	db.Create(&product)

	variantData := gin.H{
		"color": "Blue",
		"size":  "S",
		"stock": 20,
	}
	body, _ := json.Marshal(variantData)
	req, _ := http.NewRequest(http.MethodPost, "/api/products/"+strconv.Itoa(int(product.ID))+"/variants", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var createdVariant models.ProductVariant
	json.Unmarshal(rec.Body.Bytes(), &createdVariant)
	assert.Equal(t, uint(product.ID), createdVariant.ProductID)
	assert.Equal(t, "Blue", createdVariant.Color)
}

func TestProductHandler_UpdateVariant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	handler := NewProductHandler(db)
	router := gin.Default()
	api := router.Group("/api")
	handler.Register(api)

	product := models.Product{Name: "T-shirt", Price: 20}
	db.Create(&product)
	variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
	db.Create(&variant)

	updateData := gin.H{
		"color": "Red",
		"size":  "M",
		"stock": 5,
	}
	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+strconv.Itoa(int(product.ID))+"/variants/"+strconv.Itoa(int(variant.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var updatedVariant models.ProductVariant
	db.First(&updatedVariant, variant.ID)
	assert.Equal(t, "Red", updatedVariant.Color)
	assert.Equal(t, uint(5), updatedVariant.Stock)
}

func TestProductHandler_DeleteVariant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	handler := NewProductHandler(db)
	router := gin.Default()
	api := router.Group("/api")
	handler.Register(api)

	product := models.Product{Name: "T-shirt", Price: 20}
	db.Create(&product)
	variant := models.ProductVariant{ProductID: product.ID, Color: "Black", Size: "M", Stock: 10}
	db.Create(&variant)

	req, _ := http.NewRequest(http.MethodDelete, "/api/products/"+strconv.Itoa(int(product.ID))+"/variants/"+strconv.Itoa(int(variant.ID)), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	var deletedVariant models.ProductVariant
	err := db.First(&deletedVariant, variant.ID).Error
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}