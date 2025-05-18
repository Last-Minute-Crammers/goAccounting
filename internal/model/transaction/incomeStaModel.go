package transactionModel

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type IncomeUserStatistic struct {
	Statistic
	UserId uint `gorm:"primaryKey"`
}

func (i *IncomeUserStatistic) TableName() string {
	return "transaction_income_account_statistic"
}

// accumulate the new statistical data into existing statistical records
// we use int as amount becase its' unit is cent to avoid accuracy problem
func (i *IncomeUserStatistic) Accumulate(
	tradeTime time.Time, userId uint, amount int, count int, tx *gorm.DB,
) error {
	tradeTime = i.GetDate(tradeTime)
	where := tx.Model(i).Where("date = ? AND user_id = ?", tradeTime, userId)
	updatesValue := i.GetUpdatesValue(amount, count)
	update := where.Updates(updatesValue)
	err := update.Error
	if err != nil {
		return err
	}
	if update.RowsAffected == 0 {
		i.Date = tradeTime
		i.UserId = userId
		i.Amount = amount
		i.Count = count
		err = tx.Create(i).Error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = where.Updates(updatesValue).Error
		}
	}
	return err
}

type IncomeCategoryStatistic struct {
	Statistic
	CategoryId uint `gorm:"primaryKey"`
	UserId     uint
}

func (i *IncomeCategoryStatistic) TableName() string {
	return "transaction_income_category_statistic"
}

func (i *IncomeCategoryStatistic) Accumulate(
	tradeTime time.Time, userId uint, categoryId uint, amount int, count int, tx *gorm.DB,
) error {
	tradeTime = i.GetDate(tradeTime)
	where := tx.Model(i).Where(
		"date = ? AND category_id = ? AND user_id = ?", tradeTime, categoryId, userId,
	)
	updatesValue := i.GetUpdatesValue(amount, count)
	update := where.Updates(updatesValue)
	err := update.Error
	if err != nil {
		return err
	}
	if update.RowsAffected == 0 {
		i.Date = tradeTime
		i.CategoryId = categoryId
		i.UserId = userId
		i.Amount = amount
		i.Count = count
		err = tx.Create(i).Error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = where.Updates(updatesValue).Error
		}
	}
	return err
}
