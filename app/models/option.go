package models

import (
	"gorm.io/datatypes"
	"time"
)

type Option struct {
	Key       string         `json:"key" gorm:"primaryKey"`
	Value     datatypes.JSON `json:"value"`
	Label     string         `json:"label"`
	IsPublic  bool           `json:"isPublic"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
