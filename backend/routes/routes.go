package routes

import (
	"procurement-system/handlers"
	"procurement-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// API group
	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware())
	
	// Profile
	protected.Get("/profile", handlers.GetProfile)

	// Items CRUD
	items := protected.Group("/items")
	items.Get("/", handlers.GetAllItems)
	items.Get("/:id", handlers.GetItem)
	items.Post("/", handlers.CreateItem)
	items.Put("/:id", handlers.UpdateItem)
	items.Delete("/:id", handlers.DeleteItem)

	// Suppliers CRUD
	suppliers := protected.Group("/suppliers")
	suppliers.Get("/", handlers.GetAllSuppliers)
	suppliers.Get("/:id", handlers.GetSupplier)
	suppliers.Post("/", handlers.CreateSupplier)
	suppliers.Put("/:id", handlers.UpdateSupplier)
	suppliers.Delete("/:id", handlers.DeleteSupplier)

	// Purchasing
	purchases := protected.Group("/purchases")
	purchases.Get("/", handlers.GetAllPurchases)
	purchases.Get("/:id", handlers.GetPurchase)
	purchases.Post("/", handlers.CreatePurchase)
}
