package controllers

import "github.com/gofiber/fiber/v2"

// Hello godoc
// @Summary Test function.
// @Description A test function.
// @Tags test
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/hello [get]
func Hello(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "silkroad",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
