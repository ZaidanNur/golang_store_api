package usecase

import (
	"context"
	"errors"
	"test-elabram/internal/domain"
	"test-elabram/internal/dto"
)

type productUsecase struct {
	productRepository domain.ProductRepository
}

func NewProductUsecase(productRepository domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepository: productRepository,
	}
}

func (u *productUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return u.productRepository.GetAll(ctx)
}

func (u *productUsecase) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return u.productRepository.GetByID(ctx, id)
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	return u.productRepository.Create(ctx, product)
}

func (u *productUsecase) EditProduct(ctx context.Context, id int, req *dto.UpdateProductRequest) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}

	product, err := u.productRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}

	if err := u.productRepository.Edit(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid ID")
	}
	return u.productRepository.Delete(ctx, id)
}
