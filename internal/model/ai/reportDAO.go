package aiModel

import (
	"gorm.io/gorm"
)

// 创建报告
func CreateReport(db *gorm.DB, report *FinancialReport) error {
	return db.Create(report).Error
}

// 查询报告（按用户、类型、周期）
func GetReport(db *gorm.DB, userID uint, reportType ReportType, period string) (*FinancialReport, error) {
	var report FinancialReport
	err := db.Where("user_id = ? AND type = ? AND period = ?", userID, reportType, period).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// 查询历史AI报告
func GetHistoryReport(db *gorm.DB, userID uint, reportType ReportType, period string) (*FinancialReport, error) {
	return GetReport(db, userID, reportType, period)
}
