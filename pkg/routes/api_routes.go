package routes

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/app/controllers"
	"silkroad-backend/pkg/middlewares"
)

func APIRoutes(app *fiber.App) {
	route := app.Group("/api/v1")

	// 公开接口
	route.Get("/public-settings", controllers.GetPublicSettings)
	route.Post("/admin/login", controllers.AdminLogin)
	// 需要鉴权的接口
	route.Use(middlewares.JWTProtected())
	route.Get("/system/info", controllers.GetSystemInfo)
	route.Get("/settings", controllers.GetSettings)
	route.Put("/admin/settings/:key", controllers.UpdateSetting)
}
