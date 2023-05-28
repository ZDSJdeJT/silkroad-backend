package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"silkroad-backend/app/models"
	"silkroad-backend/pkg/utils"
	"silkroad-backend/platform/database"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminLogin 管理员登录接口
//
// @Summary 管理员登录
// @Description 管理员使用用户名和密码进行登录
// @Tags 管理员
// @Accept json
// @Produce json
// @Param admin body LoginForm true "管理员"
// @Success 200 {object} utils.Response "{"success":true,"message":"登录成功","result":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODUxODEwNjV9.Uj37YBTlIm4v5dcqEI4371oqNuyk632BYcuqZgYSFL8"}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"用户名或密码错误","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/login [post]
func AdminLogin(ctx *fiber.Ctx) error {
	// 从请求体中读取 JSON 数据
	body := ctx.Body()

	// 反序列化 JSON 数据为 LoginRequest 类型的对象
	var req LoginForm
	err := json.Unmarshal(body, &req)
	if err != nil {
		return err
	}

	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}

	var settings []models.Setting
	if err := db.Where("key IN (?)", []string{"ADMIN_NAME", "ADMIN_PASSWORD"}).Find(&settings).Error; err != nil {
		return err
	}

	var adminName, adminPassword string
	for _, setting := range settings {
		switch setting.Key {
		case "ADMIN_NAME":
			length := len(req.Username)
			if length < int(setting.Min) || length > int(setting.Max) {
				return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("用户名或密码错误"))
			}
			adminName = setting.TextValue
		case "ADMIN_PASSWORD":
			length := len(req.Password)
			if length < int(setting.Min) || length > int(setting.Max) {
				return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("用户名或密码错误"))
			}
			adminPassword = setting.TextValue
		}
	}

	if req.Username == adminName {
		err := utils.MatchPassword(req.Password, adminPassword)
		if err == nil {
			var token string
			token, err = utils.GenerateNewAccessToken()
			if err != nil {
				return err
			}
			return ctx.JSON(utils.SuccessWithMessage(token, "登录成功"))
		}
	}
	return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("用户名或密码错误"))
}
