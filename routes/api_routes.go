package routes

import (
	"github.com/gofiber/fiber/v2"
	v1 "silkroad-backend/controllers/v1"
	"silkroad-backend/middlewares"
)

func APIRoutes(app *fiber.App) {
	limiterMiddleware := middlewares.LimiterMiddleware()
	i18nMiddleware := middlewares.I18nMiddleware()
	JWTMiddleware := middlewares.JWTMiddleware()

	routes := app.Group("/api")
	routes.Use(i18nMiddleware)

	v1Routes := routes.Group("/v1")

	publicRoutes := v1Routes.Group("/public")
	publicRoutes.Get("/system/languages", limiterMiddleware, v1.GetSystemLanguages)
	publicRoutes.Get("/settings", limiterMiddleware, v1.GetPublicSettings)
	publicRoutes.Post("/upload/file/:uuid", v1.UploadFile)
	publicRoutes.Post("/upload/file/merge/:uuid", limiterMiddleware, v1.MergeFile)
	publicRoutes.Post("/upload/text", limiterMiddleware, v1.UploadText)
	publicRoutes.Get("/receive/:code", limiterMiddleware, v1.Receive)
	publicRoutes.Delete("/text/:id", limiterMiddleware, v1.DeleteText)
	publicRoutes.Delete("/file/:id", limiterMiddleware, v1.DeleteFile)

	adminRoutes := v1Routes.Group("/admin", limiterMiddleware)
	adminRoutes.Post("/login", v1.AdminLogin)
	adminRoutes.Use(JWTMiddleware)
	adminRoutes.Post("/logout", v1.AdminLogout)
	adminRoutes.Get("/system/info", v1.GetSystemInfo)
	adminRoutes.Get("/settings", v1.GetSettings)
	adminRoutes.Put("/settings/:key", v1.UpdateSetting)
}
