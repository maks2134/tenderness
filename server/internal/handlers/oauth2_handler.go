package handlers

import (
	"tenderness/internal/domain/models"
	"tenderness/internal/services"

	"github.com/gofiber/fiber/v2"
)

type OAuth2Handler struct {
	oauth2Service *services.OAuth2Service
}

func NewOAuth2Handler(oauth2Service *services.OAuth2Service) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2Service: oauth2Service,
	}
}

func (h *OAuth2Handler) GetAuthURL(c *fiber.Ctx) error {
	provider := c.Params("provider")
	state, err := h.oauth2Service.GenerateState()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate state",
		})
	}

	authURL, err := h.oauth2Service.GetAuthURL(provider, state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider",
		})
	}

	return c.JSON(fiber.Map{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *OAuth2Handler) Callback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	response, err := h.oauth2Service.ExchangeCode(provider, code, state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange code: " + err.Error(),
		})
	}

	// Set JWT token in cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    response.Token,
		HTTPOnly: true,
		SameSite: "lax",
		MaxAge:   86400, // 24 hours
	})

	// Redirect to frontend
	return c.Redirect("http://localhost/auth/success?token="+response.Token, fiber.StatusTemporaryRedirect)
}

func (h *OAuth2Handler) LinkAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.OAuthCallbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	err := h.oauth2Service.LinkOAuthToAccount(userID, req.Provider, req.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to link account: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account linked successfully",
	})
}

func (h *OAuth2Handler) UnlinkAccount(c *fiber.Ctx) error {
	provider := c.Params("provider")

	if provider != "google" && provider != "github" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider",
		})
	}

	// TODO: Implement unlink functionality in repository
	return c.JSON(fiber.Map{
		"message": "Account unlinked successfully",
	})
}
