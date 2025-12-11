package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"goAccounting/initialize"
	aiModel "goAccounting/internal/model/ai"
)

func main() {
	// 直接使用初始化后的数据库连接
	// initialize 包的 init() 函数会自动执行

	// 等待初始化完成
	time.Sleep(time.Second * 2)

	// 生成测试报告数据
	generateTestReports()

	fmt.Println("测试报告数据生成完成")
}

func generateTestReports() {
	// 使用 initialize.Database 而不是 db.Db
	if initialize.Database == nil {
		log.Fatal("数据库未初始化")
	}

	// 为用户ID 1生成一些测试报告
	userID := uint(1)

	// 生成月报
	monthlyReports := []struct {
		period     string
		summary    string
		suggestion string
		tags       []string
	}{
		{
			period:     "2024-01",
			summary:    "1月份财务状况良好，收入稳定，支出控制在合理范围内。",
			suggestion: "建议继续保持当前的储蓄习惯，可考虑增加投资比例。",
			tags:       []string{"收入稳定", "支出合理", "储蓄良好"},
		},
		{
			period:     "2024-02",
			summary:    "2月份支出略有增加，主要是春节期间的消费。",
			suggestion: "建议制定更详细的预算计划，控制非必要支出。",
			tags:       []string{"节日消费", "预算控制", "理性消费"},
		},
		{
			period:     "2024-03",
			summary:    "3月份财务状况回升，储蓄率有所提高。",
			suggestion: "继续保持良好的理财习惯，可以考虑多元化投资。",
			tags:       []string{"储蓄提升", "理财习惯", "投资建议"},
		},
	}

	for _, report := range monthlyReports {
		tagsJson, _ := json.Marshal(report.tags)
		financialReport := &aiModel.FinancialReport{
			UserID:     userID,
			Type:       aiModel.ReportTypeMonth,
			Period:     report.period,
			Summary:    report.summary,
			Suggestion: report.suggestion,
			Tags:       string(tagsJson),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := aiModel.CreateReport(initialize.Database, financialReport); err != nil {
			log.Printf("创建月报失败 %s: %v", report.period, err)
		} else {
			fmt.Printf("成功创建月报: %s\n", report.period)
		}
	}

	// ... 其他报告生成代码类似修改
}
