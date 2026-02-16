package services

import (
	"tenderness/internal/domain/models"
	"tenderness/internal/repository"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetProducts(page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	products, err := s.productRepo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.productRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	go s.productRepo.IncrementViews(id)

	return product, nil
}

func (s *ProductService) GetProductsByCategory(category string, page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	products, err := s.productRepo.GetByCategory(category, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.productRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) GetFeaturedProducts(limit int) ([]models.Product, error) {
	if limit < 1 {
		limit = 8
	}
	if limit > 20 {
		limit = 20
	}

	return s.productRepo.GetFeatured(limit)
}

func (s *ProductService) SearchProducts(query string, page, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	products, err := s.productRepo.Search(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.productRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}
