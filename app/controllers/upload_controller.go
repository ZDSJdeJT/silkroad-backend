package controllers

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/pkg/utils"
)

func UploadFile(ctx *fiber.Ctx) error {
	res := utils.Success("todo")
	return ctx.JSON(res)
}

func MergeFile(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	res := utils.SuccessWithMessage(id, "todo")
	return ctx.JSON(res)
}

func UploadText(ctx *fiber.Ctx) error {
	res := utils.Success("todo")
	return ctx.JSON(res)
}
