package models

import (
	"github.com/google/uuid"
	"time"
)

type Record struct {
	Id            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Code          string    `json:"code" gorm:"unique;not null"`
	Filename      string    `json:"filename"`
	DownloadTimes uint64    `json:"downloadTimes"`
	ExpireAt      time.Time `json:"expireAt"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}
