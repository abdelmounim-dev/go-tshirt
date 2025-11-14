# T-Shirt Shop API Backend

This document provides an overview of the T-Shirt Shop API backend, its architecture, how to run it, and detailed documentation of its endpoints. This project is built in Go, using the Gin web framework and GORM for database interactions with SQLite.

## 🎯 Project Goal

Build a **clean, maintainable, and well-tested REST API backend** in Go that manages products for an online T-shirt store. This API supports the product lifecycle (CRUD), product options (color, size, stock), recommendations, and shopping cart features.

## 🚀 Technologies Used

*   **Programming Language**: Go
*   **Web Framework**: Gin
*   **ORM**: GORM
*   **Database**: SQLite
*   **Validation**: `go-playground/validator`

## 🏗️ Architecture

The project follows a layered architecture:

*   **`cmd/server`**: The main application entry point.
*   **`internal/api`**: Defines the API routes and handlers.
*   **`internal/service`**: Contains the business logic (not extensively used yet, logic mostly in handlers for simplicity).
*   **`internal/repository`**: Implements the database operations (GORM handles much of this).
*   **`internal/models`**: Defines the data models (`Product`, `ProductVariant`, `Cart`, `CartItem`).
*   **`internal/config`**: Manages application configuration.
*   **`internal/db`**: Handles the database connection.

## ⚙️ Building and Running

### Build

```bash
go build ./cmd/server
```

### Run

```bash
go run ./cmd/server/main.go
```

The server will start on the address specified in the configuration (default: `:8080`).

### Testing

```bash
go test ./...
```

## 📝 API Endpoints Documentation

All endpoints are prefixed with `/api`.

### 📦 Product API

Manages the lifecycle of products and their variants.

#### 1. List all products

*   **Endpoint**: `GET /api/products`
*   **Description**: Retrieves a list of all products, including their variants.
*   **Response (200 OK)**:
    ```json
    [
      {
        "id": 1,
        "name": "Basic Tee",
        "description": "A comfortable and stylish basic tee.",
        "price": 25.00,
        "image_url": "http://example.com/basic-tee.jpg",
        "variants": [
          {
            "id": 101,
            "product_id": 1,
            "color": "Black",
            "size": "M",
            "stock": 10
          },
          {
            "id": 102,
            "product_id": 1,
            "color": "White",
            "size": "L",
            "stock": 5
          }
        ],
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-27T10:00:00Z"
      }
    ]
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to retrieve products"
    }
    ```

#### 2. Retrieve product by ID

*   **Endpoint**: `GET /api/products/:id`
*   **Description**: Retrieves a single product by its ID, including its variants.
*   **Path Parameters**:
    *   `id` (integer): The ID of the product.
*   **Response (200 OK)**:
    ```json
    {
      "id": 1,
      "name": "Basic Tee",
      "description": "A comfortable and stylish basic tee.",
      "price": 25.00,
      "image_url": "http://example.com/basic-tee.jpg",
      "variants": [
        {
          "id": 101,
          "product_id": 1,
          "color": "Black",
          "size": "M",
          "stock": 10
        }
      ],
      "created_at": "2023-10-27T10:00:00Z",
      "updated_at": "2023-10-27T10:00:00Z"
    }
    ```
*   **Error Response (404 Not Found)**:
    ```json
    {
      "error": "Product not found"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to retrieve product"
    }
    ```

#### 3. Create new product

*   **Endpoint**: `POST /api/products`
*   **Description**: Creates a new product with its associated variants.
*   **Request Body**:
    ```json
    {
      "name": "New T-Shirt",
      "description": "A brand new awesome t-shirt.",
      "price": 30.00,
      "image_url": "http://example.com/new-tshirt.jpg",
      "variants": [
        {
          "color": "Blue",
          "size": "S",
          "stock": 20
        },
        {
          "color": "Blue",
          "size": "M",
          "stock": 15
        }
      ]
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "id": 2,
      "name": "New T-Shirt",
      "description": "A brand new awesome t-shirt.",
      "price": 30.00,
      "image_url": "http://example.com/new-tshirt.jpg",
      "variants": [
        {
          "id": 201,
          "product_id": 2,
          "color": "Blue",
          "size": "S",
          "stock": 20
        },
        {
          "id": 202,
          "product_id": 2,
          "color": "Blue",
          "size": "M",
          "stock": 15
        }
      ],
      "created_at": "2023-10-27T10:05:00Z",
      "updated_at": "2023-10-27T10:05:00Z"
    }
    ```
