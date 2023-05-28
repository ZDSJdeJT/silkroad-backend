package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"silkroad-backend/app/models"
	"silkroad-backend/pkg/utils"
)

func InitDatabase() error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.Setting{}, &models.Record{})
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
		var adminPassword, err = utils.EncryptPassword("admin")
		if err != nil {
			return err
		}
		defaultOptions := []models.Setting{
			{Key: "ADMIN_NAME", IsText: true, TextValue: "admin", Min: 5, Max: 16, Label: "管理员名称", IsPublic: false},
			{Key: "ADMIN_PASSWORD", IsText: true, TextValue: adminPassword, Min: 5, Max: 16, Label: "管理员密码", IsPublic: false},
			{Key: "WEBSITE_TITLE", IsText: true, TextValue: "Silk Road", Min: 5, Max: 16, Label: "网站名称", IsPublic: false},
			{Key: "WEBSITE_DESCRIPTION", IsText: true, TextValue: "Silk Road", Min: 5, Max: 16, Label: "网站描述", IsPublic: false},
			{Key: "WEBSITE_KEYWORDS", IsText: true, TextValue: "匿名口令分享文本、文件", Min: 5, Max: 16, Label: "网站关键词", IsPublic: false},
			{Key: "FOOTER_INFO", IsText: true, TextValue: `Built with <a href="https://github.com/ZDSJdeJT/silkroad-backend" target="_blank">SilkRoad</a>`, Min: 5, Max: 16, Label: "页脚信息", IsPublic: true},
			{Key: "MAX_RETENTION_TIME", IsText: false, NumberValue: 365, Min: 5, Max: 16, Label: "最长保留时间", IsPublic: false},
			{Key: "UPLOAD_FILE_SIZE_LIMIT", IsText: false, NumberValue: 100000, Min: 5, Max: 16, Label: "上传文件大小限制", IsPublic: true},
			{Key: "UPLOAD_TEXT_LENGTH_LIMIT", IsText: false, NumberValue: 1000, Min: 5, Max: 16, Label: "上传文本长度限制", IsPublic: true},
		}
		if err := db.CreateInBatches(defaultOptions, len(defaultOptions)).Error; err != nil {
			return err
		}
	}
	return nil
}

func OpenDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
