package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"procurement-system/config"
	"procurement-system/database"
	"procurement-system/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PurchaseItemRequest struct {
	ItemID uint `json:"item_id"`
	Qty    int  `json:"qty"`
}

type CreatePurchaseRequest struct {
	SupplierID uint                  `json:"supplier_id"`
	Items      []PurchaseItemRequest `json:"items"`
}

// GetAllPurchases returns all purchases with details
func GetAllPurchases(c *fiber.Ctx) error {
	var purchases []models.Purchasing
	if result := database.DB.Preload("Supplier").Preload("User").Preload("PurchasingDetails.Item").Find(&purchases); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch purchases",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    purchases,
	})
}

// GetPurchase returns a single purchase by ID
func GetPurchase(c *fiber.Ctx) error {
	id := c.Params("id")

	var purchase models.Purchasing
	if result := database.DB.Preload("Supplier").Preload("User").Preload("PurchasingDetails.Item").First(&purchase, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Purchase not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    purchase,
	})
}

// CreatePurchase creates a new purchase transaction with ACID compliance
func CreatePurchase(c *fiber.Ctx) error {
	var req CreatePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validation
	if req.SupplierID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Supplier ID is required",
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "At least one item is required",
		})
	}

	// Get user ID from JWT context
	userID := c.Locals("userID").(uint)

	// Check if supplier exists
	var supplier models.Supplier
	if result := database.DB.First(&supplier, req.SupplierID); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Supplier not found",
		})
	}

	// Start transaction for ACID compliance
	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to start transaction",
		})
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create purchase header
	purchase := models.Purchasing{
		Date:       time.Now(),
		SupplierID: req.SupplierID,
		UserID:     userID,
		GrandTotal: 0,
	}

	if result := tx.Create(&purchase); result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create purchase header",
		})
	}

	var grandTotal float64 = 0

	// Process each item
	for _, itemReq := range req.Items {
		// Validate item request
		if itemReq.ItemID == 0 || itemReq.Qty <= 0 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid item data: item_id and qty must be positive",
			})
		}

		// Get item from database (for price and stock validation)
		var item models.Item
		if result := tx.First(&item, itemReq.ItemID); result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("Item with ID %d not found", itemReq.ItemID),
			})
		}

		// Check stock availability
		if item.Stock < itemReq.Qty {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("Insufficient stock for item '%s'. Available: %d, Requested: %d", item.Name, item.Stock, itemReq.Qty),
			})
		}

		// Calculate subtotal using price from database (NOT from request!)
		subTotal := item.Price * float64(itemReq.Qty)
		grandTotal += subTotal

		// Create purchase detail
		detail := models.PurchasingDetail{
			PurchasingID: purchase.ID,
			ItemID:       itemReq.ItemID,
			Qty:          itemReq.Qty,
			SubTotal:     subTotal,
		}

		if result := tx.Create(&detail); result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to create purchase detail",
			})
		}

		// Update item stock (deduct)
		if result := tx.Model(&item).Update("stock", item.Stock-itemReq.Qty); result.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update item stock",
			})
		}
	}

	// Update grand total
	if result := tx.Model(&purchase).Update("grand_total", grandTotal); result.Error != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update grand total",
		})
	}

	// Commit transaction
	if result := tx.Commit(); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to commit transaction",
		})
	}

	// Reload purchase with all relations
	var completePurchase models.Purchasing
	database.DB.Preload("Supplier").Preload("User").Preload("PurchasingDetails.Item").First(&completePurchase, purchase.ID)

	// Send webhook notification (async, don't block response)
	go sendWebhookNotification(completePurchase)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Purchase created successfully",
		"data":    completePurchase,
	})
}

// sendWebhookNotification sends purchase data to configured webhook URL
func sendWebhookNotification(purchase models.Purchasing) {
	webhookURL := config.AppConfig.WebhookURL
	if webhookURL == "" {
		log.Println("Webhook URL not configured, skipping notification")
		return
	}

	// Prepare webhook payload
	payload := map[string]interface{}{
		"event":      "purchase_created",
		"timestamp":  time.Now().Format(time.RFC3339),
		"order_id":   purchase.ID,
		"date":       purchase.Date.Format("2006-01-02"),
		"supplier":   purchase.Supplier.Name,
		"user":       purchase.User.Username,
		"grand_total": purchase.GrandTotal,
		"items":      make([]map[string]interface{}, 0),
	}

	for _, detail := range purchase.PurchasingDetails {
		item := map[string]interface{}{
			"item_id":   detail.ItemID,
			"item_name": detail.Item.Name,
			"qty":       detail.Qty,
			"price":     detail.Item.Price,
			"sub_total": detail.SubTotal,
		}
		payload["items"] = append(payload["items"].([]map[string]interface{}), item)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal webhook payload: %v", err)
		return
	}

	// Send HTTP POST request
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Failed to send webhook notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Webhook notification sent successfully for order #%d", purchase.ID)
	} else {
		log.Printf("Webhook notification failed with status: %d", resp.StatusCode)
	}
}
