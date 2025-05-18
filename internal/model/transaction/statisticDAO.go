package transactionModel

import (
	"errors"
	"goAccounting/global"
	"goAccounting/global/constant"
	"log"
	"time"

	"gorm.io/gorm"
)

type StatisticDao struct {
	db *gorm.DB
}

func NewStatisticDao(db ...*gorm.DB) *StatisticDao {
	if len(db) > 0 && db[0] != nil {
		return &StatisticDao{db: db[0]}
	}
	return &StatisticDao{db: global.GlobalDb}
}

func (s *StatisticDao) queryIE(ie constant.IncomeExpense, db *gorm.DB) *gorm.DB {
	if ie == constant.Expense {
		return db.Model(&ExpenseCategoryStatistic{})
	} else {

		return db.Model(&IncomeCategoryStatistic{})
	}

}

type DayStatistic struct {
	global.AmountCount
	Date time.Time
}

func (s *StatisticDao) GetDayStatisticByCondition(
	ie constant.IncomeExpense, condition StatisticCondition,
) (result []DayStatistic, err error) {
	if !condition.CheckAvailability() {
		log.Println("invalid categoriesId")
		return
	}
	query := condition.addConditionToQuery(s.db)
	// SELECT SUM(amount) as Amount, SUM(count) as Count, date
	// FROM your_table_name
	// GROUP BY date;
	query.Select("SUM(amount) as Amount,SUM(count) as Count,date").Group("date")
	err = query.Table(condition.GetStatisticTableName(ie)).Find(&result).Error
	return result, err

}

func (s *StatisticDao) GetAmountCountByCondition(condition StatisticCondition, ie constant.IncomeExpense) (
	result global.AmountCount, err error,
) {
	if !condition.CheckAvailability() {
		return
	}
	query := condition.addConditionToQuery(s.db).Table(condition.GetStatisticTableName(ie))
	err = query.Select("SUM(amount) as amount,SUM(count) as count").Scan(&result).Error
	return result, err
}

// GetIeStatisticByCondition 查询收支统计 返回 global.IEStatistic
func (s *StatisticDao) GetIeStatisticByCondition(ie *constant.IncomeExpense, condition StatisticCondition) (
	result global.IEStatistic, err error,
) {
	if !condition.CheckAvailability() {
		return result, errors.New("wrong check categories")
	}
	query := condition.addConditionToQuery(s.db)
	if ie.QueryIncome() {
		err = query.Table(condition.GetStatisticTableName(constant.Income)).Select("SUM(amount) as amount,SUM(count) as count").Scan(&result.Income).Error
		if err != nil {
			return
		}
	}
	if ie.QueryExpense() {
		err = query.Table(condition.GetStatisticTableName(constant.Expense)).Select("SUM(amount) as amount,SUM(count) as count").Scan(&result.Expense).Error
		if err != nil {
			return
		}
	}
	return result, err
}
