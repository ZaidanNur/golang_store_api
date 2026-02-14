package http

import (
	"net/http"
	"strconv"
	"test-elabram/internal/domain"

	"github.com/gin-gonic/gin"
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
	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.categoryUsecase.CreateCategory(c.Request.Context(), &category); err != nil {
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
	var category domain.Category
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.categoryUsecase.EditCategory(c.Request.Context(), id, &category); err != nil {
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
