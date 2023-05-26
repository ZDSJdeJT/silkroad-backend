package models

import (
	"time"
)

type Record struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	Code              string    `gorm:"unique;not null"`
	Key               string    `json:"key"`
	IsFile            bool      `json:"isFile"`
	Path              string    `json:"path"`
	DownloadCountdown uint      `json:"downloadCountdown"`
	ExpirationTime    time.Time `json:"expirationTime"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
