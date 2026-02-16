package repository

import (
	"tenderness/internal/domain/models"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT * FROM categories ORDER BY name`
	err := r.db.Select(&categories, query)
	return categories, err
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE id = $1`
	err := r.db.Get(&category, query, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (name, description, image_url) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, category.Name, category.Description, category.ImageURL).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}
