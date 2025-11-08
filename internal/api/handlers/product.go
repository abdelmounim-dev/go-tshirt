package handlers

import (
	"errors"
	"net/http"
	"strconv"

	custom_errors "github.com/abdelmounim-dev/go-tshirt/internal/errors"
	"github.com/abdelmounim-dev/go-tshirt/internal/models"
	"github.com/abdelmounim-dev/go-tshirt/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc service.ProductServiceInterface
}

func NewProductHandler(svc service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Register(r *gin.RouterGroup) {
	r.GET("/", h.GetAll)
	r.GET("/:id", h.GetByID)
	r.POST("/", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(&p); err != nil {
		var validationErr *custom_errors.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = uint(id)
	if err := h.svc.Update(&p); err != nil {
		var validationErr *custom_errors.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
