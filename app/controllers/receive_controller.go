package controllers

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/pkg/utils"
)

func Receive(ctx *fiber.Ctx) error {
	res := utils.Success("todo")
	return ctx.JSON(res)
}
