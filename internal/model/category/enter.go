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

	// 应用启动时为用户ID=1创建默认分类（如果不存在）
	initDefaultCategoriesForUser1()
}

// initDefaultCategoriesForUser1 为用户ID=1创建默认分类（仅用于开发测试）
func initDefaultCategoriesForUser1() {
	// 先清理无效的旧数据（account_id 为 NULL 的记录）
	db.InitDb.Where("account_id IS NULL").Delete(&Category{})

	// 检查是否已有数据
	var count int64
	db.InitDb.Model(&Category{}).Where("account_id = ?", 1).Count(&count)
	if count > 0 {
		return // 已有数据，跳过
	}

	// 创建默认分类
	CreateDefaultCategoriesForUser(1)
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
		{UserID: accountId, Name: "餐饮", Icon: "food", IncomeExpense: constant.Expense},
		{UserID: accountId, Name: "购物", Icon: "shop", IncomeExpense: constant.Expense},
		{UserID: accountId, Name: "交通", Icon: "transportation", IncomeExpense: constant.Expense},
		{UserID: accountId, Name: "住房", Icon: "house", IncomeExpense: constant.Expense},
		{UserID: accountId, Name: "娱乐", Icon: "game", IncomeExpense: constant.Expense},
	}

	incomes := []Category{
		{UserID: accountId, Name: "服游", Icon: "salary", IncomeExpense: constant.Income},
		{UserID: accountId, Name: "投资", Icon: "invest", IncomeExpense: constant.Income},
	}

	// 创建支出分类
	if err := db.InitDb.Create(&expenses).Error; err != nil {
		return err
	}

	// 创建收入分类
	return db.InitDb.Create(&incomes).Error
}
