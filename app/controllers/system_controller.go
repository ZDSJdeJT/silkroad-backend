package controllers

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"silkroad-backend/pkg/utils"
)

// GetSystemInfo 获取系统信息接口
//
// @Summary 获取系统信息
// @Description 获取系统应用程序名称和版本号
// @Tags 系统信息
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"Success","result":{"appName":"GoApp","appVersion":"1.0.0"}}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！",result:null}"
// @Router /v1/system/info [get]
func GetSystemInfo(ctx *fiber.Ctx) error {
	data := map[string]interface{}{
		"appName":    os.Getenv("APP_NAME"),
		"appVersion": os.Getenv("APP_VERSION"),
	}
	res := utils.Success(data)
	return ctx.JSON(res)
}
