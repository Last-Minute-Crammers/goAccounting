package categoryModel

import "goAccounting/global/db"

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
		{Name: "餐饮", Icon: "food", IncomeExpense: "支出"},
		{Name: "交通", Icon: "transportation", IncomeExpense: "支出"},
		{Name: "购物", Icon: "shop", IncomeExpense: "支出"},
		{Name: "娱乐", Icon: "game", IncomeExpense: "支出"},
	}

	incomes := []Category{
		{Name: "工资", Icon: "salary", IncomeExpense: "收入"},
		{Name: "奖金", Icon: "bonus", IncomeExpense: "收入"},
		{Name: "投资", Icon: "invest", IncomeExpense: "收入"},
	}

	db.InitDb.Create(&expenses)
	db.InitDb.Create(&incomes)
}
