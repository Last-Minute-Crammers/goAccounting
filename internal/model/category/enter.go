package categoryModel

import (
	"goAccounting/global/constant"
	"goAccounting/global/db"
)

func init() {
	tables := []any{
		Category{},
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}

	initDefaultCategories()
}

// in fact, I don't wonna add Icon, I'm lazy to find icons
func initDefaultCategories() {
	var count int64
	db.InitDb.Model(&Category{}).Count(&count)
	if count > 0 {
		return
	}

	expenses := []Category{
		{Name: "food", Icon: "food", IncomeExpense: constant.Expense},
		{Name: "transportation", Icon: "transportation", IncomeExpense: constant.Expense},
		{Name: "shop", Icon: "shop", IncomeExpense: constant.Expense},
		{Name: "game", Icon: "game", IncomeExpense: constant.Expense},
	}

	incomes := []Category{
		{Name: "salary", Icon: "salary", IncomeExpense: constant.Income},
		{Name: "bonus", Icon: "bonus", IncomeExpense: constant.Income},
		{Name: "invest", Icon: "invest", IncomeExpense: constant.Income},
	}

	db.InitDb.Create(&expenses)
	db.InitDb.Create(&incomes)
}
