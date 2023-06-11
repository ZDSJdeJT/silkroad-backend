package v1

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"os"
	"silkroad-backend/cache"
	"silkroad-backend/i18n"
	"silkroad-backend/models"
	"silkroad-backend/utils"
	"time"
)

type LoginForm struct {
	AdminName string `json:"adminName"`
	Password  string `json:"password"`
}

// AdminLogin 管理员登录接口
//
// @Summary 管理员登录
// @Description 管理员使用用户名和密码进行登录
// @Tags 管理员
// @Accept json
// @Produce json
// @Param admin body LoginForm true "管理员"
// @Success 200 {object} utils.Response "{"success":true,"message":"登录成功","result":null}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/admin/login [post]
func AdminLogin(ctx *fiber.Ctx) error {
	// 从请求体中读取 JSON 数据
	body := ctx.Body()

	// 反序列化 JSON 数据为 LoginRequest 类型的对象
	var req LoginForm
	err := json.Unmarshal(body, &req)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	adminName := cache.LoadTextValue(models.AdminName)
	adminPassword := cache.LoadTextValue(models.AdminPassword)

	if req.AdminName == adminName {
		err := utils.MatchPassword(req.Password, adminPassword)
		if err == nil {
			expireMinutes, err := time.ParseDuration(os.Getenv(utils.JWTExpireMinutes) + "m")
			if err != nil {
				return err
			}
			token, err := utils.GenerateJWT(adminName, expireMinutes)
			if err != nil {
				return err
			}
			cookie := fiber.Cookie{
				Name:     "JWT",
				Value:    token,
				HTTPOnly: true,
				Expires:  time.Now().Add(expireMinutes),
			}
			ctx.Cookie(&cookie)
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "loginSuccess")
			return ctx.JSON(utils.SuccessWithMessage(nil, msg))
		}
	}
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "loginFail")
	return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
}

// AdminLogout 管理员退出登录接口
//
// @Summary 管理员退出登录
// @Description 管理员退出登录
// @Tags 管理员
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"退出登录成功","result":null}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/admin/logout [post]
func AdminLogout(ctx *fiber.Ctx) error {
	expireMinutes, err := time.ParseDuration(os.Getenv(utils.JWTExpireMinutes) + "m")
	if err != nil {
		return err
	}
	ctx.Cookie(&fiber.Cookie{
		Name:    "JWT",
		Value:   "",
		Expires: time.Now().Add(-expireMinutes),
	})
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "logoutSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}
