package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"tenderness/internal/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "12"))

	products, total, err := h.productService.GetProducts(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

func (h *ProductHandler) GetProductsByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "12"))

	products, total, err := h.productService.GetProductsByCategory(category, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *ProductHandler) GetFeaturedProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "8"))

	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch featured products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
	})
}

func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "12"))

	products, total, err := h.productService.SearchProducts(query, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"query":    query,
	})
}

func (h *ProductHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.productService.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(fiber.Map{
		"categories": categories,
	})
}
