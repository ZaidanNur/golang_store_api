package http

import (
	"net/http"
	"strconv"
	"test-elabram/internal/domain"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, userUsecase domain.UserUsecase) {
	handler := &userHandler{
		userUsecase: userUsecase,
	}

	r.GET("/users", handler.GetAllUsers)
	r.GET("/users/:id", handler.GetUserByID)
	r.POST("/users", handler.CreateUser)
}

func (h *userHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	user, err := h.userUsecase.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.userUsecase.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
