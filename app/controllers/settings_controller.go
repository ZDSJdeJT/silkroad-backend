package controllers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"silkroad-backend/app/models"
	"silkroad-backend/pkg/utils"
	"silkroad-backend/platform/database"
	"strconv"
)

// GetPublicSettings 获取公开的配置项接口
//
// @Summary 获取公开的配置项
// @Description 获取公开的配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"成功","result":{"UPLOAD_FILE_SIZE_LIMIT":10,...}}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public-settings [get]
func GetPublicSettings(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}

	// 查找所有公共配置项
	var publicSettings []models.Setting
	if err := db.Where("is_public = ?", true).Find(&publicSettings).Error; err != nil {
		return err
	}

	// 将查询结果转换为 key-value 形式
	result := make(map[string]interface{})
	for _, setting := range publicSettings {
		if setting.IsText {
			result[setting.Key] = setting.TextValue
		} else {
			result[setting.Key] = setting.NumberValue
		}
	}

	// 返回 JSON 格式响应
	res := utils.Success(result)
	return ctx.JSON(res)
}

// GetSettings 获取所有配置项接口
//
// @Summary 获取所有配置项
// @Description 获取所有配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"成功","result":[{"key":"ADMIN_NAME","textValue":"admin","numberValue":0,"isText":true,"min":5,"max":16,"label":"管理员名称","isPublic":false,"createdAt":"2023-05-28T12:33:15.8278992+08:00","updatedAt":"2023-05-28T12:33:15.8278992+08:00"},{...}]}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/settings [get]
// @Security ApiKeyAuth
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

	for i := range settings {
		if settings[i].Key == "ADMIN_PASSWORD" {
			settings[i].TextValue = ""
			break
		}
	}

	// 返回 JSON 格式响应
	res := utils.Success(settings)
	return ctx.JSON(res)
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
// @Security ApiKeyAuth
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
			return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail("未找到配置项"))
		}
		return err
	}

	body := string(ctx.Body())

	if setting.IsText {
		length := len(body)
		if length < int(setting.Min) || length > int(setting.Max) {
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(fmt.Sprintf("长度必须在 %d-%d 之间", setting.Min, setting.Max)))
		}
		if key == "ADMIN_PASSWORD" {
			var adminPassword, err = utils.EncryptPassword(body)
			if err != nil {
				return err
			}
			setting.TextValue = adminPassword
		} else {
			setting.TextValue = body
			switch key {
			case "WEBSITE_TITLE":
				err := utils.ReplaceClientHTMLTitle(body)
				if err != nil {
					return err
				}
			case "WEBSITE_DESCRIPTION":
				err := utils.ReplaceClientHTMLMetaDescription(body)
				if err != nil {
					return err
				}
			case "WEBSITE_KEYWORDS":
				err := utils.ReplaceClientHTMLMetaKeywords(body)
				if err != nil {
					return err
				}
			}
		}
	} else {
		value, err := strconv.ParseUint(body, 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("请求无效或参数错误"))
		}
		setting.NumberValue = uint(value)
		if setting.NumberValue < setting.Min || setting.NumberValue > setting.Max {
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(fmt.Sprintf("值必须在 %d-%d 之间", setting.Min, setting.Max)))
		}
	}

	// 保存更改
	if err := db.Save(&setting).Error; err != nil {
		return err
	}

	// 返回 JSON 格式响应
	res := utils.SuccessWithMessage(nil, "更新成功")
	return ctx.JSON(res)
}