*   **Error Response (400 Bad Request)**:
    ```json
    {
      "error": "Invalid request body"
    }
    ```
    or (if validation fails, using raw validator output for now):
    ```json
    {
      "error": "Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to create product"
    }
    ```

#### 4. Update product details

*   **Endpoint**: `PUT /api/products/:id`
*   **Description**: Updates an existing product and its variants.
*   **Path Parameters**:
    *   `id` (integer): The ID of the product to update.
*   **Request Body**:
    ```json
    {
      "name": "Updated T-Shirt Name",
      "description": "An updated description.",
      "price": 35.00,
      "image_url": "http://example.com/updated-tshirt.jpg",
      "variants": [
        {
          "id": 101,
          "color": "Black",
          "size": "M",
          "stock": 8
        }
      ]
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "id": 1,
      "name": "Updated T-Shirt Name",
      "description": "An updated description.",
      "price": 35.00,
      "image_url": "http://example.com/updated-tshirt.jpg",
      "variants": [
        {
          "id": 101,
          "product_id": 1,
          "color": "Black",
          "size": "M",
          "stock": 8
        }
      ],
      "created_at": "2023-10-27T10:00:00Z",
      "updated_at": "2023-10-27T10:10:00Z"
    }
    ```
*   **Error Response (400 Bad Request)**:
    ```json
    {
      "error": "Invalid request body"
    }
    ```
    or (if validation fails, using raw validator output for now):
    ```json
    {
      "error": "Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'gt' tag"
    }
    ```
*   **Error Response (404 Not Found)**:
    ```json
    {
      "error": "Product not found"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to update product"
    }
    ```

#### 5. Delete product

*   **Endpoint**: `DELETE /api/products/:id`
*   **Description**: Deletes a product by its ID.
*   **Path Parameters**:
    *   `id` (integer): The ID of the product to delete.
*   **Response (204 No Content)**: (No response body)
*   **Error Response (404 Not Found)**:
    ```json
    {
      "error": "Product not found"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to delete product"
    }
    ```

### 🛒 Cart API

Manages the shopping cart functionality.

#### 1. Create a new cart

*   **Endpoint**: `POST /api/cart`
*   **Description**: Creates a new empty shopping cart.
*   **Request Body**: (None)
*   **Response (201 Created)**:
    ```json
    {
      "id": 1,
      "created_at": "2023-10-27T10:15:00Z",
      "updated_at": "2023-10-27T10:15:00Z",
      "items": []
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to create cart"
    }
    ```

#### 2. Add item to cart

*   **Endpoint**: `POST /api/cart/:cart_id/items`
*   **Description**: Adds a product variant to the specified cart. If the variant is already in the cart, its quantity is updated. Stock is checked and decremented.
*   **Path Parameters**:
    *   `cart_id` (integer): The ID of the cart to add the item to.
*   **Request Body**:
    ```json
    {
      "product_variant_id": 101,
      "quantity": 1
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "id": 1,
      "cart_id": 1,
      "product_variant_id": 101,
      "quantity": 1
    }
    ```
    (If updating an existing item, the `id` and `quantity` will reflect the updated item.)
*   **Error Response (400 Bad Request)**:
    ```json
    {
      "error": "Invalid request body"
    }
    ```
    or
    ```json
    {
      "error": "Insufficient stock"
    }
    ```
*   **Error Response (404 Not Found)**:
    ```json
    {
      "error": "Product variant not found"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Database error message"
    }
    ```

#### 3. Get cart contents

*   **Endpoint**: `GET /api/cart/:cart_id`
*   **Description**: Retrieves the contents of the specified cart, including product variant details.
*   **Path Parameters**:
    *   `cart_id` (integer): The ID of the cart to retrieve.
