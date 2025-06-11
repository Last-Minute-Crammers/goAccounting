package response

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	categoryModel "goAccounting/internal/model/category"
	transactionModel "goAccounting/internal/model/transaction"
	userModel "goAccounting/internal/model/user"
	"time"
)

type TransactionDetail struct {
	Id                 uint
	UserId             uint
	UserName           string
	Amount             int
	CategoryId         uint
	CategoryName       string
	CategoryFatherName string
	IncomeExpense      constant.IncomeExpense
	Remark             string
	RecordType         transactionModel.RecordType
	TradeTime          time.Time
	UpdateTime         time.Time
	CreateTime         time.Time
	// CategoryIcon       string
}

func (t *TransactionDetail) SetData(
	trans transactionModel.Transaction,
) error {
	var (
		user     userModel.User
		category categoryModel.Category
	)

	t.Id = trans.ID
	t.UserId = user.ID
	t.UserName = user.Username
	t.Amount = trans.Amount
	t.CategoryId = trans.CategoryId
	t.CategoryName = category.Name
	t.IncomeExpense = category.IncomeExpense
	t.Remark = trans.Remark
	t.RecordType = trans.RecordType
	t.TradeTime = trans.TradeTime
	t.UpdateTime = trans.UpdatedAt
	t.CreateTime = trans.CreatedAt
	return nil
}

type TransactionDetailList []TransactionDetail

func (t *TransactionDetailList) SetData(transList []transactionModel.Transaction) error {
	*t = make([]TransactionDetail, len(transList))
	if len(transList) == 0 {
		return nil
	}
	for i, trans := range transList {
		(*t)[i].Id = trans.ID
		(*t)[i].Amount = trans.Amount
		(*t)[i].CategoryId = trans.CategoryId
		(*t)[i].IncomeExpense = trans.IncomeExpense
		(*t)[i].Remark = trans.Remark
		(*t)[i].RecordType = trans.RecordType
		(*t)[i].TradeTime = trans.TradeTime
		(*t)[i].UpdateTime = trans.UpdatedAt
		(*t)[i].CreateTime = trans.CreatedAt
	}
	return nil
}

type TransactionGetList struct {
	List TransactionDetailList
	PageData
}

type TransactionTotal struct {
	global.IEStatistic
}

type TransactionStatistic struct {
	global.IEStatistic
	StartTime time.Time
	EndTime   time.Time
}

type TransactionDayStatistic struct {
	global.AmountCount
	Date time.Time
}

type TransactionCategoryAmountRank struct {
	Category CategoryOne
	global.AmountCount
}

type TransactionTotalStatistic struct {
	global.IEStatistic
	TotalAssets int `json:"total_assets"`
}
