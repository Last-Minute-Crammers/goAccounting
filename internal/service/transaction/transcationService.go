package transactionService

import (
	"goAccounting/global/constant"
	"goAccounting/global/db"
	transactionModel "goAccounting/internal/model/transaction"
	"log"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Transaction struct{}

func (t *Transaction) Create(
	transInfo transactionModel.Info,
	recordType transactionModel.RecordType,
	ctx context.Context,
) (transactionModel.Transaction, error) {
	log.Println("[service]: get into create")
	err := transInfo.CheckValid(db.GetDb(ctx))
	if err != nil {
		return transactionModel.Transaction{}, err
	}
	err = t.UpdateStatistic(transInfo, ctx)
	if err != nil {
		return transactionModel.Transaction{}, err
	}
	return transactionModel.NewDao(db.GetDb(ctx)).Create(transInfo, recordType)
}

func (t *Transaction) Update(
	id uint,
	updatedInfo transactionModel.Info,
	ctx context.Context,
) error {
	txDao := transactionModel.NewDao(db.GetDb(ctx))
	transaction, err := txDao.SelectById(id, true)
	if (err != nil) {
		return err
	}
	transaction.Info = updatedInfo
	err = t.UpdateStatistic(transaction.Info, ctx)
	if (err != nil) {
		return err
	}
	return db.GetDb(ctx).Save(&transaction).Error
}

func (t *Transaction) Delete(id uint, ctx context.Context) error {
	txDao := transactionModel.NewDao(db.GetDb(ctx))
	transaction, err := txDao.SelectById(id, true)
	if (err != nil) {
		return err
	}
	transInfo := transaction.Info
	transInfo.Amount = -transInfo.Amount
	err = t.UpdateStatistic(transInfo, ctx)
	if (err != nil) {
		return err
	}
	return db.GetDb(ctx).Delete(&transaction).Error
}

func (t *Transaction) UpdateStatistic(
	transInfo transactionModel.Info,
	ctx context.Context,
) error {
	dbInstance := db.GetDb(ctx)

	// 使用事务确保数据一致性
	return dbInstance.Transaction(func(tx *gorm.DB) error {
		// 更新用户级别统计
		if err := t.updateUserStatistic(transInfo, tx); err != nil {
			return errors.Wrap(err, "failed to update user statistics")
		}

		// 更新类别级别统计
		if err := t.updateCategoryStatistic(transInfo, tx); err != nil {
			return errors.Wrap(err, "failed to update category statistics")
		}

		return nil
	})
}

// updateUserStatistic 每次记账更新用户单日统计
func (t *Transaction) updateUserStatistic(transInfo transactionModel.Info, tx *gorm.DB) error {
	if transInfo.IncomeExpense == constant.Income {
		statistic := transactionModel.IncomeUserStatistic{}
		return statistic.Accumulate(
			transInfo.TradeTime, transInfo.UserId, transInfo.Amount, 1, tx,
		)
	} else if transInfo.IncomeExpense == constant.Expense {
		statistic := transactionModel.ExpenseUserStatistic{}
		return statistic.Accumulate(
			transInfo.TradeTime, transInfo.UserId, transInfo.Amount, 1, tx,
		)
	}

	return errors.New("invalid income/expense type")
}

// updateCategoryStatistic 每次记账更新用户单日单类统计
func (t *Transaction) updateCategoryStatistic(transInfo transactionModel.Info, tx *gorm.DB) error {
	if transInfo.IncomeExpense == constant.Income {
		statistic := transactionModel.IncomeCategoryStatistic{}
		return statistic.Accumulate(
			transInfo.TradeTime, transInfo.UserId, transInfo.CategoryId, transInfo.Amount, 1, tx,
		)
	} else if transInfo.IncomeExpense == constant.Expense {
		statistic := transactionModel.ExpenseCategoryStatistic{}
		return statistic.Accumulate(
			transInfo.TradeTime, transInfo.UserId, transInfo.CategoryId, transInfo.Amount, 1, tx,
		)
	}

	return errors.New("invalid income/expense type")
}
