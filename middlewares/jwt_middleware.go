package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"os"
	"silkroad-backend/cache"
	"silkroad-backend/i18n"
	"silkroad-backend/models"
	"silkroad-backend/utils"
)

func JWTMiddleware() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Cookies("JWT")
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "unauthorized")
		res := utils.Fail(msg)
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(res)
		}
		tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
			}

			return []byte(os.Getenv(utils.JWTSecretKey)), nil
		})
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(res)
		}

		claims, ok := tokenByte.Claims.(jwt.MapClaims)
		if !ok || !tokenByte.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(res)
		}

		if claims["iss"] != cache.LoadTextValue(models.AdminName) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(res)
		}

		return ctx.Next()
	}
}
