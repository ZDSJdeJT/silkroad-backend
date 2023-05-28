package models

import (
	"time"
)

type Setting struct {
	Key         string    `json:"key" gorm:"primaryKey"`
	TextValue   string    `json:"textValue" gorm:"null"`
	NumberValue uint      `json:"numberValue" gorm:"null"`
	IsText      bool      `json:"isText"`
	Min         uint      `json:"min"`
	Max         uint      `json:"max"`
	Label       string    `json:"label"`
	IsPublic    bool      `json:"isPublic"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
