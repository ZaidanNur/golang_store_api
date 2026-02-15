package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}
