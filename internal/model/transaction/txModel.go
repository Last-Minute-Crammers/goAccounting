package transactionModel

import (
	"errors"
	"goAccounting/global/constant"
	categoryModel "goAccounting/internal/model/category"
	commonModel "goAccounting/internal/model/common"
	"log"
	"time"

	"gorm.io/gorm"
)

type RecordType int8

const (
	// we'll only use the 1st one
	RecordTypeOfManual RecordType = iota
	RecordTypeOfTiming
	RecordTypeOfSync
	RecordTypeOfImport
)

type Info struct {
	UserId, AccountId, CategoryId uint
	IncomeExpense                 constant.IncomeExpense
	Amount                        int
	Remark                        string
	TradeTime                     time.Time `gorm:"type:TIMESTAMP"`
}

type Transaction struct {
	Info
	commonModel.BaseModel
	ID         uint `gorm:"primarykey"`
	RecordType RecordType
	CreatedAt  time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

func (i *Info) CheckValid(db *gorm.DB) error {
	log.Println("[model]: get into txModelCheckvalid")
	category, err := categoryModel.NewDao(db).SelectById(i.CategoryId)
	if err != nil {
		return err
	}
	if i.Amount <= 0 {
		return errors.New("transaction CheckValid: amount must be positive")
	}
	if i.IncomeExpense != category.IncomeExpense {
		return errors.New("transaction CheckValid: IncomeExpense mismatch with category")
	}
	return nil
}
