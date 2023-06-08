package v1

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"os"
	"silkroad-backend/cache"
	"silkroad-backend/database"
	"silkroad-backend/i18n"
	"silkroad-backend/models"
	"silkroad-backend/utils"
	"strconv"
	"time"
)

const ChunkBytes = /* 2 MB */ 2097152

// UploadFile 文件上传接口
//
// @Summary 上传文件切片
// @Description 上传文件切片
// @Tags 记录
// @Accept json
// @Produce json
// @Param uuid path string true "uuid"
// @Param size formData string true "size"
// @Param totalChunks formData string true "totalChunks"
// @Param index formData string true "index"
// @Param chunk formData file true "chunk"
// @Success 200 {object} utils.Response "{"success":true,"message":"","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public/upload/file/{uuid} [post]
func UploadFile(ctx *fiber.Ctx) error {
	sizeStr := ctx.FormValue("size")
	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		return err
	}
	maxUploadFileBytes := cache.LoadNumberValue(models.MaxUploadFileBytes)
	if size > maxUploadFileBytes {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "uploadFileTooLarge", map[string]interface{}{
			"Max": maxUploadFileBytes / 1048576,
		}) + "MB"
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	totalChunksStr := ctx.FormValue("totalChunks")
	totalChunks, err := strconv.ParseUint(totalChunksStr, 10, 64)
	if err != nil {
		return err
	}
	indexStr := ctx.FormValue("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		return err
	}
	if index >= totalChunks {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("todo 数据有误"))
	}
	chunk, err := ctx.FormFile("chunk")
	if err != nil {
		return err
	}
	if chunk.Size > ChunkBytes {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail("todo 切片太大"))
	}
	id := ctx.Params("uuid")
	path := "./data/chunks/" + id
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	err = ctx.SaveFile(chunk, path+"/"+indexStr)
	if err != nil {
		return err
	}
	return ctx.JSON(utils.Success(nil))
}

type MergeFileForm struct {
	KeepDays      uint64 `json:"keepDays"`
	DownloadTimes uint64 `json:"downloadTimes"`
	Filename      string `json:"filename"`
}

func MergeFile(ctx *fiber.Ctx) error {
	// 从请求体中读取 JSON 数据
	body := ctx.Body()
	var req MergeFileForm
	err := json.Unmarshal(body, &req)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	id := ctx.Params("uuid")
	recordId := uuid.New()
	err = utils.MergeFiles("./data/chunks/"+id, "./data/files/"+recordId.String()+"/"+req.Filename)
	if err != nil {
		return err
	}

	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}

	// 随机生成一个接收码
	var count int64
	db.Model(&models.Record{}).Count(&count)
	code, err := utils.GenerateReceiveCode(int(count))
	if err != nil {
		return err
	}

	// 存入数据库
	record := models.Record{
		Id:            recordId,
		Code:          code,
		Content:       req.Filename,
		IsFile:        true,
		DownloadTimes: req.DownloadTimes,
		ExpireAt:      time.Now().AddDate(0, 0, int(req.KeepDays)),
	}
	db.Create(&record)

	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "uploadFileSuccess")
	return ctx.JSON(utils.SuccessWithMessage(code, msg))
}

type UploadTextForm struct {
	KeepDays      uint64 `json:"keepDays"`
	DownloadTimes uint64 `json:"downloadTimes"`
	Text          string `json:"text"`
}

