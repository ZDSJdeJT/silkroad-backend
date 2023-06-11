package v1

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"os"
	"silkroad-backend/cache"
	"silkroad-backend/database"
	"silkroad-backend/i18n"
	"silkroad-backend/models"
	"silkroad-backend/utils"
	"strconv"
	"time"
)

// UploadFile 文件上传接口
//
// @Summary 上传文件切片
// @Description 上传文件切片
// @Tags 记录
// @Accept json
// @Produce json
// @Param uuid path string true "uuid"
// @Param size formData string true "文件总大小"
// @Param total formData string true "文件总切片数"
// @Param index formData string true "文件切片索引"
// @Param chunk formData file true "文件切片"
// @Success 200 {object} utils.Response "{"success":true,"message":"","result":null}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/upload/files/{uuid} [post]
func UploadFile(ctx *fiber.Ctx) error {
	sizeStr := ctx.FormValue("size")
	size, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	uploadFileBytes := cache.LoadNumberValue(models.UploadFileBytes)
	if size > uploadFileBytes {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "uploadFileTooLarge", map[string]interface{}{
			"Max": uploadFileBytes / 1048576,
		}) + "MB"
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	totalStr := ctx.FormValue("total")
	total, err := strconv.ParseUint(totalStr, 10, 64)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	indexStr := ctx.FormValue("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	if index >= total {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	uploadChunkBytes := cache.LoadNumberValue(models.UploadChunkBytes)
	if uint64(math.Ceil(float64(uploadFileBytes)/float64(uploadChunkBytes))) < total {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "uploadFileTooLarge", map[string]interface{}{
			"Max": uploadFileBytes / 1048576,
		}) + "MB"
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	chunk, err := ctx.FormFile("chunk")
	if err != nil {
		return err
	}
	chunkSize := uint64(chunk.Size)
	if chunkSize > uploadChunkBytes || index*uploadChunkBytes+chunkSize > uploadFileBytes {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "uploadFileTooLarge", map[string]interface{}{
			"Max": uploadFileBytes / 1048576,
		}) + "MB"
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}
	id := ctx.Params("uuid")
	path := database.ChunksDir + id
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

// MergeFile 合并文件接口
//
// @Summary 合并文件
// @Description 将文件切片合并
// @Tags 记录
// @Accept json
// @Produce json
// @Param uuid path string true "uuid"
// @Param merge body MergeFileForm true "合并信息"
// @Success 200 {object} utils.Response "{"success":true,"message":"文件上传成功","result":null}"
// @Failure 400 {object} utils.Response "{"success":false,"message":"请求无效或参数错误","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/upload/files/merge/{uuid} [post]
func MergeFile(ctx *fiber.Ctx) error {
	// 从请求体中读取 JSON 数据
	body := ctx.Body()
	var req MergeFileForm
	err := json.Unmarshal(body, &req)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxKeepDays := cache.LoadNumberValue(models.KeepDays)
	if req.KeepDays > maxKeepDays || req.KeepDays < 1 {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxKeepDaysInvalid", map[string]interface{}{
			"KeepDays": maxKeepDays,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxDownloadTimes := cache.LoadNumberValue(models.DownloadTimes)
	if req.DownloadTimes > maxDownloadTimes {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxDownloadTimesInvalid", map[string]interface{}{
			"DownloadTimes": maxDownloadTimes,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	if req.Filename == "" {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	id := ctx.Params("uuid")
	recordId := uuid.New()
	err = utils.MergeFiles(database.ChunksDir+id, database.DataDir+recordId.String()+"/"+req.Filename)
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
	err = db.Model(&models.Record{}).Count(&count).Error
	if err != nil {
		return err
	}
	code, err := utils.GenerateReceiveCode(int(count))
	if err != nil {
		return err
	}

	// 存入数据库
	record := models.Record{
		Id:            recordId,
		Code:          code,
		Filename:      req.Filename,
		DownloadTimes: req.DownloadTimes,
		ExpireAt:      time.Now().AddDate(0, 0, int(req.KeepDays)),
	}
	err = db.Create(&record).Error
	if err != nil {
		return err
	}

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
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/upload/texts [post]
func UploadText(ctx *fiber.Ctx) error {
	// 从请求体中读取 JSON 数据
	body := ctx.Body()
	var req UploadTextForm
	err := json.Unmarshal(body, &req)
	if err != nil {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "badRequest")
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxKeepDays := cache.LoadNumberValue(models.KeepDays)
	if req.KeepDays > maxKeepDays || req.KeepDays < 1 {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxKeepDaysInvalid", map[string]interface{}{
			"KeepDays": maxKeepDays,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxUploadTextLength := cache.LoadNumberValue(models.UploadTextLength)
	textLength := uint64(len(req.Text))
	if textLength > maxUploadTextLength || textLength < 1 {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxUploadTextLengthInvalid", map[string]interface{}{
			"UploadTextLength": maxUploadTextLength,
		})
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.Fail(msg))
	}

	maxDownloadTimes := cache.LoadNumberValue(models.DownloadTimes)
	if req.DownloadTimes > maxDownloadTimes {
		msg := i18n.GetLocalizedMessageWithTemplate(ctx.Locals("lang").(string), "maxDownloadTimesInvalid", map[string]interface{}{
			"DownloadTimes": maxDownloadTimes,
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
	err = db.Model(&models.Record{}).Count(&count).Error
	if err != nil {
		return err
	}
	code, err := utils.GenerateReceiveCode(int(count))
	if err != nil {
		return err
	}

	id := uuid.New()
	// 存入数据库
	record := models.Record{
		Id:            id,
		Code:          code,
		DownloadTimes: req.DownloadTimes,
		ExpireAt:      time.Now().AddDate(0, 0, int(req.KeepDays)),
	}
	err = db.Create(&record).Error
	if err != nil {
		return err
	}

	err = utils.WriteToFile(database.DataDir+id.String()+"/"+database.TextFilename, req.Text)
	if err != nil {
		return err
	}

	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "uploadTextSuccess")
	return ctx.JSON(utils.SuccessWithMessage(code, msg))
}

// DeleteRecord 删除记录接口
//
// @Summary 删除记录
// @Description 删除记录
// @Tags 记录
// @Accept json
// @Produce json
// @Param id path string true "记录 id"
// @Success 200 {object} utils.Response "{"success":true,"message":"记录删除成功","result":null}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"未找到记录","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/records/{id} [delete]
func DeleteRecord(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}
	res := db.Where("expire_at >= ? AND id = ?", time.Now(), id).Delete(models.Record{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteRecordFail")
		return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
	}
	err = os.RemoveAll(database.DataDir + id)
	if err != nil {
		return err
	}
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteRecordSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}

// DeleteExpiredTextRecords 删除过期文本接口
//
// @Summary 删除过期文本
// @Description 删除过期文本
// @Tags 记录
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"删除过期文本成功","result":null}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/records/expired/text [delete]
func DeleteExpiredTextRecords(ctx *fiber.Ctx) error {
	database.DeleteExpiredTextRecords()
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteExpiredTextsRecordsSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}

// DeleteExpiredFileRecords 删除过期文件接口
//
// @Summary 删除过期文件
// @Description 删除过期文件
// @Tags 记录
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"删除过期文件成功","result":null}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/records/expired/file [delete]
func DeleteExpiredFileRecords(ctx *fiber.Ctx) error {
	database.DeleteExpiredFileRecords()
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteExpiredFilesRecordsSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}

// DeleteExpiredChunks 删除过期文件切片接口
//
// @Summary 删除过期文件切片
// @Description 删除过期文件切片
// @Tags 记录
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "{"success":true,"message":"删除过期文件切片成功","result":null}"
// @Failure 401 {object} utils.Response "{"success":false,"message":"请登录后再试",result:null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Router /v1/admin/expired/chunks [delete]
func DeleteExpiredChunks(ctx *fiber.Ctx) error {
	utils.DeleteExpiredChunks()
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "deleteOldChunksSuccess")
	return ctx.JSON(utils.SuccessWithMessage(nil, msg))
}

// GetRecordByCode 获取记录接口
//
// @Summary 获取记录
// @Description 根据接收码获取记录
// @Tags 记录
// @Accept json
// @Produce json
// @Param code path string true "接收码"
// @Success 200 {object} utils.Response "{"success":true,"message":"","result":{"success":true,"message":"","result":{"id":"c04ff62e-49ae-4320-9f7f-7ad8582235f4","code":"045151","filename":"","downloadTimes":1,"expireAt":"2023-06-12T00:01:28.2012091+08:00"}}}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"接收码无效","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/records/{code} [get]
func GetRecordByCode(ctx *fiber.Ctx) error {
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
	return ctx.JSON(utils.Success(record))
}

// ReceiveText 接收文本接口
//
// @Summary 接收文本
// @Description 接收文本
// @Tags 记录
// @Accept json
// @Produce json
// @Param code path string true "接收码"
// @Success 200 {object} utils.Response "{"success":true,"message":"","result":"text"}"
// @Failure 404 {object} utils.Response "{"success":false,"message":"接收码无效","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/receive/texts/{code} [get]
func ReceiveText(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}
	now := time.Now()
	var record models.Record
	if result := db.Where("code = ? AND expire_at >= ? AND filename = \"\"", code, now).First(&record); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "receiveFail")
			return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
		}
		return err
	}

	if record.DownloadTimes != 0 {
		if record.DownloadTimes == 1 {
			// 逻辑删除
			// 定时任务会进行物理删除
			if err := db.Model(&record).Select("expire_at").Update("expire_at", now).Error; err != nil {
				return err
			}
		} else {
			if err := db.Model(&record).Select("download_times").Update("download_times", record.DownloadTimes-1).Error; err != nil {
				return err
			}
		}
	}

	text, err := os.ReadFile(database.DataDir + record.Id.String() + "/" + database.TextFilename)
	msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "receiveSuccess")
	return ctx.JSON(utils.SuccessWithMessage(string(text), msg))
}

