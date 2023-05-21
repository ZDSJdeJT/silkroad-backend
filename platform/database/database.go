package database

import (
	"gorm.io/datatypes"
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

	err = db.AutoMigrate(&models.Option{})
	if err != nil {
		return err
	}

	err = initOptions(db)
	if err != nil {
		return err
	}
	return nil
}

func initOptions(db *gorm.DB) error {
	var options []models.Option
	data := db.Find(&options)
	if data.RowsAffected == 0 {
		defaultOptions := []models.Option{
			{Key: "ADMIN_NAME", Value: datatypes.JSON(`{"data":"admin"}`), Label: "管理员名称", IsPublic: false},
			{Key: "ADMIN_PASSWORD", Value: datatypes.JSON(`{"data":"admin"}`), Label: "管理员密码", IsPublic: false},
			{Key: "WEBSITE_NAME", Value: datatypes.JSON(`{"data":"Silk Road"}`), Label: "网站名称", IsPublic: true},
			{Key: "WEBSITE_DESCRIPTION", Value: datatypes.JSON(`{"data":"This is my website"}`), Label: "网站描述", IsPublic: true},
			{Key: "WEBSITE_KEYWORDS", Value: datatypes.JSON(`{"data":"website, keywords, search"}`), Label: "网站关键词", IsPublic: true},
			{Key: "MAX_RETENTION_TIME", Value: datatypes.JSON(`{"data":365}`), Label: "最长保留时间", IsPublic: true},
			{Key: "DELETE_EXPIRED_INTERVAL", Value: datatypes.JSON(`{"data":30}`), Label: "删除过期间隔", IsPublic: false},
			{Key: "UPLOAD_SIZE_LIMIT", Value: datatypes.JSON(`{"data":10}`), Label: "上传大小限制", IsPublic: true},
			{Key: "FOOTER_INFO", Value: datatypes.JSON(`{"data":"Copyright © 2023 Silk Road"}`), Label: "页脚信息", IsPublic: true},
		}
		if err := db.CreateInBatches(defaultOptions, len(defaultOptions)).Error; err != nil {
			return err
		}
	}
	return nil
}
