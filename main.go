package main

import (
	_ "backend/docs"
	"encoding/json"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

// @title Swagger
// @version 2.0
// @description This is a simple server.

// @BasePath /
// @schemes http
func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	env := os.Getenv("APP_ENV")
	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	case "testing":
		envFile = ".env.testing"
		app.Use(cors.New())
		app.Get("/swagger/*", swagger.HandlerDefault)
	default:
		envFile = ".env.development"
		app.Use(cors.New())
		app.Get("/swagger/*", swagger.HandlerDefault)
	}
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %s", envFile, err)
	}

	app.Use(compress.New())
	app.Use(csrf.New())
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(favicon.New(favicon.Config{
		File: "./static/favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Get("/", Hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e := app.Listen(":" + port)
	if e != nil {
		log.Fatalf("Error starting: %s", e)
		return
	}
}

// Hello godoc
// @Summary Test function.
// @Description A test function.
// @Tags test
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func Hello(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "silkroad",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