*   **Response (200 OK)**:
    ```json
    {
      "id": 1,
      "created_at": "2023-10-27T10:15:00Z",
      "updated_at": "2023-10-27T10:18:00Z",
      "items": [
        {
          "id": 1,
          "cart_id": 1,
          "product_variant_id": 101,
          "product_variant": {
            "id": 101,
            "product_id": 1,
            "color": "Black",
            "size": "M",
            "stock": 9
          },
          "quantity": 1
        }
      ]
    }
    ```
*   **Error Response (404 Not Found)**:
    ```json
    {
      "error": "Cart not found"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to retrieve cart"
    }
    ```

#### 4. Remove item from cart

*   **Endpoint**: `DELETE /api/cart/:cart_id/items/:item_id`
*   **Description**: Removes a specific item from the specified cart by its `CartItem` ID.
*   **Path Parameters**:
    *   `cart_id` (integer): The ID of the cart from which to remove the item.
    *   `item_id` (integer): The ID of the cart item to remove.
*   **Response (204 No Content)**: (No response body)
*   **Error Response (404 Not Found)**:
    ```json
    {
      "error": "Cart item not found"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Failed to delete cart item"
    }
    ```

### ✨ Recommendations API

Provides product recommendations.

#### 1. Get recommendations by color

*   **Endpoint**: `GET /api/recommendations`
*   **Description**: Retrieves a list of products that have variants of a specified color.
*   **Query Parameters**:
    *   `color` (string, required): The color to filter recommendations by (e.g., `Black`, `White`).
*   **Response (200 OK)**:
    ```json
    [
      {
        "id": 1,
        "name": "Basic Tee",
        "description": "A comfortable and stylish basic tee.",
        "price": 25.00,
        "image_url": "http://example.com/basic-tee.jpg",
        "variants": [
          {
            "id": 101,
            "product_id": 1,
            "color": "Black",
            "size": "M",
            "stock": 9
          }
        ],
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-27T10:00:00Z"
      }
    ]
    ```
*   **Error Response (400 Bad Request)**:
    ```json
    {
      "error": "color query parameter is required"
    }
    ```
*   **Error Response (500 Internal Server Error)**:
    ```json
    {
      "error": "Database error message"
    }
    ```

## 📋 Data Models

### `Product`

Represents a T-shirt product.

```go
type Product struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description"`
	Price       float64          `json:"price" validate:"required,gt=0"`
	ImageURL    string           `json:"image_url"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID" validate:"dive"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
```

### `ProductVariant`

Represents a specific color and size combination for a product, with its stock.

```go
type ProductVariant struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ProductID uint   `json:"product_id"`
	Color     string `json:"color" validate:"required"`
	Size      string `json:"size" validate:"required"`
	Stock     uint   `json:"stock" validate:"required,gte=0"`
}
```

### `Cart`

Represents a shopping cart.

```go
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
}
```

### `CartItem`

Represents an item within a shopping cart, linked to a specific `ProductVariant`.

```go
type CartItem struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CartID           uint           `json:"cart_id"`
	ProductVariantID uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID" validate:"omitempty"`
	Quantity         uint           `json:"quantity" validate:"required,gte=1"`
}
```

## ⚠️ Current Status & Notes for Interviewers

*   **Phase 1 (Product CRUD)**: Fully implemented and tested.
*   **Phase 2 (Validation and Error Handling)**: Input validation is in place using `go-playground/validator`. However, the error responses for validation failures currently return raw validator messages. A future improvement would be to standardize these into more user-friendly formats.
*   **Phase 3 (Product Options and Availability)**: Fully implemented. Products support variants with color, size, and stock tracking.
*   **Phase 4 (Shopping Cart API)**:
    *   **Stock Check**: Implemented and tested. Items are only added if stock is available, and stock is decremented upon addition.
    *   **Multi-Cart Support**: The cart API now supports multiple carts, with cart IDs specified in the URL for adding, retrieving, and removing items.
*   **Phase 5 (Recommendations API)**: Implemented with a basic recommendation logic (by color).

This project demonstrates a solid foundation for a Go REST API, adhering to TDD principles and clean architecture, with clear next steps for further development.
