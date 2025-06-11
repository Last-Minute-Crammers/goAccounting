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
	log.Printf("[model]: 验证交易信息 - CategoryId: %d, IncomeExpense: %s", i.CategoryId, i.IncomeExpense)
	
	category, err := categoryModel.NewDao(db).SelectById(i.CategoryId)
	if err != nil {
		log.Printf("[model]: 查询分类失败 - CategoryId: %d, Error: %v", i.CategoryId, err)
		return errors.New("找不到指定的分类")
	}
	
	log.Printf("[model]: 分类信息 - ID: %d, Name: %s, IncomeExpense: %s", category.ID, category.Name, category.IncomeExpense)
	
	if i.Amount <= 0 {
		return errors.New("transaction CheckValid: amount must be positive")
	}
	
	if i.IncomeExpense != category.IncomeExpense {
		log.Printf("[model]: IncomeExpense不匹配 - 交易: %s, 分类: %s", i.IncomeExpense, category.IncomeExpense)
		return errors.New("交易的收支类型与分类不匹配")
	}
	
	return nil
}
