package aiModel

import (
	"time"

	"gorm.io/gorm"
)

type reprotType int8

const (
	reportWeekly  reprotType = iota // 周报
	reportMonthly                   // 月报
	reportYearly                    // 年报
)

type reportRecord struct {
	ID           uint `gorm:"primarykey"`
	UserId       uint `gorm:"index"`
	reportRecord reprotType
	CreatedAt    time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
