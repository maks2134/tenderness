package services

import (
	"errors"
	"reflect"
	"tenderness/internal/domain/models"
	"tenderness/internal/middleware"
	"tenderness/internal/repository"

	"github.com/go-playground/validator/v10"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwt       *middleware.JWTMiddleware
	validator *validator.Validate
}

func NewAuthService(userRepo *repository.UserRepository, jwt *middleware.JWTMiddleware) *AuthService {
	v := validator.New()

	// Регистрируем кастомный валидатор для использования JSON имен вместо имен полей структуры
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("json")
		if name == "" {
			return field.Name
		}
		return name
	})

	return &AuthService{
		userRepo:  userRepo,
		jwt:       jwt,
		validator: v,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("user with this email already exists")
	}

	user := &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return &models.AuthResponse{
		User:  *userResponse,
		Token: token,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := s.userRepo.ValidatePassword(req.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return &models.AuthResponse{
		User:  *userResponse,
		Token: token,
	}, nil
}

func (s *AuthService) GetProfile(userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return userResponse, nil
}

func (s *AuthService) UpdateProfile(userID int, req *models.User) error {
	if err := s.validator.Struct(req); err != nil {
		return err
	}

	user := &models.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	return s.userRepo.Update(user)
}

func (s *AuthService) ChangePassword(userID int, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := s.userRepo.ValidatePassword(currentPassword, user.Password); err != nil {
		return errors.New("current password is incorrect")
	}

	return s.userRepo.UpdatePassword(userID, newPassword)
}

func (s *AuthService) DeleteAccount(userID int) error {
	return s.userRepo.Delete(userID)
}
