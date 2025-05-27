package transactionModel

import (
	"goAccounting/global/constant"
	"goAccounting/util/timeTool"

	"time"

	"gorm.io/gorm"
)

// ForeignKeyCondition 交易外键查询条件 用于交易记录和统计的查询
type ForeignKeyCondition struct {
	UserId      uint
	CategoryIds *[]uint
}

type TimeCondition struct {
	TradeStartTime *time.Time
	TradeEndTime   *time.Time
}

type ExtensionCondition struct {
	MinAmount, MaxAmount *int
}

type Condition struct {
	ForeignKeyCondition
	TimeCondition
	ExtensionCondition
	IncomeExpense *constant.IncomeExpense
}

func (f *ForeignKeyCondition) addConditionToQuery(db *gorm.DB) *gorm.DB {
	query := db.Where("user_id = ?", f.UserId)
	if f.CategoryIds != nil {
		query = query.Where("category_id IN (?)", *f.CategoryIds)
	}
	return query
}

func (f *ForeignKeyCondition) getInComeStatisticModel() statisticModel {
	if f.CategoryIds == nil {
		return &IncomeUserStatistic{}

	} else {
		return &IncomeCategoryStatistic{}
	}
}

func (f *ForeignKeyCondition) getExpendStatisticModel() statisticModel {
	if f.CategoryIds == nil {
		return &ExpenseUserStatistic{}

	} else {
		return &ExpenseCategoryStatistic{}
	}
}

func (f *ForeignKeyCondition) GetStatisticTableName(ie constant.IncomeExpense) string {
	var model statisticModel
	if ie == constant.Income {
		model = f.getInComeStatisticModel()

	} else {
		model = f.getExpendStatisticModel()

	}
	if model == nil {
		return "transaction"
	}

	return model.TableName()
}

func NewTimeCondition() *TimeCondition {
	return &TimeCondition{}
}

func (tc *TimeCondition) SetTradeTimes(startTime, endTime time.Time) {
	tc.TradeStartTime = &startTime
	tc.TradeEndTime = &endTime
}

func (tc *TimeCondition) addConditionToQuery(query *gorm.DB) *gorm.DB {
	if tc.TradeStartTime != nil {
		query = query.Where("trade_time >= ?", *tc.TradeStartTime)

	}
	if tc.TradeEndTime != nil {
		query = query.Where("trade_time <= ?", *tc.TradeEndTime)

	}
	return query
}

func (ec *ExtensionCondition) IsSet() bool {
	return ec != nil && (ec.MinAmount != nil || ec.MaxAmount != nil)
}

func (ec *ExtensionCondition) addConditionToQuery(query *gorm.DB) *gorm.DB {
	if ec.MinAmount != nil {
		query.Where("amount >= ?", *ec.MinAmount)

	}
	if ec.MaxAmount != nil {
		query.Where("amount <= ?", *ec.MaxAmount)

	}
	return query
}

func (c *Condition) addConditionToQuery(db *gorm.DB) *gorm.DB {
	query := c.ForeignKeyCondition.addConditionToQuery(db)
	query = c.TimeCondition.addConditionToQuery(query)
	query = c.ExtensionCondition.addConditionToQuery(query)
	if c.IncomeExpense != nil {
		query = query.Where("income_expense = ?", *c.IncomeExpense)

	}
	return query
}

// time location搞默认
type StatisticCondition struct {
	ForeignKeyCondition
	StartTime time.Time
	EndTime   time.Time
}

func (s *StatisticCondition) getLocation() *time.Location {
	return time.Local
}

func (s *StatisticCondition) addConditionToQuery(db *gorm.DB) *gorm.DB {
	query := s.ForeignKeyCondition.addConditionToQuery(db)
	switch true {
	case !s.StartTime.IsZero() && !s.EndTime.IsZero():
		query = query.Where(
			"date BETWEEN ? AND ?", timeTool.ToDay(s.StartTime.In(s.getLocation())),
			timeTool.ToDay(s.EndTime.In(s.getLocation())),
		)
	case !s.StartTime.IsZero():
		query = query.Where("date >=", timeTool.ToDay(s.StartTime.In(s.getLocation())))
	case !s.EndTime.IsZero():
		query = query.Where("date <=", timeTool.ToDay(s.EndTime.In(s.getLocation())))
	}
	return query
}

func (s *StatisticCondition) CheckAvailability() bool {
	if s.CategoryIds != nil && len(*s.CategoryIds) == 0 {
		return false
	}
	return true
}

type StatisticConditionBuilder struct {
	condition StatisticCondition
}

func NewStatisticConditionBuilder(userId uint) *StatisticConditionBuilder {
	return &StatisticConditionBuilder{
		condition: StatisticCondition{
			ForeignKeyCondition: ForeignKeyCondition{UserId: userId},
		},
	}
}

// WithCategoryIds 设置交易类型ids
func (b *StatisticConditionBuilder) WithCategoryIds(Ids []uint) *StatisticConditionBuilder {
	b.condition.CategoryIds = &Ids
	return b
}

// WithDate 设置时间范围
func (b *StatisticConditionBuilder) WithDate(startTime, endTime time.Time) *StatisticConditionBuilder {
	b.condition.StartTime = startTime
	b.condition.EndTime = endTime
	return b
}

// Build 构建 StatisticCondition 实例
func (b *StatisticConditionBuilder) Build() *StatisticCondition {
	return &b.condition
}
