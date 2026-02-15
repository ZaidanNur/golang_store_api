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

type PaginationQuery struct {
	Page  int `form:"page,default=1" binding:"omitempty,min=1"`
	Limit int `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

type ProductFilterParams struct {
	Name       string `form:"name"`
	CategoryID *uint  `form:"category_id"`
	PriceMin   *int   `form:"price_min"`
	PriceMax   *int   `form:"price_max"`
	StockMin   *int   `form:"stock_min"`
	StockMax   *int   `form:"stock_max"`
	SortBy     string `form:"sort_by,default=created_at"`
	SortOrder  string `form:"sort_order,default=desc"`
}

type ProductReportResponse struct {
	TotalProducts int                 `json:"total_products"`
	TotalStock    int64               `json:"total_stock"`
	AveragePrice  float64             `json:"average_price"`
	Products      []ProductReportItem `json:"products"`
}

type ProductReportItem struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CategoryName  string `json:"category_name"`
	Price         int    `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}
