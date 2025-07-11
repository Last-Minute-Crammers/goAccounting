package transactionModel

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/util/timeTool"
	"time"

	"log"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionDao struct {
	db *gorm.DB
}

func NewDao(db ...*gorm.DB) *TransactionDao {
	if len(db) > 0 {
		return &TransactionDao{db: db[0]}
	}
	return &TransactionDao{global.GlobalDb}
}

func (t *TransactionDao) SelectById(id uint, forUpdate bool) (result Transaction, err error) {
	if forUpdate {
		err = t.db.Set("goem:query_option", "FOR UPDATE").First(&result, id).Error

	} else {
		err = t.db.First(&result, id).Error
	}
	return
}

func (t *TransactionDao) Create(info Info, recordType RecordType) (result Transaction, err error) {
	result.Info, result.RecordType = info, recordType
	return result, t.db.Create(&result).Error
}

func (t *TransactionDao) GetListByCondition(condition Condition, offset int, limit int) (
	result []Transaction, err error,
) {
	query := condition.addConditionToQuery(t.db)
	err = query.Limit(limit).Offset(offset).Order("trade_time DESC").Find(&result).Error
	return

}

func (t *TransactionDao) setTimeRangeForQuery(query *gorm.DB, startTime, endTime time.Time) (*gorm.DB, error) {
	switch true {
	case !startTime.IsZero() && !endTime.IsZero():
		query = query.Where("trade_time BETWEEN ? AND ?", startTime, endTime)
	case !startTime.IsZero():
		query = query.Where("trade_time >=", startTime)
	case !endTime.IsZero():
		query = query.Where("trade_time <=", endTime)
	}
	return query, nil
}

func (t *TransactionDao) getAmountCountStatistic(query *gorm.DB, ie constant.IncomeExpense) (
	result global.AmountCount, err error,
) {
	err = query.Where("income_expense = ? ", ie).Select("COUNT(*) as Count,SUM(amount) as Amount").Scan(&result).Error
	return
}

func (t *TransactionDao) getIEStatisticByWhere(ie *constant.IncomeExpense, query *gorm.DB) (
	result global.IEStatistic, err error,
) {
	// 如果ie为nil，查询收入和支出
	if ie == nil {
		result.Income, err = t.getAmountCountStatistic(query, constant.Income)
		if err != nil {
			return
		}
		result.Expense, err = t.getAmountCountStatistic(query, constant.Expense)
		if err != nil {
			return
		}
	} else {
		if ie.QueryIncome() {
			result.Income, err = t.getAmountCountStatistic(query, constant.Income)
			if err != nil {
				return
			}
		}
		if ie.QueryExpense() {
			result.Expense, err = t.getAmountCountStatistic(query, constant.Expense)
			if err != nil {
				return
			}
		}
	}
	return
}

func (t *TransactionDao) GetIeStatisticByCondition(
	ie *constant.IncomeExpense, condition StatisticCondition, extCond *ExtensionCondition,
) (result global.IEStatistic, err error) {
	log.Printf("[txDAO]: GetIeStatisticByCondition - ie: %v, condition: %+v, extCond: %+v", ie, condition, extCond)

	if extCond.IsSet() {
		// transaction table select
		log.Printf("[txDAO]: Using transaction table query due to extension conditions")
		query := t.db.Model(&Transaction{})
		query = condition.ForeignKeyCondition.addConditionToQuery(query)
		query, err = t.setTimeRangeForQuery(
			query, timeTool.ToDay(condition.StartTime), timeTool.ToDay(condition.EndTime),
		)
		if err != nil {
			return
		}
		query = extCond.addConditionToQuery(query)
		result, err = t.getIEStatisticByWhere(ie, query)
	} else {
		// statistic table select
		log.Printf("[txDAO]: Using statistic table query")
		result, err = NewStatisticDao(t.db).GetIeStatisticByCondition(ie, condition)
	}
	if err != nil {
		err = errors.Wrap(err, "transactionDao.GetIeStatisticByCondition")
	}
	log.Printf("[txDAO]: Final result: %+v", result)
	return
}

func (t *TransactionDao) GetAmountRank(
	userId uint, ie constant.IncomeExpense, timeCond TimeCondition,
) (result []Transaction, err error) {
	limit := 10
	query := timeCond.addConditionToQuery(t.db)
	query = query.Where("user_id = ?", userId).Where("income_expense = ?", ie)
	return result, query.Limit(limit).Order("amount DESC").Find(&result).Error
}
