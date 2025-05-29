package aiAPI

import (
	"encoding/json"
	"net/http"
	"strconv"
	aiService "goAccounting/internal/service/thirdparty/ai"
)

type ReportHandler struct {
	reportService *aiService.ReportService
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{
		reportService: aiService.NewReportService(),
	}
}

func (h *ReportHandler) GenerateWeeklyReport(w http.ResponseWriter, r *http.Request) {
	var data aiService.FinancialData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GenerateWeeklyReport(data, r.Context())
	if err != nil {
		http.Error(w, "生成周报失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GenerateMonthlyReport(w http.ResponseWriter, r *http.Request) {
	var data aiService.FinancialData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GenerateMonthlyReport(data, r.Context())
	if err != nil {
		http.Error(w, "生成月报失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetUserReports(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	period := r.URL.Query().Get("period")
	limitStr := r.URL.Query().Get("limit")

	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	reports, err := h.reportService.GetUserReports(uint(userId), period, limit)
	if err != nil {
		http.Error(w, "获取报告失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
