package models

import (
	"gorm.io/gorm"

	"github.com/fish/ai-tools/service"
)

// LabPrompt model
type LabPrompt struct {
	ID   int64  `gorm:"column:id" json:"id,omitempty"`
	Text string `gorm:"column:text" json:"text"`
	Type int    `gorm:"column:type" json:"type"`
}

// DB the db instance of MSoundComment model
func (m LabPrompt) DB() *gorm.DB {
	return service.DB.Table(m.TableName())
}

// TableName for current model
func (LabPrompt) TableName() string {
	return "lab_prompt"
}
