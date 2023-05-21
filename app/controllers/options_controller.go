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
)

// GetPublicOptions 获取公开的配置项接口
//
// @Summary 获取公开的配置项
// @Description 获取公开的配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"成功","result":[{"key":"WEBSITE_NAME","value":{"data":"Silk Road"},"label":"网站名称","isPublic":true,"CreatedAt":"2023-05-21T15:29:42.6390127+08:00","UpdatedAt":"2023-05-21T15:29:42.6390127+08:00"}, {...}]}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public-options [get]
func GetPublicOptions(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}

	// 查找所有公共配置项
	var publicOptions []models.Option
	if err := db.Where("is_public = ?", true).Find(&publicOptions).Error; err != nil {
		return err
	}

	// 返回 JSON 格式响应
	res := utils.Success(publicOptions)
	return ctx.JSON(res)
}

// UpdateOption 更新配置项接口
//
// @Summary 更新配置项
// @Description 更新配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Param option body object true "配置项"
// @Success 200 {object} utils.Response "{"success":true,"message":"成功","result":{"key":"WEBSITE_NAME","value":{"data":"New Silk Road"},"label":"网站名称","isPublic":true,"CreatedAt":"2023-05-21T15:29:42.6390127+08:00","UpdatedAt":"2023-05-23T10:31:12.1234567+08:00"}}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"未找到配置项","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/options [put]
func UpdateOption(ctx *fiber.Ctx) error {
	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}

	// 解析请求体中的 JSON 数据
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"Value" binding:"required"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("请求无效或参数错误"))
	}

	// 获取要更新的配置项键
	id := req.Key

	// 检查配置项是否存在
	var option models.Option
	if err := db.Where("key = ?", id).First(&option).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail("未找到配置项"))
		}
		return err
	}

	// 更新配置项值
	data := struct {
		Data string `json:"data"`
	}{
		Data: req.Value,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		// 处理错误
	}
	option.Value = jsonData

	// 保存更改
	if err := db.Save(&option).Error; err != nil {
		return err
	}

	// 返回 JSON 格式响应
	res := utils.Success(option)
	return ctx.JSON(res)
}
