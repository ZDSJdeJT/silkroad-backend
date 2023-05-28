package middlewares

import (
	"os"
	"silkroad-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"

	jwtMiddleware "github.com/gofiber/jwt/v2"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(ctx *fiber.Ctx, _ error) error {
	// Return status 401 and failed authentication error.
	return ctx.Status(fiber.StatusUnauthorized).JSON(utils.Fail("请登录后再试"))
}
