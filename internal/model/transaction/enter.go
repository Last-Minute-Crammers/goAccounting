package transactionModel

import "goAccounting/global/db"

func init() {
	tables := []any{
		Transaction{},
		ExpenseUserStatistic{}, IncomeUserStatistic{},
		ExpenseCategoryStatistic{}, IncomeCategoryStatistic{},
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}
