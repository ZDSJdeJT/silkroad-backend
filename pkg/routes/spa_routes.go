package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SPARoutes(app *fiber.App) {
	app.Static("/", "client")
	app.Get("*", func(c *fiber.Ctx) error {
		err := c.SendFile("client/index.html")
		if err != nil {
			return err
		}
		return nil
	})
}
