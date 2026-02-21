package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"tenderness/internal/domain/models"
	"tenderness/internal/middleware"
	"tenderness/internal/repository"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2Service struct {
	userRepo *repository.UserRepository
	jwt      *middleware.JWTMiddleware
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func NewOAuth2Service(userRepo *repository.UserRepository, jwt *middleware.JWTMiddleware) *OAuth2Service {
	return &OAuth2Service{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

// Google OAuth2 configuration
var googleOAuth2Config = &oauth2.Config{
	ClientID:     getEnv("GOOGLE_CLIENT_ID", "your-google-client-id"),
	ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "your-google-client-secret"),
	RedirectURL:  "http://localhost/api/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// GitHub OAuth2 configuration
var githubOAuth2Config = &oauth2.Config{
	ClientID:     getEnv("GITHUB_CLIENT_ID", "your-github-client-id"),
	ClientSecret: getEnv("GITHUB_CLIENT_SECRET", "your-github-client-secret"),
	RedirectURL:  "http://localhost/api/auth/github/callback",
	Scopes:       []string{"user:email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *OAuth2Service) GetAuthURL(provider, state string) (string, error) {
	switch provider {
	case "google":
		return googleOAuth2Config.AuthCodeURL(state), nil
	case "github":
		return githubOAuth2Config.AuthCodeURL(state), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *OAuth2Service) ExchangeCode(provider, code, state string) (*models.AuthResponse, error) {
	switch provider {
	case "google":
		return s.handleGoogleCallback(code, state)
	case "github":
		return s.handleGitHubCallback(code, state)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *OAuth2Service) handleGoogleCallback(code, state string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := googleOAuth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from Google
	client := googleOAuth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetByGoogleID(userInfo.ID)
	if err != nil {
		// User doesn't exist, create new user
		names := strings.Split(userInfo.Name, " ")
		firstName := names[0]
		lastName := ""
		if len(names) > 1 {
			lastName = strings.Join(names[1:], " ")
		}

		newUser := &models.User{
			Email:        userInfo.Email,
			FirstName:    firstName,
			LastName:     lastName,
			AvatarURL:    userInfo.Picture,
			GoogleID:     userInfo.ID,
			AuthProvider: "google",
			IsActive:     true,
		}

		createdUser, err := s.userRepo.CreateOAuth(newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		user = createdUser
	}

	// Generate JWT token
	jwtToken, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	userResponse := &models.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt,
		AvatarURL:    user.AvatarURL,
		AuthProvider: user.AuthProvider,
	}

	return &models.AuthResponse{
		User:  *userResponse,
		Token: jwtToken,
	}, nil
}

func (s *OAuth2Service) handleGitHubCallback(code, state string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := githubOAuth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from GitHub
	client := githubOAuth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Get user email (requires separate API call)
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}
	defer emailResp.Body.Close()

	emailBody, err := io.ReadAll(emailResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read email response: %w", err)
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.Unmarshal(emailBody, &emails); err != nil {
		return nil, fmt.Errorf("failed to parse emails: %w", err)
	}

	var primaryEmail string
	for _, email := range emails {
		if email.Primary {
			primaryEmail = email.Email
			break
		}
	}

	if primaryEmail == "" {
		return nil, errors.New("no primary email found")
	}

	// Check if user exists
	user, err := s.userRepo.GetByGithubID(fmt.Sprintf("%d", userInfo.ID))
	if err != nil {
		// User doesn't exist, create new user
		names := strings.Split(userInfo.Name, " ")
		firstName := names[0]
		lastName := ""
		if len(names) > 1 {
			lastName = strings.Join(names[1:], " ")
		}

		newUser := &models.User{
			Email:        primaryEmail,
			FirstName:    firstName,
			LastName:     lastName,
			AvatarURL:    userInfo.AvatarURL,
			GithubID:     fmt.Sprintf("%d", userInfo.ID),
			AuthProvider: "github",
			IsActive:     true,
		}

		createdUser, err := s.userRepo.CreateOAuth(newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		user = createdUser
	}

	// Generate JWT token
	jwtToken, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	userResponse := &models.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt,
		AvatarURL:    user.AvatarURL,
		AuthProvider: user.AuthProvider,
	}

	return &models.AuthResponse{
		User:  *userResponse,
		Token: jwtToken,
	}, nil
}

func (s *OAuth2Service) GenerateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *OAuth2Service) LinkOAuthToAccount(userID int, provider, code string) error {
	var user *models.User
	var err error

	switch provider {
	case "google":
		user, err = s.handleGoogleUserCreation(code)
	case "github":
		user, err = s.handleGitHubUserCreation(code)
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		return fmt.Errorf("failed to get %s user: %w", provider, err)
	}

	// Link OAuth to existing user
	return s.userRepo.LinkOAuth(userID, provider, user.GoogleID, user.GithubID, user.AvatarURL)
}

func (s *OAuth2Service) handleGoogleUserCreation(code string) (*models.User, error) {
	token, err := googleOAuth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := googleOAuth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	names := strings.Split(userInfo.Name, " ")
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	return &models.User{
		GoogleID:  userInfo.ID,
		AvatarURL: userInfo.Picture,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func (s *OAuth2Service) handleGitHubUserCreation(code string) (*models.User, error) {
	token, err := githubOAuth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := githubOAuth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	names := strings.Split(userInfo.Name, " ")
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	return &models.User{
		GithubID:  fmt.Sprintf("%d", userInfo.ID),
		AvatarURL: userInfo.AvatarURL,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
