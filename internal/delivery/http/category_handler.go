package http

import (
	"net/http"
	"strconv"
	"test-elabram/internal/domain"
	"test-elabram/internal/dto"

	"errors"
	"test-elabram/internal/delivery/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type categoryHandler struct {
	categoryUsecase domain.CategoryUsecase
}

func NewCategoryHandler(r *gin.Engine, categoryUsecase domain.CategoryUsecase) {
	handler := &categoryHandler{
		categoryUsecase: categoryUsecase,
	}

	r.GET("/category", handler.GetAllCategories)
	r.GET("/category/:id", handler.GetCategoryByID)
	r.POST("/category", handler.CreateCategory)
	r.PUT("/category/:id", handler.EditCategory)
	r.DELETE("/category/:id", handler.DeleteCategory)
}

func (h *categoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryUsecase.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "get categories failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get categories success",
		"data":    categories,
	})
}

func (h *categoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	category, err := h.categoryUsecase.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get category success",
		"data":    category,
	})
}

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
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

	category := domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.categoryUsecase.CreateCategory(c, &category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "category created successfully",
		"data":    category,
	})
}

func (h *categoryHandler) EditCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req dto.UpdateCategoryRequest
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
	category, err := h.categoryUsecase.EditCategory(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "category updated successfully",
		"data":    category,
	})
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.categoryUsecase.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "category deleted successfully",
	})
}
