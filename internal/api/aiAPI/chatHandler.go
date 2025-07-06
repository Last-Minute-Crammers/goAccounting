// 使用 curl 测试 AI 对话接口：
// curl -X POST -H "Content-Type: application/json" -d '{"message":"你好"}' http://localhost:8080/api/public/ai/chat

package aiAPI

import (
	"encoding/json"
	"fmt"
	aiService "goAccounting/internal/service/thirdparty/ai"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	contextFunc "goAccounting/internal/api/util"
	aiModel "goAccounting/internal/model/ai"
	transactionModel "goAccounting/internal/model/transaction"
	transactionService "goAccounting/internal/service/transaction"
	"goAccounting/global/constant"
	userModel "goAccounting/internal/model/user"
	"gorm.io/gorm"
)

type ChatRequest struct {
	Message   string `json:"message" binding:"required"`
	SessionId string `json:"sessionId,omitempty"`
}

type ChatResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Debug   string `json:"debug,omitempty"` // 添加调试信息
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.FormValue("input")
	chatService := &aiService.ChatService{}
	responseText, err := chatService.GetChatResponse(userInput, r.Context())
	if err != nil {
		http.Error(w, "对话失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("智能体回复: %s", responseText)))
}

// Gin 适配器
func GinChatHandler(ctx *gin.Context) {
	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		ctx.JSON(http.StatusBadRequest, ChatResponse{
			Success: false,
			Error:   "输入内容不能为空",
		})
		return
	}

	log.Printf("收到AI对话请求: %s", req.Message)

	// 获取当前用户ID
	userId := contextFunc.ContextFunc.GetUserId(ctx)
	chatService := &aiService.ChatService{}
	var responseText string
	var err error
	if req.SessionId != "" {
		responseText, err = chatService.ContinueSession(req.Message, userId, req.SessionId, ctx.Request.Context())
	} else {
		responseText, err = chatService.GetChatResponseWithUser(req.Message, userId, ctx.Request.Context())
	}
	if err != nil {
		log.Printf("AI对话失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, ChatResponse{
			Success: false,
			Error:   "对话失败: " + err.Error(),
			Debug:   fmt.Sprintf("请求内容: %s", req.Message),
		})
		return
	}

	log.Printf("AI对话成功，响应长度: %d", len(responseText))
	log.Printf("AI响应内容: %s", responseText)
	
	// 检查返回内容是否为空
	if responseText == "" {
		log.Printf("警告: AI返回了空内容")
		ctx.JSON(http.StatusOK, ChatResponse{
			Success: false,
			Error:   "AI返回了空内容",
			Debug:   "API调用成功但返回空字符串",
		})
		return
	}

	ctx.JSON(http.StatusOK, ChatResponse{
		Success: true,
		Data:    responseText,
	})
}

// 聊天历史记录请求
// GET /user/ai/chat/history?offset=0&limit=20
func GinChatHistoryHandler(ctx *gin.Context) {
	userId := contextFunc.ContextFunc.GetUserId(ctx)
	offset := 0
	limit := 20
	if v := ctx.Query("offset"); v != "" {
		offsetInt, err := strconv.Atoi(v)
		if err == nil {
			offset = offsetInt
		}
	}
	if v := ctx.Query("limit"); v != "" {
		limitInt, err := strconv.Atoi(v)
		if err == nil {
			limit = limitInt
		}
	}
	chatService := &aiService.ChatService{}
	history, err := chatService.GetChatHistory(userId, offset, limit, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "获取历史失败: " + err.Error()})
		return
	}
	// 分组：每个sessionId只取最早一条（即每个会话的第一个问题）
	sessionMap := make(map[string]aiModel.ChatRecord)
	for i := len(history) - 1; i >= 0; i-- { // 倒序，保证最早的在后面覆盖
		record := history[i]
		sessionMap[record.SessionId] = record
	}
	var sessions []gin.H
	for _, record := range sessionMap {
		sessions = append(sessions, gin.H{
			"sessionId":     record.SessionId,
			"firstQuestion": record.Input,
			"createdAt":     record.CreatedAt,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": sessions})
}

// 获取指定会话的全部聊天记录
// GET /user/ai/chat/session?sessionId=xxx
func GinChatSessionDetailHandler(ctx *gin.Context) {
	userId := contextFunc.ContextFunc.GetUserId(ctx)
	sessionId := ctx.Query("sessionId")
	if sessionId == "" {
		ctx.JSON(400, gin.H{"success": false, "error": "缺少sessionId参数"})
		return
	}
	chatService := &aiService.ChatService{}
	records, err := chatService.GetSessionHistory(sessionId, 100, ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"success": false, "error": "获取会话详情失败: " + err.Error()})
		return
	}
	// 只返回属于当前用户的消息
	var filtered []aiModel.ChatRecord
	for _, r := range records {
		if r.UserId == userId || r.UserId == 0 {
			filtered = append(filtered, r)
		}
	}
	ctx.JSON(200, gin.H{"success": true, "data": filtered})
}

