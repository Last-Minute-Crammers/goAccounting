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

	// 不在这里创建默认分类，改为在用户注册时创建
}

// CreateDefaultCategoriesForUser 为新用户创建默认分类
func CreateDefaultCategoriesForUser(accountId uint) error {
	// 检查用户是否已有分类
	var count int64
	db.InitDb.Model(&Category{}).Where("account_id = ?", accountId).Count(&count)
	if count > 0 {
		return nil // 已有分类，不重复创建
	}

	// 为新用户创建默认分类
	expenses := []Category{
		{AccountID: accountId, Name: "餐饮", Icon: "food", IncomeExpense: constant.Expense},
		{AccountID: accountId, Name: "购物", Icon: "shop", IncomeExpense: constant.Expense},
		{AccountID: accountId, Name: "交通", Icon: "transportation", IncomeExpense: constant.Expense},
		{AccountID: accountId, Name: "住房", Icon: "house", IncomeExpense: constant.Expense},
		{AccountID: accountId, Name: "娱乐", Icon: "game", IncomeExpense: constant.Expense},
	}

	incomes := []Category{
		{AccountID: accountId, Name: "服游", Icon: "salary", IncomeExpense: constant.Income},
		{AccountID: accountId, Name: "投资", Icon: "invest", IncomeExpense: constant.Income},
	}

	// 创建支出分类
	if err := db.InitDb.Create(&expenses).Error; err != nil {
		return err
	}

	// 创建收入分类
	return db.InitDb.Create(&incomes).Error
}
