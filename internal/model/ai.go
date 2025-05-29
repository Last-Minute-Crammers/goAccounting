package model

import (
	"time"

	"gorm.io/gorm"
)

// AIPetStatus AI宠物状态记录
type AIPetStatus struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Mood          string    `gorm:"size:50;not null" json:"mood"`
	Encouragement string    `gorm:"type:text" json:"encouragement"`
	GoalProgress  float64   `gorm:"type:decimal(10,2)" json:"goal_progress"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AIReport AI生成的财务报告
type AIReport struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	ReportType string         `gorm:"size:20;not null" json:"report_type"` // weekly, monthly, yearly
	Content    string         `gorm:"type:text" json:"content"`
	Summary    string         `gorm:"type:text" json:"summary"`
	Advice     string         `gorm:"type:text" json:"advice"`
	Score      int            `gorm:"default:0" json:"score"`
	Period     string         `gorm:"size:50" json:"period"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// AIChatHistory AI对话历史记录
type AIChatHistory struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	UserInput  string    `gorm:"type:text;not null" json:"user_input"`
	AIResponse string    `gorm:"type:text" json:"ai_response"`
	SessionID  string    `gorm:"size:100;index" json:"session_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// AIProcessLog AI处理日志
type AIProcessLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	ServiceType string    `gorm:"size:50;not null" json:"service_type"` // voice, ocr, chat, report, pet
	InputData   string    `gorm:"type:text" json:"input_data"`
	OutputData  string    `gorm:"type:text" json:"output_data"`
	ProcessTime int       `gorm:"default:0" json:"process_time"`           // 处理时间(毫秒)
	Status      string    `gorm:"size:20;default:'success'" json:"status"` // success, failed
	ErrorMsg    string    `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}

func (AIPetStatus) TableName() string {
	return "ai_pet_status"
}

func (AIReport) TableName() string {
	return "ai_reports"
}

func (AIChatHistory) TableName() string {
	return "ai_chat_history"
}

func (AIProcessLog) TableName() string {
	return "ai_process_logs"
}
