package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"silkroad-backend/app/models"
	"silkroad-backend/pkg/utils"
	"strconv"
)

// GetPublicSettings 获取公开的配置项接口
//
// @Summary 获取公开的配置项
// @Description 获取公开的配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"成功","result":{"UPLOAD_FILE_SIZE_LIMIT":10, ...}}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public-settings [get]
func GetPublicSettings(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
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
		var memo map[string]interface{}
		err = json.Unmarshal(setting.Value, &memo)
		temp, _ := memo["data"]
		result[setting.Key] = temp
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
// @Success 200 {object} utils.Response "{"success":true,"message":"成功","result":[{"key":"UPLOAD_FILE_SIZE_LIMIT","value":{"data":10},"label":"上传大小限制","isPublic":true,"createdAt":"2023-05-22T15:10:40.7958637+08:00","updatedAt":"2023-05-22T15:10:40.7958637+08:00"}, {...}]}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/settings [get]
func GetSettings(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}

	// 查找所有公共配置项
	var publicSettings []models.Setting
	if err := db.Find(&publicSettings).Error; err != nil {
		return err
	}

	// 返回 JSON 格式响应
	res := utils.Success(publicSettings)
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
// @Success 200 {object} utils.Response "{"success":true,"message":"更新成功","result":{"key":"WEBSITE_TITLE","value":{"data":"New Silk Road"},"label":"网站名称","isPublic":true,"CreatedAt":"2023-05-21T15:29:42.6390127+08:00","UpdatedAt":"2023-05-23T10:31:12.1234567+08:00"}}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"未找到配置项","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/settings/{key} [put]
func UpdateSetting(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
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

	// 根据原配置项值的类型进行不同的操作
	var memo map[string]interface{}
	err = json.Unmarshal(setting.Value, &memo)
	if err != nil {
		return err
	}
	temp, _ := memo["data"]
	switch temp.(type) {
	case string:
		// "data" 是字符串类型，解析字符串并存储
		data := struct {
			Data string `json:"data"`
		}{
			Data: body,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		setting.Value = jsonData
	case float64:
		// "data" 是浮点数类型，解析数字并存储
		value, err := strconv.ParseFloat(body, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("请求无效或参数错误"))
		}
		data := struct {
			Data float64 `json:"data"`
		}{
			Data: value,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		setting.Value = jsonData
	}

	// 保存更改
	if err := db.Save(&setting).Error; err != nil {
		return err
	}

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

	// 返回 JSON 格式响应
	res := utils.SuccessWithMessage(setting, "更新成功")
	return ctx.JSON(res)
}
