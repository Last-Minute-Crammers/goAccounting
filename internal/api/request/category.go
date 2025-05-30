package request

import "goAccounting/global/constant"

type CategoryCreateOne struct {
	Name          string                     `json:"name" binding:"required"`
	IncomeExpense constant.IncomeExpense     `json:"income_expense" binding:"required,oneof=income expense"`
	Icon          string                     `json:"icon"`
}

type CategoryUpdateOne struct {
	Name          string                     `json:"name"`
	IncomeExpense constant.IncomeExpense     `json:"income_expense" binding:"omitempty,oneof=income expense"`
	Icon          string                     `json:"icon"`
}