// UploadText 上传文本接口
//
// @Summary 上传文本
// @Description 上传文本
// @Tags 记录
// @Accept json
// @Produce json
// @Param admin body UploadTextForm true "上传信息"
// @Success 200 {object} utils.Response "{"success":true,"message":"文本上传成功","result":"973758"}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public/upload/text [post]
func UploadText(ctx *fiber.Ctx) error {
	// 从请求体中读取 JSON 数据
	body := ctx.Body()
	var req UploadTextForm
	err := json.Unmarshal(body, &req)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxKeepDays := cache.LoadNumberValue(models.MaxKeepDays)
	if req.KeepDays > maxKeepDays || req.KeepDays < 1 {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxKeepDaysInvalid", map[string]interface{}{
			"MaxKeepDays": maxKeepDays,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxUploadTextLength := cache.LoadNumberValue(models.MaxUploadTextLength)
	textLength := uint64(len(req.Text))
	if textLength > maxUploadTextLength || textLength < 1 {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxUploadTextLengthInvalid", map[string]interface{}{
			"MaxUploadTextLength": maxUploadTextLength,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxDownloadTimes := cache.LoadNumberValue(models.MaxDownloadTimes)
	if req.DownloadTimes > maxDownloadTimes {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxDownloadTimesInvalid", map[string]interface{}{
			"MaxDownloadTimes": maxDownloadTimes,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}

	// 随机生成一个接收码
	var count int64
	db.Model(&models.Record{}).Count(&count)
	code, err := utils.GenerateReceiveCode(int(count))
	if err != nil {
		return err
	}

	// 存入数据库
	record := models.Record{
		Id:            uuid.New(),
		Code:          code,
		Content:       req.Text,
		IsFile:        false,
		DownloadTimes: req.DownloadTimes,
		ExpireAt:      time.Now().AddDate(0, 0, int(req.KeepDays)),
	}
	db.Create(&record)

	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "uploadTextSuccess")
	return ctx.JSON(utils.SuccessWithMessage(code, msg))
}

// Receive 接收接口
//
// @Summary 接收
// @Description 接收文件或文本
// @Tags 记录
// @Accept json
// @Produce json
// @Param code path string true "接收码"
// @Success 200 {object} utils.Response "{"success":true,"message":"接收成功","result":{"id":"09cb82b3-20dc-4218-bcfc-dc33ca1ddb6a","code":"579186","content":"string","isFile":false,"downloadTimes":2,"expireAt":"2023-06-08T21:05:55.2348526+08:00"}}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"接收码无效","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public/receive/{code} [get]
func Receive(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}
	var record models.Record
	if result := db.Where("code = ? AND expire_at >= ?", code, time.Now()).First(&record); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "receiveFail")
			return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
		}
		return err
	}
	if record.IsFile {
		if record.DownloadTimes != 0 {
			if record.DownloadTimes == 1 {
				if err := db.Delete(record).Error; err != nil {
					return err
				}
				// todo 删除文件
			} else {
				record.DownloadTimes--
				if err := db.Save(record).Error; err != nil {
					return err
				}
				record.DownloadTimes++
			}
		}
		file, err := os.Open("./data/files/" + record.Id.String() + "/" + record.Content)
		if err != nil {
			return err
		}

		// 设置响应头，指定 Content-Type 和 Content-Disposition
		ctx.Set("Content-Type", "application/octet-stream")
		ctx.Set("Content-Disposition", "attachment; filename=\""+record.Content+"\"")

		// 使用 JSON 编码器将 JSON 数据写入响应体中
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "receiveSuccess")
		res := utils.SuccessWithMessage(record, msg)
		if err := json.NewEncoder(ctx.Response().BodyWriter()).Encode(res); err != nil {
			return err
		}

		// 发送文件流到 ResponseWriter 中
		if _, err := io.Copy(ctx.Response().BodyWriter(), file); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}

		return nil
	} else {
		if record.DownloadTimes != 0 {
			if record.DownloadTimes == 1 {
				if err := db.Delete(record).Error; err != nil {
					return err
				}
			} else {
				record.DownloadTimes--
				if err := db.Save(record).Error; err != nil {
					return err
				}
				record.DownloadTimes++
			}
		}
	}
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "receiveSuccess")
	return ctx.JSON(utils.SuccessWithMessage(record, msg))
}

// DeleteText 删除文本接口
//
// @Summary 删除文本
// @Description 删除文本
// @Tags 记录
// @Accept json
// @Produce json
// @Param id path string true "文本 id"
// @Success 200 {object} utils.Response "{"success":true,"message":"文本删除成功","result":null}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"未找到文本","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public/text/{id} [delete]
func DeleteText(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}
	res := db.Where("is_file = false AND expire_at >= ? AND id = ?", time.Now(), id).Delete(models.Record{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteTextFail")
		return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
	}
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteTextSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}

// DeleteFile 删除文件接口
//
// @Summary 删除文件
// @Description 删除文件
// @Tags 记录
// @Accept json
// @Produce json
// @Param id path string true "文件 id"
// @Success 200 {object} utils.Response "{"success":true,"message":"文件删除成功","result":null}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"未找到文件","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/public/file/{id} [delete]
func DeleteFile(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}
	res := db.Where("is_file = true AND expire_at >= ? AND id = ?", time.Now(), id).Delete(models.Record{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteFileFail")
		return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
	}
	// todo 删除文件
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteFileSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}
