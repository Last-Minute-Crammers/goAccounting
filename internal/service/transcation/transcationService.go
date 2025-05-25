package transcationService

import (
	"goAccounting/global/constant"
	"goAccounting/global/db"
	transactionModel "goAccounting/internal/model/transaction"

	"golang.org/x/net/context"
)

type Transaction struct{}

func (t *Transaction) Create(
	transInfo transactionModel.Info,
	recordType transactionModel.RecordType,
	ctx context.Context,
) (transactionModel.Transaction, error) {
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
	if err != nil {
		return err
	}
	transaction.Info = updatedInfo
	err = t.UpdateStatistic(transaction.Info, ctx)
	if err != nil {
		return err
	}
	return db.GetDb(ctx).Save(&transaction).Error
}

func (t *Transaction) Delete(id uint, ctx context.Context) error {
	txDao := transactionModel.NewDao(db.GetDb(ctx))
	transaction, err := txDao.SelectById(id, true)
	if err != nil {
		return err
	}
	transInfo := transaction.Info
	transInfo.Amount = -transInfo.Amount
	err = t.UpdateStatistic(transInfo, ctx)
	if err != nil {
		return err
	}
	return db.GetDb(ctx).Delete(&transaction).Error
}

func (t *Transaction) UpdateStatistic(
	transInfo transactionModel.Info,
	ctx context.Context,
) error {
	dbInstance := db.GetDb(ctx)
	if transInfo.IncomeExpense == constant.Income {
		statistic := transactionModel.IncomeUserStatistic{}
		err := statistic.Accumulate(
			transInfo.TradeTime, transInfo.UserId, transInfo.Amount, 1, dbInstance,
		)
		if err != nil {
			return err
		}
	} else if transInfo.IncomeExpense == constant.Expense {
		statistic := transactionModel.ExpenseUserStatistic{}
		err := statistic.Accumulate(
			transInfo.TradeTime, transInfo.UserId, transInfo.Amount, 1, dbInstance,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
