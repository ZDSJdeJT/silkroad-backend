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
	publicRoutes.Post("/upload/files/:uuid", v1.UploadFile)
	publicRoutes.Post("/upload/files/merge/:uuid", limiterMiddleware, v1.MergeFile)
	publicRoutes.Post("/upload/texts", limiterMiddleware, v1.UploadText)
	publicRoutes.Get("/records/:code", limiterMiddleware, v1.GetRecordByCode)
	publicRoutes.Get("/receive/texts/:code", limiterMiddleware, v1.ReceiveText)
	publicRoutes.Get("/receive/files/:code", limiterMiddleware, v1.ReceiveFile)
	publicRoutes.Delete("/records/:id", limiterMiddleware, v1.DeleteRecord)

	adminRoutes := v1Routes.Group("/admin", limiterMiddleware)
	adminRoutes.Post("/login", v1.AdminLogin)
	adminRoutes.Use(JWTMiddleware)
	adminRoutes.Post("/logout", v1.AdminLogout)
	adminRoutes.Get("/system/info", v1.GetSystemInfo)
	adminRoutes.Get("/settings", v1.GetSettings)
	adminRoutes.Put("/settings/:key", v1.UpdateSetting)
	adminRoutes.Delete("/records/expired/text", v1.DeleteExpiredTextRecords)
	adminRoutes.Delete("/records/expired/file", v1.DeleteExpiredFileRecords)
	adminRoutes.Delete("/expired/chunks", v1.DeleteExpiredChunks)
}
