package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"silkroad-backend/pkg/utils"
	"time"
)

func FiberMiddlewares(app *fiber.App, enableCors bool) {
	if enableCors {
		app.Use(cors.New())
	}
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(idempotency.New())
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(limiter.New(limiter.Config{
		Max:               40,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(ctx *fiber.Ctx) error {
			res := utils.Fail("请求过于频繁，请稍后再试！")
			return ctx.JSON(res)
		},
	}))
}
