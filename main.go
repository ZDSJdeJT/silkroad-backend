package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	_ "silkroad-backend/docs"
	"silkroad-backend/pkg/middlewares"
	"silkroad-backend/pkg/routes"
	"silkroad-backend/pkg/utils"
	"silkroad-backend/platform/database"
)

// @title Silk Road
// @version 1.0.0
// @description The API doc of Silk Road.

// @BasePath /api
// @schemes http
func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Get error message from the error interface and send it as response
			errorMessage := fmt.Errorf("%v", err)
			return ctx.Status(code).JSON(utils.Fail(errorMessage.Error()))
		},
	})

	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		err := godotenv.Load(".env.production")
		if err != nil {
			log.Fatalf("Error loading .env.production file: %s.", err)
		}
		if _, err := utils.CheckEnvVarsExist([]string{"PORT", "APP_NAME", "APP_VERSION", "DATABASE_DSN"}); err != nil {
			log.Fatalf("Error starting: %s.", err)
		}
		middlewares.FiberMiddlewares(app, false)
	default:
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatalf("Error loading .env.development file: %s.", err)
		}
		if _, err := utils.CheckEnvVarsExist([]string{"PORT", "APP_NAME", "APP_VERSION", "DATABASE_DSN"}); err != nil {
			log.Fatalf("Error starting: %s.", err)
		}
		middlewares.FiberMiddlewares(app, true)
		routes.SwaggerRoutes(app)
	}

	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Error initing database: %s.", err)
		return
	}

	err := utils.InitClientHTML()
	if err != nil {
		log.Fatalf("Error initing client HTML: %s.", err)
		return
	}

	routes.APIRoutes(app)
	routes.SPARoutes(app)

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Error starting: %s.", err)
	}
}
