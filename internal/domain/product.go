package domain

import "time"

type Product struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	Name          string    `json:"name" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
	Price         int       `json:"price" gorm:"not null"`
	StockQuantity int       `json:"stock_quantity" gorm:"not null"`
	IsActive      bool      `json:"is_active" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(product *Product) error
}

type ProductUsecase interface {
	GetAllProducts() ([]Product, error)
	GetProductByID(id int) (*Product, error)
	CreateProduct(product *Product) error
}
