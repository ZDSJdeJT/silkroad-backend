package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"silkroad-backend/i18n"
	"silkroad-backend/utils"
	"time"
)

func LimiterMiddleware() func(*fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:               30,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(ctx *fiber.Ctx) error {
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "limiter")
			return ctx.Status(fiber.StatusTooManyRequests).JSON(utils.Fail(msg))
		},
	})
}
