package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func CommonMiddlewares(app *fiber.App) {
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(helmet.New())
	app.Use(idempotency.New())
	app.Use(recover.New())
	app.Use(logger.New())
}
