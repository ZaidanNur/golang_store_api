package dto

type CreateProductRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Price         int    `json:"price" binding:"required,gt=0"`
	StockQuantity int    `json:"stock_quantity" binding:"required,gte=0"`
	IsActive      bool   `json:"is_active"`
	CategoryID    uint   `json:"category_id" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Price         *int    `json:"price" binding:"omitempty,gt=0"`
	StockQuantity *int    `json:"stock_quantity" binding:"omitempty,gte=0"`
	IsActive      *bool   `json:"is_active"`
	CategoryID    *uint   `json:"category_id" binding:"omitempty,gt=0"`
}
