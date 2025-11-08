package service

import (

	"errors"

	"testing"



	custom_errors "github.com/abdelmounim-dev/go-tshirt/internal/errors"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"

	"github.com/abdelmounim-dev/go-tshirt/internal/repository/mocks"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"

	"gorm.io/gorm"

)



func TestProductService_Create(t *testing.T) {



	ctrl := gomock.NewController(t)



	defer ctrl.Finish()







	mockRepo := mocks.NewMockProductRepository(ctrl)



	productService := NewProductService(mockRepo)







	t.Run("should return validation error when product name is empty", func(t *testing.T) {



		// Arrange



		product := &models.Product{



			Name:        "",



			Description: "A nice t-shirt",



			Price:       10.0,



			Variants: []models.ProductVariant{



				{Color: "Black", Size: "M", Stock: 10},



			},



		}







		// Act



		err := productService.Create(product)







		// Assert



		assert.Error(t, err)



		var validationErr *custom_errors.ValidationError



		assert.True(t, errors.As(err, &validationErr))



		assert.Equal(t, "product name cannot be empty", err.Error())



	})







	t.Run("should return validation error when product price is zero or less", func(t *testing.T) {



		// Arrange



		product := &models.Product{



			Name:        "T-shirt",



			Description: "A nice t-shirt",



			Price:       0,



			Variants: []models.ProductVariant{



				{Color: "Black", Size: "M", Stock: 10},



			},



		}







		// Act



		err := productService.Create(product)







		// Assert



		assert.Error(t, err)



		var validationErr *custom_errors.ValidationError



		assert.True(t, errors.As(err, &validationErr))



		assert.Equal(t, "product price must be greater than zero", err.Error())



	})







	t.Run("should return validation error when product has no variants", func(t *testing.T) {



		// Arrange



		product := &models.Product{



			Name:        "T-shirt",



			Description: "A nice t-shirt",



			Price:       10.0,



			Variants:    []models.ProductVariant{},



		}







		// Act



		err := productService.Create(product)







		// Assert



		assert.Error(t, err)



		var validationErr *custom_errors.ValidationError



		assert.True(t, errors.As(err, &validationErr))



		assert.Equal(t, "product must have at least one variant", err.Error())



	})







	t.Run("should create product successfully", func(t *testing.T) {



		// Arrange



		product := &models.Product{



			Name:        "T-shirt",



			Description: "A nice t-shirt",



			Price:       10.0,



			Variants: []models.ProductVariant{



				{Color: "Black", Size: "M", Stock: 10},



			},



		}



		mockRepo.EXPECT().Create(product).Return(nil)







		// Act



		err := productService.Create(product)







		// Assert



		assert.NoError(t, err)



	})







	t.Run("should return error when repository fails", func(t *testing.T) {



		// Arrange



		product := &models.Product{



			Name:        "T-shirt",



			Description: "A nice t-shirt",



			Price:       10.0,



			Variants: []models.ProductVariant{



				{Color: "Black", Size: "M", Stock: 10},



			},



		}



		expectedErr := errors.New("database error")



		mockRepo.EXPECT().Create(product).Return(expectedErr)







		// Act



		err := productService.Create(product)







		// Assert



		assert.Error(t, err)



		assert.Equal(t, expectedErr, err)



	})



}



func TestProductService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockRepo)

	t.Run("should return validation error when product name is empty", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			ID:          1,
			Name:        "",
			Description: "A nice t-shirt",
			Price:       10.0,
			Variants: []models.ProductVariant{
				{Color: "Black", Size: "M", Stock: 10},
			},
		}

		// Act
		err := productService.Update(product)

		// Assert
		assert.Error(t, err)
		var validationErr *custom_errors.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "product name cannot be empty", err.Error())
	})

	t.Run("should return validation error when product price is zero or less", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			ID:          1,
			Name:        "T-shirt",
			Description: "A nice t-shirt",
			Price:       0,
			Variants: []models.ProductVariant{
				{Color: "Black", Size: "M", Stock: 10},
			},
		}

		// Act
		err := productService.Update(product)

		// Assert
		assert.Error(t, err)
		var validationErr *custom_errors.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "product price must be greater than zero", err.Error())
	})

	t.Run("should return validation error when product has no variants", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			ID:          1,
			Name:        "T-shirt",
			Description: "A nice t-shirt",
			Price:       10.0,
			Variants:    []models.ProductVariant{},
		}

		// Act
		err := productService.Update(product)

		// Assert
		assert.Error(t, err)
		var validationErr *custom_errors.ValidationError
		assert.True(t, errors.As(err, &validationErr))
		assert.Equal(t, "product must have at least one variant", err.Error())
	})

	t.Run("should update product successfully", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			ID:          1,
			Name:        "T-shirt",
			Description: "A nice t-shirt",
			Price:       10.0,
			Variants: []models.ProductVariant{
				{Color: "Black", Size: "M", Stock: 10},
			},
		}
		mockRepo.EXPECT().Update(product).Return(nil)

		// Act
		err := productService.Update(product)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			ID:          1,
			Name:        "T-shirt",
			Description: "A nice t-shirt",
			Price:       10.0,
			Variants: []models.ProductVariant{
				{Color: "Black", Size: "M", Stock: 10},
			},
		}
		expectedErr := errors.New("database error")
		mockRepo.EXPECT().Update(product).Return(expectedErr)

		// Act
		err := productService.Update(product)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestProductService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockRepo)

	t.Run("should return product successfully", func(t *testing.T) {
		// Arrange
		expectedProduct := &models.Product{
			ID:   1,
			Name: "T-shirt",
			Price: 10.0,
			Variants: []models.ProductVariant{
				{Color: "Black", Size: "M", Stock: 10},
			},
		}
		mockRepo.EXPECT().GetByID(uint(1)).Return(expectedProduct, nil)

		// Act
		product, err := productService.GetByID(uint(1))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("should return NotFoundError when product is not found", func(t *testing.T) {
		// Arrange
		mockRepo.EXPECT().GetByID(uint(1)).Return(nil, gorm.ErrRecordNotFound)

		// Act
		product, err := productService.GetByID(uint(1))

		// Assert
		assert.Error(t, err)
		var notFoundErr *custom_errors.NotFoundError
		assert.True(t, errors.As(err, &notFoundErr))
		assert.Nil(t, product)
	})

	t.Run("should return generic error when repository fails", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database error")
		mockRepo.EXPECT().GetByID(uint(1)).Return(nil, expectedErr)

		// Act
		product, err := productService.GetByID(uint(1))

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, product)
	})
}
