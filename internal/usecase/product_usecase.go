package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"test-elabram/internal/cache"
	"test-elabram/internal/domain"
	"test-elabram/internal/dto"
	"time"
)

const reportCacheTTL = 5 * time.Minute

var productCacheKey = map[string]string{
	"report": "product:report",
}

type productUsecase struct {
	productRepository domain.ProductRepository
	cache             *cache.RedisCache
}

func NewProductUsecase(productRepository domain.ProductRepository, redisCache *cache.RedisCache) domain.ProductUsecase {
	return &productUsecase{
		productRepository: productRepository,
		cache:             redisCache,
	}
}

func (u *productUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return u.productRepository.GetAll(ctx)
}

func (u *productUsecase) GetAllProductsPaginated(ctx context.Context, params dto.ProductFilterParams, pq dto.PaginationQuery) (*dto.PaginatedResponse, error) {
	if pq.Page <= 0 {
		pq.Page = 1
	}
	if pq.Limit <= 0 {
		pq.Limit = 10
	}
	if pq.Limit > 100 {
		pq.Limit = 100
	}

	products, total, err := u.productRepository.GetAllPaginated(ctx, params, pq)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pq.Limit)))

	return &dto.PaginatedResponse{
		Data:       products,
		Page:       pq.Page,
		Limit:      pq.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (u *productUsecase) GetProductReport(ctx context.Context) (*dto.ProductReportResponse, error) {
	if u.cache != nil && u.cache.IsAvailable() {
		cached, err := u.cache.Get(ctx, productCacheKey["report"])
		if err == nil && cached != nil {
			var report dto.ProductReportResponse
			if json.Unmarshal(cached, &report) == nil {
				return &report, nil
			}
		}
	}

	report, err := u.productRepository.GetProductReport(ctx)
	if err != nil {
		return nil, err
	}

	if u.cache != nil && u.cache.IsAvailable() {
		data, err := json.Marshal(report)
		if err == nil {
			if cacheErr := u.cache.Set(ctx, productCacheKey["report"], data, reportCacheTTL); cacheErr != nil {
				log.Printf("[CACHE] Failed to cache report: %v", cacheErr)
			}
		}
	}

	return report, nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return u.productRepository.GetByID(ctx, id)
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	err := u.productRepository.Create(ctx, product)
	if err == nil {
		u.invalidateReportCache(ctx, productCacheKey["report"])
	}
	return err
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
	u.invalidateReportCache(ctx, productCacheKey["report"])
	return product, nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid ID")
	}
	err := u.productRepository.Delete(ctx, id)
	if err == nil {
		u.invalidateReportCache(ctx, productCacheKey["report"])
	}
	return err
}

func (u *productUsecase) invalidateReportCache(ctx context.Context, cacheKey string) {
	if u.cache != nil && u.cache.IsAvailable() {
		if err := u.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("[CACHE] Failed to invalidate report cache: %v", err)
		}
	}
}
