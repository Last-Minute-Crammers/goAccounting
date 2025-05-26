package aiModel

import (
	"time"

	"gorm.io/gorm"
)

type ChatRecord struct {
	ID        uint           `gorm:"primarykey"`
	UserId    uint           `gorm:"index"`
	Input     string         `gorm:"type:text"`
	Response  string         `gorm:"type:text"`
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
