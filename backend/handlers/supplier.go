package handlers

import (
	"procurement-system/database"
	"procurement-system/models"

	"github.com/gofiber/fiber/v2"
)

type CreateSupplierRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type UpdateSupplierRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// GetAllSuppliers returns all suppliers
func GetAllSuppliers(c *fiber.Ctx) error {
	var suppliers []models.Supplier
	if result := database.DB.Find(&suppliers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch suppliers",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    suppliers,
	})
}

// GetSupplier returns a single supplier by ID
func GetSupplier(c *fiber.Ctx) error {
	id := c.Params("id")

	var supplier models.Supplier
	if result := database.DB.First(&supplier, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Supplier not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    supplier,
	})
}

// CreateSupplier creates a new supplier
func CreateSupplier(c *fiber.Ctx) error {
	var req CreateSupplierRequest
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
			"message": "Supplier name is required",
		})
	}

	supplier := models.Supplier{
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
	}

	if result := database.DB.Create(&supplier); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create supplier",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Supplier created successfully",
		"data":    supplier,
	})
}

// UpdateSupplier updates an existing supplier
func UpdateSupplier(c *fiber.Ctx) error {
	id := c.Params("id")

	var supplier models.Supplier
	if result := database.DB.First(&supplier, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Supplier not found",
		})
	}

	var req UpdateSupplierRequest
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
			"message": "Supplier name is required",
		})
	}

	supplier.Name = req.Name
	supplier.Email = req.Email
	supplier.Address = req.Address

	if result := database.DB.Save(&supplier); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update supplier",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Supplier updated successfully",
		"data":    supplier,
	})
}

// DeleteSupplier soft deletes a supplier
func DeleteSupplier(c *fiber.Ctx) error {
	id := c.Params("id")

	var supplier models.Supplier
	if result := database.DB.First(&supplier, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Supplier not found",
		})
	}

	if result := database.DB.Delete(&supplier); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete supplier",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Supplier deleted successfully",
	})
}
