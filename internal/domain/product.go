package domain

import (
	"context"
	"test-elabram/internal/dto"
	"time"
)

type Product struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	Name          string    `json:"name" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
	Price         int       `json:"price" gorm:"not null"`
	StockQuantity int       `json:"stock_quantity" gorm:"not null"`
	IsActive      bool      `json:"is_active" gorm:"not null"`
	CategoryID    uint      `json:"category_id" gorm:"not null;index"`
	Category      Category  `json:"category" gorm:"foreignKey:CategoryID"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductRepository interface {
	GetAll(ctx context.Context) ([]Product, error)
	GetAllPaginated(ctx context.Context, params dto.ProductFilterParams, pq dto.PaginationQuery) ([]Product, int64, error)
	GetProductReport(ctx context.Context) (*dto.ProductReportResponse, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Edit(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int) error
}

type ProductUsecase interface {
	GetAllProducts(ctx context.Context) ([]Product, error)
	GetAllProductsPaginated(ctx context.Context, params dto.ProductFilterParams, pq dto.PaginationQuery) (*dto.PaginatedResponse, error)
	GetProductReport(ctx context.Context) (*dto.ProductReportResponse, error)
	GetProductByID(ctx context.Context, id int) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	EditProduct(ctx context.Context, id int, req *dto.UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id int) error
}
