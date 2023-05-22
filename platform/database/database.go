package database

import (
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"silkroad-backend/app/models"
)

func InitDatabase() error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Setting{})
	if err != nil {
		return err
	}

	err = initSettings(db)
	if err != nil {
		return err
	}
	return nil
}

func initSettings(db *gorm.DB) error {
	var settings []models.Setting
	data := db.Find(&settings)
	if data.RowsAffected == 0 {
		defaultOptions := []models.Setting{
			{Key: "ADMIN_NAME", Value: json.RawMessage(`{"data":"admin"}`), Label: "管理员名称", IsPublic: false},
			{Key: "ADMIN_PASSWORD", Value: json.RawMessage(`{"data":""}`), Label: "管理员密码", IsPublic: false},
			{Key: "WEBSITE_TITLE", Value: json.RawMessage(`{"data":"Silk Road"}`), Label: "网站名称", IsPublic: false},
			{Key: "WEBSITE_DESCRIPTION", Value: json.RawMessage(`{"data":"Silk Road"}`), Label: "网站描述", IsPublic: false},
			{Key: "WEBSITE_KEYWORDS", Value: json.RawMessage(`{"data":"匿名口令分享文本、文件"}`), Label: "网站关键词", IsPublic: false},
			{Key: "MAX_RETENTION_TIME", Value: json.RawMessage(`{"data":365}`), Label: "最长保留时间", IsPublic: false},
			{Key: "DELETE_EXPIRED_INTERVAL", Value: json.RawMessage(`{"data":30}`), Label: "删除过期间隔", IsPublic: false},
			{Key: "UPLOAD_FILE_SIZE_LIMIT", Value: json.RawMessage(`{"data":10}`), Label: "上传文件大小限制", IsPublic: true},
			{Key: "UPLOAD_TEXT_LENGTH_LIMIT", Value: json.RawMessage(`{"data":1000}`), Label: "上传文本长度限制", IsPublic: true},
			{Key: "FOOTER_INFO", Value: json.RawMessage(`{"data":"Copyright © 2023 Silk Road"}`), Label: "页脚信息", IsPublic: true},
		}
		if err := db.CreateInBatches(defaultOptions, len(defaultOptions)).Error; err != nil {
			return err
		}
	}
	return nil
}