// GinAIReportHandler 处理AI财务报告生成请求
func GinAIReportHandler(ctx *gin.Context) {
	type Req struct {
		Type  string                 `json:"type"`
		Stats map[string]interface{} `json:"stats"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	userId := contextFunc.ContextFunc.GetUserId(ctx)
	log.Printf("=== 开始AI报告生成请求 ===")
	log.Printf("用户ID: %d", userId)
	log.Printf("报告类型: %s", req.Type)
	log.Printf("前端传递的统计数据: %+v", req.Stats)
	
	// 检查前端是否传递了统计数据
	var stats map[string]interface{}
	var err error
	
	if req.Stats != nil && len(req.Stats) > 0 {
		// 使用前端传递的统计数据
		stats = req.Stats
		log.Printf("使用前端传递的统计数据: %+v", stats)
	} else {
		// 前端没有传递数据，从数据库获取
		stats, err = getAIReportData(userId, req.Type, ctx)
		if err != nil {
			log.Printf("获取AI报告数据失败: %v", err)
			ctx.JSON(500, gin.H{"error": "获取数据失败: " + err.Error()})
			return
		}
	}

	log.Printf("最终传递给AI的统计数据: %+v", stats)

	prompt := buildAIReportPrompt(req.Type, stats)
	log.Printf("生成的AI提示词: %s", prompt)

	// 调用蓝心大模型
	chatService := &aiService.ChatService{}
	resp, err := chatService.GetChatResponse(prompt, ctx)
	if err != nil {
		log.Printf("AI生成失败: %v", err)
		ctx.JSON(500, gin.H{"error": "AI生成失败"})
		return
	}

	log.Printf("AI原始响应: %s", resp)

	// 假设AI返回内容为JSON字符串，解析出summary、suggestion、tags
	var aiResult struct {
		Summary    string   `json:"summary"`
		Suggestion string   `json:"suggestion"`
		Tags       []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(resp), &aiResult); err != nil {
		log.Printf("AI返回解析失败: %v", err)
		ctx.JSON(500, gin.H{"error": "AI返回解析失败", "raw": resp})
		return
	}

	log.Printf("解析后的AI结果: %+v", aiResult)

	// 保存报告到数据库
	db := ctx.MustGet("db").(*gorm.DB)
	tagsJson, _ := json.Marshal(aiResult.Tags)
	now := time.Now()
	var startTime, endTime time.Time
	switch req.Type {
	case "week":
		endTime = now
		startTime = now.AddDate(0, 0, -6)
	case "month":
		endTime = now
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "year":
		endTime = now
		startTime = now.AddDate(0, -11, 1)
	default:
		endTime = now
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	
	// 生成period字符串
	var period string
	switch req.Type {
	case "week":
		year := now.Year()
		week := (now.YearDay() + 6) / 7
		period = fmt.Sprintf("%d-W%02d", year, week)
	case "month":
		period = fmt.Sprintf("%d-%02d", now.Year(), now.Month())
	case "year":
		period = fmt.Sprintf("%d", now.Year())
	default:
		period = fmt.Sprintf("%d-%02d", now.Year(), now.Month())
	}

	err = aiService.GenerateAndSaveReport(db, userId, aiModel.ReportType(req.Type), period, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), aiResult.Summary, aiResult.Suggestion, string(tagsJson))
	if err != nil {
		log.Printf("保存AI报告失败: %v", err)
		// 即使保存失败，也返回生成的报告
	} else {
		log.Printf("AI报告保存成功")
	}

	log.Printf("=== AI报告生成完成 ===")

	ctx.JSON(200, gin.H{
		"summary": aiResult.Summary,
		"suggestion": aiResult.Suggestion,
		"tags": aiResult.Tags,
	})
}

// getAIReportData 根据报告类型获取相应的统计数据
func getAIReportData(userId uint, reportType string, ctx *gin.Context) (map[string]interface{}, error) {
	now := time.Now()
	var startTime, endTime time.Time
	
	// 根据报告类型设置时间范围
	switch reportType {
	case "week":
		endTime = now
		startTime = now.AddDate(0, 0, -6)
	case "month":
		endTime = now
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "year":
		endTime = now
		startTime = now.AddDate(0, -11, 1)
	default:
		endTime = now
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	
	log.Printf("查询时间范围: %s 到 %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	
	db := ctx.MustGet("db").(*gorm.DB)
	dao := transactionModel.NewStatisticDao(db)
	service := transactionService.NewStatisticService()
	
	result := make(map[string]interface{})
	
	// 1. 获取基础统计数据
	condition := transactionModel.NewStatisticConditionBuilder(userId).WithDate(startTime, endTime).Build()
	baseStats, err := dao.GetIeStatisticByCondition(nil, *condition)
	if err != nil {
		log.Printf("获取基础统计数据失败: %v", err)
	} else {
		result["base_statistics"] = baseStats
		log.Printf("基础统计数据: %+v", baseStats)
	}
	
	// 2. 根据报告类型获取不同的统计数据
	switch reportType {
	case "week":
		// 周报：最近一周的统计数据 + 类别统计数据 + 单日统计数据
		err = getWeeklyReportData(userId, startTime, endTime, dao, service, result, ctx)
	case "month":
		// 月报：本月的统计数据 + 类别统计数据 + 单周统计数据
		err = getMonthlyReportData(userId, startTime, endTime, dao, service, result, ctx)
	case "year":
		// 年报：近12个月的统计数据 + 类别统计数据
		err = getYearlyReportData(userId, startTime, endTime, dao, service, result, ctx)
	}
	
	if err != nil {
		log.Printf("获取%s报告数据失败: %v", reportType, err)
		return result, err
	}
	
	return result, nil
}

// getWeeklyReportData 获取周报数据
func getWeeklyReportData(userId uint, startTime, endTime time.Time, dao *transactionModel.StatisticDao, service *transactionService.StatisticService, result map[string]interface{}, ctx *gin.Context) error {
	// 1. 获取最近一周的类别统计数据
	incomeCategoryStats, err := dao.GetCategoryAmountRank(constant.Income, transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: userId},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil)
	if err != nil {
		log.Printf("获取收入类别统计失败: %v", err)
	} else {
		result["income_category_stats"] = incomeCategoryStats
		log.Printf("收入类别统计: %+v", incomeCategoryStats)
	}
	
	expenseCategoryStats, err := dao.GetCategoryAmountRank(constant.Expense, transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: userId},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil)
	if err != nil {
		log.Printf("获取支出类别统计失败: %v", err)
	} else {
		result["expense_category_stats"] = expenseCategoryStats
		log.Printf("支出类别统计: %+v", expenseCategoryStats)
	}
	
	// 2. 获取单日统计数据
	incomeDailyStats, err := dao.GetDayStatisticByCondition(constant.Income, *transactionModel.NewStatisticConditionBuilder(userId).WithDate(startTime, endTime).Build())
	if err != nil {
		log.Printf("获取收入单日统计失败: %v", err)
	} else {
		result["income_daily_stats"] = incomeDailyStats
		log.Printf("收入单日统计: %+v", incomeDailyStats)
	}
	
	expenseDailyStats, err := dao.GetDayStatisticByCondition(constant.Expense, *transactionModel.NewStatisticConditionBuilder(userId).WithDate(startTime, endTime).Build())
	if err != nil {
		log.Printf("获取支出单日统计失败: %v", err)
	} else {
		result["expense_daily_stats"] = expenseDailyStats
		log.Printf("支出单日统计: %+v", expenseDailyStats)
	}
	
	return nil
}

// getMonthlyReportData 获取月报数据
func getMonthlyReportData(userId uint, startTime, endTime time.Time, dao *transactionModel.StatisticDao, service *transactionService.StatisticService, result map[string]interface{}, ctx *gin.Context) error {
	// 1. 获取本月的类别统计数据
	incomeCategoryStats, err := dao.GetCategoryAmountRank(constant.Income, transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: userId},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil)
	if err != nil {
		log.Printf("获取收入类别统计失败: %v", err)
	} else {
		result["income_category_stats"] = incomeCategoryStats
		log.Printf("收入类别统计: %+v", incomeCategoryStats)
	}
	
	expenseCategoryStats, err := dao.GetCategoryAmountRank(constant.Expense, transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: userId},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil)
	if err != nil {
		log.Printf("获取支出类别统计失败: %v", err)
	} else {
		result["expense_category_stats"] = expenseCategoryStats
		log.Printf("支出类别统计: %+v", expenseCategoryStats)
	}
	
	// 2. 获取单周统计数据
	weeklyStats, err := service.GetPeriodStatistics(userId, transactionService.Weekly, startTime, endTime, nil, nil, ctx)
	if err != nil {
		log.Printf("获取单周统计失败: %v", err)
	} else {
		result["weekly_stats"] = weeklyStats
		log.Printf("单周统计: %+v", weeklyStats)
	}
	
	return nil
}

// getYearlyReportData 获取年报数据
func getYearlyReportData(userId uint, startTime, endTime time.Time, dao *transactionModel.StatisticDao, service *transactionService.StatisticService, result map[string]interface{}, ctx *gin.Context) error {
	// 1. 获取近12个月的统计数据
	monthlyStats, err := service.GetPeriodStatistics(userId, transactionService.Monthly, startTime, endTime, nil, nil, ctx)
	if err != nil {
		log.Printf("获取月度统计失败: %v", err)
	} else {
		result["monthly_stats"] = monthlyStats
		log.Printf("月度统计: %+v", monthlyStats)
	}
	
	// 2. 获取12个月的类别统计数据
	incomeCategoryStats, err := dao.GetCategoryAmountRank(constant.Income, transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: userId},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil)
	if err != nil {
		log.Printf("获取收入类别统计失败: %v", err)
	} else {
		result["income_category_stats"] = incomeCategoryStats
		log.Printf("收入类别统计: %+v", incomeCategoryStats)
	}
	
	expenseCategoryStats, err := dao.GetCategoryAmountRank(constant.Expense, transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: userId},
		StartTime: startTime,
		EndTime:   endTime,
	}, nil)
	if err != nil {
		log.Printf("获取支出类别统计失败: %v", err)
	} else {
		result["expense_category_stats"] = expenseCategoryStats
		log.Printf("支出类别统计: %+v", expenseCategoryStats)
	}
	
	return nil
}

// buildAIReportPrompt 组装AI报表prompt
func buildAIReportPrompt(reportType string, stats map[string]interface{}) string {
	var typeText string
	switch reportType {
	case "week":
		typeText = "本周财务数据（金额单位：元）："
	case "month":
		typeText = "本月财务数据（金额单位：元）："
	case "year":
		typeText = "本年财务数据（金额单位：元）："
	default:
		typeText = "财务数据（金额单位：元）："
	}
	
	// 构建更详细的数据描述
	var dataDescription string
	if baseStats, ok := stats["base_statistics"]; ok {
		dataDescription += fmt.Sprintf("基础收支统计（元）：%+v\n", baseStats)
	}
	
	if incomeCategoryStats, ok := stats["income_category_stats"]; ok {
		dataDescription += fmt.Sprintf("收入类别统计（元）：%+v\n", incomeCategoryStats)
	}
	
	if expenseCategoryStats, ok := stats["expense_category_stats"]; ok {
		dataDescription += fmt.Sprintf("支出类别统计（元）：%+v\n", expenseCategoryStats)
	}
	
	if reportType == "week" {
		if incomeDailyStats, ok := stats["income_daily_stats"]; ok {
			dataDescription += fmt.Sprintf("收入单日统计（元）：%+v\n", incomeDailyStats)
		}
		if expenseDailyStats, ok := stats["expense_daily_stats"]; ok {
			dataDescription += fmt.Sprintf("支出单日统计（元）：%+v\n", expenseDailyStats)
		}
	} else if reportType == "month" {
		if weeklyStats, ok := stats["weekly_stats"]; ok {
			dataDescription += fmt.Sprintf("单周统计（元）：%+v\n", weeklyStats)
		}
	} else if reportType == "year" {
		if monthlyStats, ok := stats["monthly_stats"]; ok {
			dataDescription += fmt.Sprintf("月度统计（元）：%+v\n", monthlyStats)
		}
	}
	
	statsJson, _ := json.Marshal(stats)
	
	return typeText + "\n" + dataDescription + "\n完整数据（金额单位：元）：" + string(statsJson) + `\n请你作为理财助手，输出如下JSON格式：{ "summary": "收支总结", "suggestion": "理财建议", "tags": ["标签1", "标签2", "标签3"] }。summary为收支总结，请尽量多结合用户的实际情况，给出有参考意义的总结，注意假如用户长期没有收支记录，可能是用户之前还没有使用此记账软件。suggestion为理财建议，请给出具体的建议，不要过于笼统，保持专业性和实用性。tags为3个简短标签，可适当俏皮有趣一些。所有金额数据都已转换为元为单位。`
}

// 查询历史AI报告
func GetHistoryReportHandler(ctx *gin.Context) {
	log.Printf("GetHistoryReportHandler 被调用")
	userID := contextFunc.ContextFunc.GetUserId(ctx)
	reportType := ctx.Query("type")
	period := ctx.Query("period")

	log.Printf("GetHistoryReportHandler - userID: %d, reportType: %s, period: %s", userID, reportType, period)

	if userID == 0 || reportType == "" || period == "" {
		log.Printf("参数验证失败 - userID: %d, reportType: %s, period: %s", userID, reportType, period)
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数缺失"})
		return
	}

	db := ctx.MustGet("db").(*gorm.DB)
	report, err := aiModel.GetHistoryReport(db, userID, aiModel.ReportType(reportType), period)
	if err != nil {
		log.Printf("查询历史报告失败: %v", err)
		// 返回200状态码，但data为null，表示没有找到历史报告
		ctx.JSON(http.StatusOK, gin.H{"data": nil, "msg": "未找到历史报告"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": report})
}
