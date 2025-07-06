package aiModel

import (
	"time"
)

type ReportType string

const (
	ReportTypeWeek  ReportType = "week"
	ReportTypeMonth ReportType = "month"
	ReportTypeYear  ReportType = "year"
)

type FinancialReport struct {
	ID         uint       `gorm:"primaryKey"`
	UserID     uint       `gorm:"index"`
	Type       ReportType `gorm:"type:varchar(10);index"` // week/month/year
	Period     string     `gorm:"type:varchar(32);index"` // 2024-W27/2024-07/2024
	StartTime  time.Time
	EndTime    time.Time
	Summary    string
	Suggestion string
	Tags       string // JSON数组字符串
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
