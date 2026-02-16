package repository

import (
	"tenderness/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	query := `SELECT * FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(&products, query, limit, offset)
	return products, err
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var product models.Product
	query := `SELECT * FROM products WHERE id = $1`
	err := r.db.Get(&product, query, id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByCategory(category string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	query := `SELECT * FROM products WHERE category = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.Select(&products, query, category, limit, offset)
	return products, err
}

func (r *ProductRepository) GetFeatured(limit int) ([]models.Product, error) {
	var products []models.Product
	query := `SELECT * FROM products WHERE in_stock = true ORDER BY rating DESC, views DESC LIMIT $1`
	err := r.db.Select(&products, query, limit)
	return products, err
}

func (r *ProductRepository) Search(query string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	searchPattern := "%" + query + "%"
	sqlQuery := `SELECT * FROM products WHERE name ILIKE $1 OR description ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.Select(&products, sqlQuery, searchPattern, limit, offset)
	return products, err
}

func (r *ProductRepository) IncrementViews(id int) error {
	query := `UPDATE products SET views = views + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ProductRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM products`
	err := r.db.Get(&count, query)
	return count, err
}
