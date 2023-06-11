package v1

import (
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/i18n"
	"silkroad-backend/utils"
)

// GetSystemInfo 获取系统信息接口
//
// @Summary 获取系统信息
// @Description 获取系统应用程序名称和版本号
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"Success","result":{"appName":"Silk Road","appVersion":"1.0.0"}}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！",result:null}"
// @Router /v1/admin/system/info [get]
func GetSystemInfo(ctx *fiber.Ctx) error {
	res := utils.Success(fiber.Map{
		"appName":    utils.APPName,
		"appVersion": utils.APPVersion,
	})
	return ctx.JSON(res)
}

// GetSystemLanguages 获取系统语言接口
//
// @Summary 获取系统语言
// @Description 获取系统支持的所有语言
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"Success","result":["zh-CN","en-US"]}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！",result:null}"
// @Router /v1/public/system/languages [get]
func GetSystemLanguages(ctx *fiber.Ctx) error {
	return ctx.JSON(utils.Success(i18n.Languages))
}
