package handlers

import (
	"bytes"
	"encoding/json"
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