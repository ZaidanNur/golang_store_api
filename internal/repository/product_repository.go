package repository

import (
	"context"
	"fmt"
	"test-elabram/internal/domain"
	"test-elabram/internal/dto"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.WithContext(ctx).Preload("Category").Find(&products).Error
	return products, err
}

var allowedSortColumns = map[string]bool{
	"name":           true,
	"price":          true,
	"stock_quantity": true,
	"created_at":     true,
	"category_id":    true,
}

func (r *productRepository) GetAllPaginated(ctx context.Context, params dto.ProductFilterParams, pq dto.PaginationQuery) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{})

	if params.Name != "" {
		query = query.Where("name ILIKE ?", "%"+params.Name+"%")
	}
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}
	if params.PriceMin != nil {
		query = query.Where("price >= ?", *params.PriceMin)
	}
	if params.PriceMax != nil {
		query = query.Where("price <= ?", *params.PriceMax)
	}
	if params.StockMin != nil {
		query = query.Where("stock_quantity >= ?", *params.StockMin)
	}
	if params.StockMax != nil {
		query = query.Where("stock_quantity <= ?", *params.StockMax)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "created_at"
	if allowedSortColumns[params.SortBy] {
		sortBy = params.SortBy
	}
	sortOrder := "desc"
	if params.SortOrder == "asc" {
		sortOrder = "asc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	offset := (pq.Page - 1) * pq.Limit
	err := query.Offset(offset).Limit(pq.Limit).Preload("Category").Find(&products).Error
	return products, total, err
}

func (r *productRepository) GetProductReport(ctx context.Context) (*dto.ProductReportResponse, error) {
	var report dto.ProductReportResponse

	row := r.db.WithContext(ctx).Model(&domain.Product{}).
		Select("COUNT(*) as total_products, COALESCE(SUM(stock_quantity), 0) as total_stock, COALESCE(AVG(price), 0) as average_price").
		Row()
	if err := row.Scan(&report.TotalProducts, &report.TotalStock, &report.AveragePrice); err != nil {
		return nil, err
	}

	var products []domain.Product
	err := r.db.WithContext(ctx).
		Select("id, name, price, stock_quantity, category_id").
		Preload("Category").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	report.Products = make([]dto.ProductReportItem, len(products))
	for i, p := range products {
		report.Products[i] = dto.ProductReportItem{
			ID:            p.ID,
			Name:          p.Name,
			CategoryName:  p.Category.Name,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
		}
	}

	return &report, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Edit(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}
