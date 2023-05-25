package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"silkroad-backend/app/models"
)

// LoginRequest 定义LoginRequest结构体类型
type LoginRequest struct {
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
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} utils.Response "{"success":true,"message":"登录成功！","result":null}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"用户名或密码错误！","result":null}"
// @Router /v1/admin/login [post]
func AdminLogin(ctx *fiber.Ctx) error {
	// 定义结构体类型
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 从请求体中读取JSON数据
	body := ctx.Body()

	// 反序列化JSON数据为LoginRequest类型的对象
	var req LoginRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		return err
	}

	// 访问LoginRequest对象的Username字段来获取username
	username := req.Username
	password := req.Username

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}

	adminName, err := getAdminName(db)
	adminPassword, err := getAdminName(db)

	if err != nil {
		// 处理错误
	}

	if username == adminName && password == adminPassword {
		return ctx.JSON(200)
	} else {
		return ctx.JSON(400)
	}
}

func getAdminName(db *gorm.DB) (string, error) {
	var setting models.Setting
	err := db.Where("key = ?", "ADMIN_NAME").First(&setting).Error
	if err != nil {
		return "", err
	}

	var name struct {
		Data string `json:"data"`
	}
	err = json.Unmarshal(setting.Value, &name)
	if err != nil {
		return "", err
	}

	return name.Data, nil
}
