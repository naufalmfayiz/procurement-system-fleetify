package handlers

import (
	"procurement-system/database"
	"procurement-system/models"

	"github.com/gofiber/fiber/v2"
)

type CreateItemRequest struct {
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}

type UpdateItemRequest struct {
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}

// GetAllItems returns all items
func GetAllItems(c *fiber.Ctx) error {
	var items []models.Item
	if result := database.DB.Find(&items); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch items",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
	})
}

// GetItem returns a single item by ID
func GetItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Item
	if result := database.DB.First(&item, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Item not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    item,
	})
}

// CreateItem creates a new item
func CreateItem(c *fiber.Ctx) error {
	var req CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validation
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Item name is required",
		})
	}

	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Price cannot be negative",
		})
	}

	if req.Stock < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Stock cannot be negative",
		})
	}

	item := models.Item{
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	if result := database.DB.Create(&item); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Item created successfully",
		"data":    item,
	})
}

// UpdateItem updates an existing item
func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Item
	if result := database.DB.First(&item, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Item not found",
		})
	}

	var req UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validation
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Item name is required",
		})
	}

	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Price cannot be negative",
		})
	}

	if req.Stock < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Stock cannot be negative",
		})
	}

	item.Name = req.Name
	item.Stock = req.Stock
	item.Price = req.Price

	if result := database.DB.Save(&item); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update item",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Item updated successfully",
		"data":    item,
	})
}

// DeleteItem soft deletes an item
func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Item
	if result := database.DB.First(&item, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Item not found",
		})
	}

	if result := database.DB.Delete(&item); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete item",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Item deleted successfully",
	})
}
