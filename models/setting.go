package models

import (
	"encoding/json"
	"time"
)

type Label struct {
	EnUS string `json:"en-US"`
	ZhCN string `json:"zh-CN"`
}

type Setting struct {
	Key         string          `json:"key" gorm:"primaryKey"`
	TextValue   string          `json:"textValue"`
	NumberValue uint64          `json:"numberValue"`
	IsText      bool            `json:"isText"`
	Min         uint64          `json:"min"`
	Max         uint64          `json:"max"`
	Label       json.RawMessage `json:"label" gorm:"type:json"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
}

const (
	AdminName          = "ADMIN_NAME"
	AdminPassword      = "ADMIN_PASSWORD"
	WebsiteDescription = "WEBSITE_DESCRIPTION"
	WebsiteKeywords    = "WEBSITE_KEYWORDS"
	KeepDays           = "KEEP_DAYS"
	DownloadTimes      = "DOWNLOAD_TIMES"
	UploadFileBytes    = "UPLOAD_FILE_BYTES"
	UploadChunkBytes   = "UPLOAD_Chunk_BYTES"
	UploadTextLength   = "UPLOAD_TEXT_LENGTH"
)
