package request

import (
	"goAccounting/global/constant"
	transactionModel "goAccounting/internal/model/transaction"
	"time"
)

type TransactionCreateOne struct {
	Amount        int
	CategoryId    uint
	IncomeExpense constant.IncomeExpense
	Remark        string
	TradeTime     time.Time
}

type TransactionUpdateOne struct {
	UserId        uint
	Amount        int
	CategoryId    uint
	IncomeExpense constant.IncomeExpense
	Remark        string
	TradeTime     time.Time
}

type TransactionQueryCondition struct {
	CategoryIds   *[]uint
	IncomeExpense *constant.IncomeExpense `binding:"omitempty,oneof=income expense"`
	MinimumAmount *int                    `binding:"omitempty,min=0"`
	MaximumAmount *int                    `binding:"omitempty,min=0"`
	TimeFrame
}

func (t *TransactionQueryCondition) GetCondition() transactionModel.Condition {
	startTime := t.TimeFrame.StartTime
	endTime := t.TimeFrame.EndTime
	return transactionModel.Condition{
		IncomeExpense:       t.IncomeExpense,
		TimeCondition:       transactionModel.TimeCondition{TradeStartTime: &startTime, TradeEndTime: &endTime},
		ForeignKeyCondition: t.GetForeignKeyCondition(),
		ExtensionCondition:  t.GetExtensionCondition(),
	}
}

func (t *TransactionQueryCondition) GetForeignKeyCondition() transactionModel.ForeignKeyCondition {
	return transactionModel.ForeignKeyCondition{
		CategoryIds: t.CategoryIds,
	}
}

func (t *TransactionQueryCondition) GetStatisticCondition() transactionModel.StatisticCondition {
	return transactionModel.StatisticCondition{
		ForeignKeyCondition: t.GetForeignKeyCondition(),
		StartTime:           t.StartTime,
		EndTime:             t.EndTime,
	}
}

func (t *TransactionQueryCondition) GetExtensionCondition() transactionModel.ExtensionCondition {
	return transactionModel.ExtensionCondition{
		MinAmount: t.MinimumAmount,
		MaxAmount: t.MaximumAmount,
	}
}

type TransactionGetList struct {
	TransactionQueryCondition
	PageData
}

type TransactionTotal struct {
	TransactionQueryCondition
}

type TransactionMonthStatistic struct {
	TransactionQueryCondition
}

type TransactionDayStatistic struct {
	CategoryIds   *[]uint
	IncomeExpense *constant.IncomeExpense `binding:"omitempty,oneof=income expense"`
	TimeFrame
}

type TransactionCategoryAmountRank struct {
	IncomeExpense constant.IncomeExpense `binding:"required,oneof=income expense"`
	Limit         *int                   `binding:"omitempty"`
	TimeFrame
}

type TransactionAmountRank struct {
	IncomeExpense constant.IncomeExpense `binding:"required,oneof=income expense"`
	TimeFrame
}
