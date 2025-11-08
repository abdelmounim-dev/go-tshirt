package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	custom_errors "github.com/abdelmounim-dev/go-tshirt/internal/errors"
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/abdelmounim-dev/go-tshirt/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/products", handler.Create)

	t.Run("should return 400 when validation fails", func(t *testing.T) {
		// Arrange
		product := &models.Product{Name: "", Price: 10}
		mockService.EXPECT().Create(gomock.Any()).Return(&custom_errors.ValidationError{Message: "product name cannot be empty"})

		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 201 when product is created successfully", func(t *testing.T) {
		// Arrange
		product := &models.Product{Name: "T-shirt", Price: 10}
		mockService.EXPECT().Create(gomock.Any()).Return(nil)

		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestProductHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/products/:id", handler.Update)

	t.Run("should return 400 when validation fails", func(t *testing.T) {
		// Arrange
		product := &models.Product{Name: "", Price: 10}
		mockService.EXPECT().Update(gomock.Any()).Return(&custom_errors.ValidationError{Message: "product name cannot be empty"})

		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 200 when product is updated successfully", func(t *testing.T) {
		// Arrange
		product := &models.Product{Name: "T-shirt", Price: 10}
		mockService.EXPECT().Update(gomock.Any()).Return(nil)

		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestProductHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/products", handler.GetAll)

	t.Run("should return 200 with products and variants", func(t *testing.T) {
		// Arrange
		expectedProducts := []models.Product{
			{
				ID:   1,
				Name: "T-shirt Black",
				Price: 25.00,
				Variants: []models.ProductVariant{
					{ID: 1, ProductID: 1, Color: "Black", Size: "M", Stock: 10},
				},
			},
		}
		mockService.EXPECT().GetAll().Return(expectedProducts, nil)

		req, _ := http.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		var products []models.Product
		err := json.Unmarshal(rec.Body.Bytes(), &products)
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, expectedProducts[0].Name, products[0].Name)
		assert.Len(t, products[0].Variants, 1)
		assert.Equal(t, expectedProducts[0].Variants[0].Color, products[0].Variants[0].Color)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetAll().Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestProductHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/products/:id", handler.GetByID)

	t.Run("should return 200 with product and variants", func(t *testing.T) {
		// Arrange
		expectedProduct := &models.Product{
			ID:   1,
			Name: "T-shirt Black",
			Price: 25.00,
			Variants: []models.ProductVariant{
				{ID: 1, ProductID: 1, Color: "Black", Size: "M", Stock: 10},
			},
		}
		mockService.EXPECT().GetByID(uint(1)).Return(expectedProduct, nil)

		req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusOK, rec.Code)
		var product models.Product
		err := json.Unmarshal(rec.Body.Bytes(), &product)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct.Name, product.Name)
		assert.Len(t, product.Variants, 1)
		assert.Equal(t, expectedProduct.Variants[0].Color, product.Variants[0].Color)
	})

	t.Run("should return 404 when product not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetByID(uint(1)).Return(nil, &custom_errors.NotFoundError{Message: "product not found"})

		req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetByID(uint(1)).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
		rec := httptest.NewRecorder()

		// Act
		router.ServeHTTP(rec, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}