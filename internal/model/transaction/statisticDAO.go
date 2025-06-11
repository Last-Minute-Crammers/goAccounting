package transactionModel

import (
	"errors"
	"goAccounting/global"
	"goAccounting/global/constant"
	userModel "goAccounting/internal/model/user"
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

// GetTotalStatistics gets total income and expense statistics from user registration
func (s *StatisticDao) GetTotalStatistics(userId uint) (result global.IEStatistic, err error) {
	// Get total income
	err = s.db.Table("transaction_income_account_statistic").
		Where("user_id = ?", userId).
		Select("COALESCE(SUM(amount), 0) as amount, COALESCE(SUM(count), 0) as count").
		Scan(&result.Income).Error
	if err != nil {
		return
	}
	
	// Get total expense
	err = s.db.Table("transaction_expense_account_statistic").
		Where("user_id = ?", userId).
		Select("COALESCE(SUM(amount), 0) as amount, COALESCE(SUM(count), 0) as count").
		Scan(&result.Expense).Error
	
	return result, err
}

// set time condition
type CategoryAmountRankCondition struct {
	User      userModel.User // in fact, we don't need it
	StartTime time.Time
	EndTime   time.Time
}

func (c *CategoryAmountRankCondition) Local() {
	location := time.Local
	c.StartTime = c.StartTime.In(location)
	c.EndTime = c.EndTime.In(location)
}

// for CategoryStatistic result
type CategoryAmountRank struct {
	CategoryId uint
	global.AmountCount
}

func (s *StatisticDao) GetCategoryAmountRank(
	ie constant.IncomeExpense, condition CategoryAmountRankCondition, limit *int,
) (result []CategoryAmountRank, err error) {
	condition.Local()
	query := s.db.Where("user_id = ?", condition.User.ID)
	query = query.Where("data BETWEEN ? AND ?", condition.StartTime, condition.EndTime)
	// group means group :)
	query = query.Select("SUM(amount) as Amount,SUM(count) as Count,category_id").Group("category_id")
	if limit != nil {
		query = query.Limit(*limit)
	}
	err = s.queryIE(ie, query).Order("Amount desc").Find(&result).Error
	return result, err
}
