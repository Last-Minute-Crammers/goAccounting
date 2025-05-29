package aiModel

import (
	"time"
	"gorm.io/gorm"
)

type FinancialReport struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserId      uint           `gorm:"index" json:"user_id"`
	Period      string         `json:"period"`      // weekly, monthly, yearly
	Summary     string         `gorm:"type:text" json:"summary"`
	Advice      string         `gorm:"type:text" json:"advice"`
	Score       int            `json:"score"`
	Income      float64        `json:"income"`
	Expense     float64        `json:"expense"`
	Savings     float64        `json:"savings"`
	GeneratedAt time.Time      `json:"generated_at"`
	CreatedAt   time.Time      `gorm:"type:TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ReportDAO struct {
	db *gorm.DB
}

func NewReportDAO() *ReportDAO {
	return &ReportDAO{}
}

func (dao *ReportDAO) Create(report *FinancialReport) error {
	return dao.db.Create(report).Error
}

func (dao *ReportDAO) GetByUserAndPeriod(userId uint, period string, limit int) ([]FinancialReport, error) {
	var reports []FinancialReport
	err := dao.db.Where("user_id = ? AND period = ?", userId, period).
		Order("generated_at DESC").
		Limit(limit).
		Find(&reports).Error
	return reports, err
}

func (dao *ReportDAO) GetLatestByUser(userId uint) (*FinancialReport, error) {
	var report FinancialReport
	err := dao.db.Where("user_id = ?", userId).
		Order("generated_at DESC").
		First(&report).Error
	return &report, err
}
