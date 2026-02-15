package http

import (
	"errors"
	"net/http"
	"strconv"

	"test-elabram/internal/delivery/helper"
	"test-elabram/internal/domain"
	"test-elabram/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(r *gin.Engine, productUsecase domain.ProductUsecase) {
	handler := &ProductHandler{
		productUsecase: productUsecase,
	}

	r.GET("/products", handler.GetAllProducts)
	r.GET("/products/:id", handler.GetProductByID)
	r.POST("/products", handler.CreateProduct)
	r.PUT("/products/:id", handler.EditProduct)
	r.DELETE("/products/:id", handler.DeleteProduct)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productUsecase.GetAllProducts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get products success",
		"data":    products,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	product, err := h.productUsecase.GetProductByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get product success",
		"data":    product,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErrors := make(map[string]string)
			for _, fe := range ve {
				fieldErrors[fe.Field()] = helper.MsgForTag(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Validation failed",
				"errors":  fieldErrors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	product := domain.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		IsActive:      req.IsActive,
		CategoryID:    req.CategoryID,
	}

	if err := h.productUsecase.CreateProduct(c, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "product created successfully",
		"data":    product,
	})
}

func (h *ProductHandler) EditProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErrors := make(map[string]string)
			for _, fe := range ve {
				fieldErrors[fe.Field()] = helper.MsgForTag(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Validation failed",
				"errors":  fieldErrors,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
		})
		return
	}

	product, err := h.productUsecase.EditProduct(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product updated successfully",
		"data":    product,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.productUsecase.DeleteProduct(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product deleted successfully",
	})
}
