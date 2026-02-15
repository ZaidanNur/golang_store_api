package usecase

import (
	"context"
	"errors"
	"test-elabram/internal/domain"
	"test-elabram/internal/dto"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecase) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.GetAll(ctx)
}

func (u *categoryUsecase) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return u.categoryRepo.GetByID(ctx, id)
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) error {
	if category.Name == "" || category.Description == "" {
		return errors.New("name and description are required")
	}
	return u.categoryRepo.Create(ctx, category)
}

func (u *categoryUsecase) EditCategory(ctx context.Context, id int, category *dto.UpdateCategoryRequest) (*domain.Category, error) {
	existingCategory, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category.Name != "" {
		existingCategory.Name = category.Name
	}
	if category.Description != "" {
		existingCategory.Description = category.Description
	}

	if err := u.categoryRepo.Edit(ctx, existingCategory); err != nil {
		return nil, err
	}
	return existingCategory, nil
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, id int) error {
	return u.categoryRepo.Delete(ctx, id)
}
