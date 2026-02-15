package domain

import (
	"context"
	"test-elabram/internal/dto"
	"time"
)

type Category struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	Create(ctx context.Context, category *Category) error
	Edit(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
}

type CategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetCategoryByID(ctx context.Context, id int) (*Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	EditCategory(ctx context.Context, id int, category *dto.UpdateCategoryRequest) (*Category, error)
	DeleteCategory(ctx context.Context, id int) error
}
