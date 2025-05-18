package transactionModel

import (
	commonModel "goAccounting/internal/model/common"
	"time"

	"gorm.io/gorm"
)

type statisticModel interface {
	GetUpdatesValue(amount, count int) map[string]interface{}
	GetDate(tradeTime time.Time) time.Time
	TableName() string
}

type Statistic struct {
	Date   time.Time `gorm:"primaryKey;type:TIMESTAMP"`
	Amount int
	Count  int
	commonModel.BaseModel
}

func (s *Statistic) GetUpdatesValue(amount, count int) map[string]interface{} {
	return map[string]interface{}{
		"amount": gorm.Expr("amount + ?", amount),
		"count":  gorm.Expr("count + ?", count),
	}
}

func (s *Statistic) GetDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
