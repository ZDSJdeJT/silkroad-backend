package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func SwaggerRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
