package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	_ "silkroad-backend/docs"
	"silkroad-backend/pkg/middlewares"
	"silkroad-backend/pkg/routes"
)

// @title Swagger
// @version 2.0
// @description This is a simple server.

// @BasePath /api
// @schemes http
func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		err := godotenv.Load(".env.production")
		if err != nil {
			log.Fatalf("Error loading .env.production file: %s.", err)
		}
		middlewares.FiberMiddlewares(app, false)
		routes.APIRoutes(app)
		routes.SPARoutes(app)
	default:
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatalf("Error loading .env.development file: %s.", err)
		}
		middlewares.FiberMiddlewares(app, true)
		routes.SwaggerRoutes(app)
		routes.APIRoutes(app)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Please add a PORT field in the env file.")
		return
	}
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Error starting: %s.", err)
	}
}
