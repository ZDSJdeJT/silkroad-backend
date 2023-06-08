package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/i18n"
	"strings"
)

func SetLocals(ctx *fiber.Ctx) {
	acceptLanguage := ctx.Get("Accept-Language", i18n.DefaultLanguage)
	for _, lang := range i18n.Languages {
		if strings.Contains(acceptLanguage, lang) {
			ctx.Locals("lang", lang)
			return
		}
	}
	ctx.Locals("lang", i18n.DefaultLanguage)
}

func I18nMiddleware() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		SetLocals(ctx)
		return ctx.Next()
	}
}
