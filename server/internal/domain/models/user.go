package models

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Phone     string    `db:"phone" json:"phone"`
	IsActive  bool      `db:"is_active" json:"is_active"`

	// OAuth2 fields
	GoogleID     string `db:"google_id" json:"google_id,omitempty"`
	GithubID     string `db:"github_id" json:"github_id,omitempty"`
	AvatarURL    string `db:"avatar_url" json:"avatar_url,omitempty"`
	AuthProvider string `db:"auth_provider" json:"auth_provider,omitempty"` // "email", "google", "github"
}

type UserResponse struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	AuthProvider string    `json:"auth_provider,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=15"`
}

type OAuthRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google github"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
}

type OAuthCallbackRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google github"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
