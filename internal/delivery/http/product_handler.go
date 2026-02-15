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

	r.GET("/products/report", handler.GetProductReport)
	r.GET("/products", handler.GetAllProducts)
	r.GET("/products/:id", handler.GetProductByID)
	r.POST("/products", handler.CreateProduct)
	r.PUT("/products/:id", handler.EditProduct)
	r.DELETE("/products/:id", handler.DeleteProduct)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid pagination params",
		})
		return
	}

	var filters dto.ProductFilterParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid filter params",
		})
		return
	}

	result, err := h.productUsecase.GetAllProductsPaginated(c, filters, pq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      http.StatusOK,
		"message":     "get products success",
		"data":        result.Data,
		"page":        result.Page,
		"limit":       result.Limit,
		"total_items": result.TotalItems,
		"total_pages": result.TotalPages,
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

// GetProductReport returns a dashboard-style report of all products.
func (h *ProductHandler) GetProductReport(c *gin.Context) {
	report, err := h.productUsecase.GetProductReport(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get product report success",
		"data":    report,
	})
}
