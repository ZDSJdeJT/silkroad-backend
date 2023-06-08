package v1

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"silkroad-backend/cache"
	"silkroad-backend/database"
	"silkroad-backend/i18n"
	"silkroad-backend/models"
	"silkroad-backend/utils"
	"strconv"
)

// GetPublicSettings 获取公开的配置项接口
//
// @Summary 获取公开的配置项
// @Description 获取公开的配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"","result":{"maxKeepDays":14,"maxUploadFileBytes":100000,"maxUploadTextLength":100000,"maxDownloadTimes":5}}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public/settings [get]
func GetPublicSettings(ctx *fiber.Ctx) error {
	return ctx.JSON(utils.Success(fiber.Map{
		"maxKeepDays":         cache.LoadNumberValue(models.MaxKeepDays),
		"maxUploadFileBytes":  cache.LoadNumberValue(models.MaxUploadFileBytes),
		"maxUploadTextLength": cache.LoadNumberValue(models.MaxUploadTextLength),
		"maxDownloadTimes":    cache.LoadNumberValue(models.MaxDownloadTimes),
	}))
}

// GetSettings 获取所有配置项接口
//
// @Summary 获取所有配置项
// @Description 获取所有配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"","result":[{"key":"ADMIN_NAME","textValue":"admin","numberValue":0,"isText":true,"min":5,"max":20,"label":{"en-US":"Admin name","zh-CN":"管理员名称"}},{...}]}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/settings [get]
func GetSettings(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}

	// 查找所有配置项
	var settings []models.Setting
	if err := db.Find(&settings).Error; err != nil {
		return err
	}

	// 管理员密码需置空
	for i := range settings {
		if settings[i].Key == "ADMIN_PASSWORD" {
			settings[i].TextValue = ""
			break
		}
	}

	return ctx.JSON(utils.Success(settings))
}

// UpdateSetting 更新配置项接口
//
// @Summary 更新配置项
// @Description 更新配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Param key path string true "配置项键"
// @Param value body string true "配置项新值"
// @Success 200 {object} utils.Response "{"success":true,"message":"更新成功","result":null}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"未找到配置项","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/settings/{key} [put]
func UpdateSetting(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}

	// 从请求路径中获取要更新的配置项键
	key := ctx.Params("key")

	// 检查配置项是否存在
	var setting models.Setting
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "updateSettingFail")
			return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
		}
		return err
	}

	body := string(ctx.Body())

	if setting.IsText {
		length := len(body)
		if length < int(setting.Min) || length > int(setting.Max) {
			msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "textSettingInvalid", map[string]interface{}{
				"Min": setting.Min,
				"Max": setting.Max,
			})
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
		}
		if key == models.AdminPassword {
			setting.TextValue, err = utils.EncryptPassword(body)
			if err != nil {
				return err
			}
		} else {
			setting.TextValue = body
		}
	} else {
		value, err := strconv.ParseUint(body, 10, 64)
		if err != nil {
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
		}
		if value < setting.Min || value > setting.Max {
			msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "numberSettingInvalid", map[string]interface{}{
				"Min": setting.Min,
				"Max": setting.Max,
			})
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
		}
		setting.NumberValue = value
	}

	// 保存更改
	if err := db.Save(&setting).Error; err != nil {
		return err
	}

	err = cache.HandleSetting(setting)
	if err != nil {
		return err
	}
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "updateSettingSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}
