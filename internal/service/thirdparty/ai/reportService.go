// report服务为用户提供周、月、年的财务状况和理财建议报告，通过调用云端大模型api实现
package aiService

import (
	"context"
	//"fmt"
	"time"

	"goAccounting/internal/model/ai"
	"gorm.io/gorm"
)

type ReportService struct{}

type FinancialData struct {
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Savings  float64 `json:"savings"`
	Period   string  `json:"period"`
	Category string  `json:"category"`
}

type FinancialReport struct {
	Period      string    `json:"period"`
	Summary     string    `json:"summary"`
	Advice      string    `json:"advice"`
	Score       int       `json:"score"`
	GeneratedAt time.Time `json:"generated_at"`
}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (rs *ReportService) GenerateWeeklyReport(data FinancialData, ctx context.Context) (*FinancialReport, error) {
	// 调用大模型API生成周报
	//prompt := fmt.Sprintf("基于以下财务数据生成周报：收入%.2f，支出%.2f，储蓄%.2f",
	//	data.Income, data.Expense, data.Savings)

	report := &FinancialReport{
		Period:      "weekly",
		Summary:     "本周财务状况良好，支出控制在合理范围内",
		Advice:      "建议继续保持当前的储蓄习惯，可考虑增加投资比例",
		Score:       85,
		GeneratedAt: time.Now(),
	}

	return report, nil
}

func (rs *ReportService) GenerateMonthlyReport(data FinancialData, ctx context.Context) (*FinancialReport, error) {
	// 调用大模型API生成月报
	//prompt := fmt.Sprintf("基于以下财务数据生成月报：收入%.2f，支出%.2f，储蓄%.2f",
	//	data.Income, data.Expense, data.Savings)

	report := &FinancialReport{
		Period:      "monthly",
		Summary:     "本月整体财务表现稳定，达成了预期目标",
		Advice:      "建议下月制定更详细的预算计划，优化支出结构",
		Score:       78,
		GeneratedAt: time.Now(),
	}

	return report, nil
}

func (rs *ReportService) GenerateYearlyReport(data FinancialData, ctx context.Context) (*FinancialReport, error) {
	// 调用大模型API生成年报
	//prompt := fmt.Sprintf("基于以下财务数据生成年报：收入%.2f，支出%.2f，储蓄%.2f",
	//	data.Income, data.Expense, data.Savings)

	report := &FinancialReport{
		Period:      "yearly",
		Summary:     "年度财务目标基本达成，财务健康度良好",
		Advice:      "建议明年设定更高的储蓄目标，考虑多元化投资策略",
		Score:       82,
		GeneratedAt: time.Now(),
	}

	return report, nil
}

func (rs *ReportService) GetUserReports(userId uint, period string, limit int) ([]FinancialReport, error) {
	// TODO: Implement database query
	reports := []FinancialReport{
		{
			Period:      period,
			Summary:     "示例报告",
			Advice:      "示例建议",
			Score:       80,
			GeneratedAt: time.Now(),
		},
	}
	return reports, nil
}

// 生成并保存AI报告
func GenerateAndSaveReport(db *gorm.DB, userID uint, reportType aiModel.ReportType, period string, startTime, endTime string, summary, suggestion, tags string) error {
	// 解析时间字符串
	startTimeParsed, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		// 如果解析失败，使用当前时间
		startTimeParsed = time.Now()
	}
	
	endTimeParsed, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		// 如果解析失败，使用当前时间
		endTimeParsed = time.Now()
	}
	
	report := &aiModel.FinancialReport{
		UserID:     userID,
		Type:       reportType,
		Period:     period,
		StartTime:  startTimeParsed,
		EndTime:    endTimeParsed,
		Summary:    summary,
		Suggestion: suggestion,
		Tags:       tags,
	}
	return aiModel.CreateReport(db, report)
}

// 查询历史AI报告
func GetHistoryReport(db *gorm.DB, userID uint, reportType aiModel.ReportType, period string) (*aiModel.FinancialReport, error) {
	return aiModel.GetReport(db, userID, reportType, period)
}

// 定时批量生成历史AI报告（接口定义，具体实现可后续完善）
func BatchGenerateReports(db *gorm.DB) error {
	// 1. 查询所有活跃用户
	// 2. 按周期（周/月/年）遍历，生成报告内容
	// 3. 调用 GenerateAndSaveReport 保存
	// 4. 可加并发/分片优化
	return nil
}
