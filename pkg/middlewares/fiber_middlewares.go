package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"silkroad-backend/pkg/utils"
	"strings"
	"time"
)

func FiberMiddlewares(app *fiber.App, enableCors bool) {
	if enableCors {
		app.Use(cors.New())
	}
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(helmet.New())
	app.Use(idempotency.New())
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(skip.New(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(ctx *fiber.Ctx) error {
			res := utils.Fail("请求过于频繁，请稍后再试！")
			return ctx.Status(fiber.StatusTooManyRequests).JSON(res)
		},
	}), func(ctx *fiber.Ctx) bool {
		return !strings.HasPrefix(ctx.Path(), "/api")
	}))
}
