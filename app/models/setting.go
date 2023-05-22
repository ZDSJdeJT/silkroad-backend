package models

import (
	"encoding/json"
	"time"
)

type Setting struct {
	Key       string          `json:"key" gorm:"primaryKey"`
	Value     json.RawMessage `json:"value" gorm:"type:json"`
	Label     string          `json:"label"`
	IsPublic  bool            `json:"isPublic"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}
