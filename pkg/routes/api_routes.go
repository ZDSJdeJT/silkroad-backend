package routes

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/app/controllers"
)

func APIRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	route.Get("/system/info", controllers.GetSystemInfo)

	route.Get("/public-options", controllers.GetPublicOptions)

	route.Put("/admin/options", controllers.UpdateOption)
}
