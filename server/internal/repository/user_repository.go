package repository

import (
	"tenderness/internal/domain/models"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (email, password, first_name, last_name, phone, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(query,
		user.Email,
		string(hashedPassword),
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, email, password, first_name, last_name, phone, is_active 
			  FROM users WHERE email = $1 AND is_active = true`

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, email, password, first_name, last_name, phone, is_active 
			  FROM users WHERE id = $1 AND is_active = true`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET first_name = $2, last_name = $3, phone = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query, user.ID, user.FirstName, user.LastName, user.Phone)
	return err
}

func (r *UserRepository) UpdatePassword(userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err = r.db.Exec(query, userID, string(hashedPassword))
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.Get(&exists, query, email)
	return exists, err
}

func (r *UserRepository) ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// OAuth2 methods
func (r *UserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, email, first_name, last_name, phone, is_active, google_id, github_id, avatar_url, auth_provider FROM users WHERE google_id = $1`
	err := r.db.Get(&user, query, googleID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByGithubID(githubID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, created_at, updated_at, email, first_name, last_name, phone, is_active, google_id, github_id, avatar_url, auth_provider FROM users WHERE github_id = $1`
	err := r.db.Get(&user, query, githubID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateOAuth(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (email, first_name, last_name, phone, is_active, google_id, github_id, avatar_url, auth_provider, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at, email, first_name, last_name, phone, is_active, google_id, github_id, avatar_url, auth_provider`

	var createdUser models.User
	err := r.db.QueryRow(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
		user.GoogleID,
		user.GithubID,
		user.AvatarURL,
		user.AuthProvider,
	).Scan(
		&createdUser.ID,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.Email,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Phone,
		&createdUser.IsActive,
		&createdUser.GoogleID,
		&createdUser.GithubID,
		&createdUser.AvatarURL,
		&createdUser.AuthProvider,
	)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) LinkOAuth(userID int, provider string, googleID, githubID, avatarURL string) error {
	query := `
		UPDATE users 
		SET google_id = $2, github_id = $3, avatar_url = $4, auth_provider = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	_, err := r.db.Exec(query, provider, userID, googleID, githubID, avatarURL)
	return err
}
