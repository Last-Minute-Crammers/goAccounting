package aiModel

import (
	"time"
	"gorm.io/gorm"
)

type ChatRecord struct {
	ID        uint           `gorm:"primarykey"`
	SessionId string         `gorm:"type:varchar(36);index;not null"` // 会话ID
	RequestId string         `gorm:"type:varchar(36);index;not null"` // 问答ID
	UserId    uint           `gorm:"index"`
	Input     string         `gorm:"type:text"`
	Response  string         `gorm:"type:text"`
	HistorySummary string    `gorm:"type:text"` // 整理后的对话历史
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
