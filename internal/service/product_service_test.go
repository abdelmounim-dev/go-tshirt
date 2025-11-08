package service

import (

	"errors"

	"testing"



	custom_errors "github.com/abdelmounim-dev/go-tshirt/internal/errors"

	"github.com/abdelmounim-dev/go-tshirt/internal/models"

	"github.com/abdelmounim-dev/go-tshirt/internal/repository/mocks"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"

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

		}



		// Act

		err := productService.Create(product)



		// Assert

		assert.Error(t, err)

		var validationErr *custom_errors.ValidationError

		assert.True(t, errors.As(err, &validationErr))

		assert.Equal(t, "product price must be greater than zero", err.Error())

	})



	t.Run("should create product successfully", func(t *testing.T) {

		// Arrange

		product := &models.Product{

			Name:        "T-shirt",

			Description: "A nice t-shirt",

			Price:       10.0,

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

		}



		// Act

		err := productService.Update(product)



		// Assert

		assert.Error(t, err)

		var validationErr *custom_errors.ValidationError

		assert.True(t, errors.As(err, &validationErr))

		assert.Equal(t, "product price must be greater than zero", err.Error())

	})



	t.Run("should update product successfully", func(t *testing.T) {

		// Arrange

		product := &models.Product{

			ID:          1,

			Name:        "T-shirt",

			Description: "A nice t-shirt",

			Price:       10.0,

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