// ReceiveFile 接收文件接口
//
// @Summary 接收文件
// @Description 接收文件
// @Tags 记录
// @Accept json
// @Produce json
// @Param code path string true "接收码"
// @Success 200
// @Failure 404 {object} utils.Response "{"success":false,"message":"接收码无效","result":null}"
// @Failure 429 {object} utils.Response "{"success":false,"message":"请求过于频繁，请稍后再试！","result":null}"
// @Failure 500 {object} utils.Response "{"success":false,"message":"服务器错误","result":null}"
// @Router /v1/public/receive/files/{code} [get]
func ReceiveFile(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	// 打开数据库连接
	db, err := database.OpenDBConnection()
	if err != nil {
		return err
	}
	now := time.Now()
	var record models.Record
	if result := db.Where("code = ? AND expire_at >= ?", code, now).First(&record); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			msg := i18n.GetLocalizedMessage(ctx.Locals("lang").(string), "receiveFail")
			return ctx.Status(fiber.StatusNotFound).JSON(utils.Fail(msg))
		}
		return err
	}

	if record.DownloadTimes != 0 {
		if record.DownloadTimes == 1 {
			// 逻辑删除
			// 定时任务会进行物理删除
			if err := db.Model(&record).Select("expire_at").Update("expire_at", now).Error; err != nil {
				return err
			}
		} else {
			if err := db.Model(&record).Select("download_times").Update("download_times", record.DownloadTimes-1).Error; err != nil {
				return err
			}
		}
	}
	ctx.Set(fiber.HeaderContentType, "application/octet-stream")
	ctx.Set(fiber.HeaderContentDisposition, "attachment; filename=\""+record.Filename+"\"")
	return ctx.SendFile(database.DataDir + record.Id.String() + "/" + record.Filename)
}
