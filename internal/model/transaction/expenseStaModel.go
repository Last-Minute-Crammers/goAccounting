package transactionModel

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type ExpenseUserStatistic struct {
	Statistic
	UserId uint `gorm:"primaryKey"`
}

func (e *ExpenseUserStatistic) TableName() string {
	return "transaction_expense_account_statistic"
}

// accumulate the new statistical data into existing statistical records
// we use int as amount becase its' unit is cent to avoid accuracy problem
func (e *ExpenseUserStatistic) Accumulate(
	tradeTime time.Time, userId uint, amount int, count int, tx *gorm.DB,
) error {
	tradeTime = e.GetDate(tradeTime)
	where := tx.Model(e).Where("date = ? AND user_id = ?", tradeTime, userId)
	updatesValue := e.GetUpdatesValue(amount, count)
	update := where.Updates(updatesValue)
	err := update.Error
	if err != nil {
		return err
	}
	if update.RowsAffected == 0 {
		e.Date = tradeTime
		e.UserId = userId
		e.Amount = amount
		e.Count = count
		err = tx.Create(e).Error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = where.Updates(updatesValue).Error
		}
	}
	return err
}

type ExpenseCategoryStatistic struct {
	Statistic
	CategoryId uint `gorm:"primaryKey"`
	UserId     uint
}

func (e *ExpenseCategoryStatistic) TableName() string {
	return "transaction_expense_category_statistic"
}

func (e *ExpenseCategoryStatistic) Accumulate(
	tradeTime time.Time, userId uint, categoryId uint, amount int, count int, tx *gorm.DB,
) error {
	tradeTime = e.GetDate(tradeTime)
	where := tx.Model(e).Where(
		"date = ? AND category_id = ? AND user_id = ?", tradeTime, categoryId, userId,
	)
	updatesValue := e.GetUpdatesValue(amount, count)
	update := where.Updates(updatesValue)
	err := update.Error
	if err != nil {
		return err
	}
	if update.RowsAffected == 0 {
		e.Date = tradeTime
		e.CategoryId = categoryId
		e.UserId = userId
		e.Amount = amount
		e.Count = count
		err = tx.Create(e).Error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = where.Updates(updatesValue).Error
		}
	}
	return err
}
