package main

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"silkroad-backend/cache"
	"silkroad-backend/cron"
	"silkroad-backend/database"
	_ "silkroad-backend/docs"
	"silkroad-backend/i18n"
	"silkroad-backend/middlewares"
	"silkroad-backend/routes"
	"silkroad-backend/utils"
)

// @title Silk Road
// @version 1.0.0
// @description The API doc of Silk Road.

// @BasePath /api
// @schemes http
func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		EnablePrintRoutes: true,
		AppName:           utils.APPName,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			if _, ok := ctx.Locals("lang").(string); !ok {
				middlewares.SetLocals(ctx)
			}

			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "internalServerError")
			return ctx.Status(code).JSON(utils.Fail(msg))
		},
	})

	env := os.Getenv(utils.AppMode)
	switch env {
	case "production":
		err := godotenv.Load(".env.production")
		if err != nil {
			log.Fatalf("Error loading .env.production file: %s.", err)
		}
		err = utils.CheckEnvVarsExist()
		if err != nil {
			log.Fatalf("Error starting: %s.", err)
		}
	default:
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatalf("Error loading .env.development file: %s.", err)
		}
		err = utils.CheckEnvVarsExist()
		if err != nil {
			log.Fatalf("Error starting: %s.", err)
		}
		routes.SwaggerRoutes(app)
	}

	settings, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %s.", err)
	}

	err = cache.InitCacheAndClientHTML(settings)
	if err != nil {
		log.Fatalf("Error initializing cache and client HTML: %s.", err)
	}

	middlewares.CommonMiddlewares(app)

	routes.APIRoutes(app)
	routes.SPARoutes(app)

	cron.Start()

	if err := app.Listen(":" + utils.APPPort); err != nil {
		log.Fatalf("Error starting: %s.", err)
	}
}
