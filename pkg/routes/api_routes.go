package routes

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/app/controllers"
)

func APIRoutes(app *fiber.App) {
	// Create routes group.
	route := app.Group("/api/v1")

	// Routes for GET method:
	route.Get("/hello", controllers.Hello)
}
