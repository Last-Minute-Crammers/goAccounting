package v1

import (
	aiModel "goAccounting/internal/model/ai"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// 查询历史AI报告
func GetHistoryReportHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	reportType := c.Query("type")
	period := c.Query("period")

	if userID == 0 || reportType == "" || period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数缺失"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	report, err := aiModel.GetHistoryReport(db, userID, aiModel.ReportType(reportType), period)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "未找到历史报告"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
} 