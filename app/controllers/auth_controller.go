package controllers

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/pkg/utils"
)

func AdminLogin(ctx *fiber.Ctx) error {
	// todo
	return ctx.JSON(utils.Success(nil))
}
