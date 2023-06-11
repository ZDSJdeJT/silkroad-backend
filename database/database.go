package database

import (
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"silkroad-backend/models"
	"silkroad-backend/utils"
)

func InitDatabase() ([]models.Setting, error) {
	db, err := OpenDBConnection()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Setting{}, &models.Record{})
	if err != nil {
		return nil, err
	}

	settings, err := initSettings(db)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func initSettings(db *gorm.DB) ([]models.Setting, error) {
	var settings []models.Setting
	data := db.Find(&settings)
	if data.Error != nil {
		return nil, data.Error
	}
	if data.RowsAffected == 0 {
		var adminPassword, err = utils.EncryptPassword("admin")
		if err != nil {
			return nil, err
		}
		settings = []models.Setting{
			{Key: models.AdminName, IsText: true, TextValue: "admin", Min: 5, Max: 20, Label: json.RawMessage(`{"en-US":"Admin name","zh-CN":"管理员名称"}`)},
			{Key: models.AdminPassword, IsText: true, TextValue: adminPassword, Min: 5, Max: 20, Label: json.RawMessage(`{"en-US":"Admin password","zh-CN":"管理员密码"}`)},
			{Key: models.WebsiteDescription, IsText: true, TextValue: "File Express Cabinet - It allows for anonymous sharing of text and files, similar to using express delivery to quickly send files.", Min: 0, Max: 160, Label: json.RawMessage(`{"en-US":"Website description","zh-CN":"网站描述"}`)},
			{Key: models.WebsiteKeywords, IsText: true, TextValue: "File Express Cabinet, anonymous sharing, text files, express delivery", Min: 0, Max: 160, Label: json.RawMessage(`{"en-US":"Website keywords","zh-CN":"网站关键词"}`)},
			{Key: models.KeepDays, IsText: false, NumberValue: 7, Min: 1, Max: 7, Label: json.RawMessage(`{"en-US":"Max keep days","zh-CN":"最长保留天数"}`)},
			{Key: models.DownloadTimes, IsText: false, NumberValue: 5, Min: 1, Max: 10, Label: json.RawMessage(`{"en-US":"Max download times","zh-CN":"最大下载次数"}`)},
			{Key: models.UploadFileBytes, IsText: false, NumberValue: /* 1 GB */ 1073741824, Min: 0, Max: /* 1 GB */ 1073741824, Label: json.RawMessage(`{"en-US":"Max upload file bytes","zh-CN":"最大上传文件字节数"}`)},
			{Key: models.UploadChunkBytes, IsText: false, NumberValue: /* 2 MB */ 2097152, Min: /* 1 MB */ 1048576, Max: /* 3 MB */ 3145728, Label: json.RawMessage(`{"en-US":"Max upload chunks bytes","zh-CN":"最大上传文件切片字节数"}`)},
			{Key: models.UploadTextLength, IsText: false, NumberValue: 1000000, Min: 1, Max: 1000000, Label: json.RawMessage(`{"en-US":"Max upload text length","zh-CN":"最大上传文本长度"}`)},
		}
		if err := db.CreateInBatches(settings, len(settings)).Error; err != nil {
			return nil, err
		}
	}
	return settings, nil
}

func OpenDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(utils.DatabaseDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
