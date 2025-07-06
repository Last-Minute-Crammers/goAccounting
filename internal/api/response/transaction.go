package response

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/global/db"
	categoryModel "goAccounting/internal/model/category"
	transactionModel "goAccounting/internal/model/transaction"
	userModel "goAccounting/internal/model/user"
	"time"
)

type TransactionDetail struct {
	Id            uint
	UserId        uint
	UserName      string
	Amount        int
	CategoryId    uint
	CategoryIcon  string
	IncomeExpense constant.IncomeExpense
	Remark        string
	RecordType    transactionModel.RecordType
	TradeTime     time.Time
	UpdateTime    time.Time
	CreateTime    time.Time
}

func (t *TransactionDetail) SetData(
	trans transactionModel.Transaction,
) error {
	// 查询用户信息
	var user userModel.User
	if err := db.Db.First(&user, trans.UserId).Error; err != nil {
		// 如果用户不存在，使用默认值
		user = userModel.User{Username: "未知用户"}
	}

	// 查询分类信息
	var category categoryModel.Category
	if err := db.Db.First(&category, trans.CategoryId).Error; err != nil {
		// 如果分类不存在，使用默认值
		category = categoryModel.Category{
			Icon:          "unknown",
			IncomeExpense: trans.IncomeExpense, // 使用交易中的收支类型作为后备
		}
	}

	t.Id = trans.ID
	t.UserId = trans.UserId
	t.UserName = user.Username
	t.Amount = trans.Amount
	t.CategoryId = trans.CategoryId
	t.CategoryIcon = category.Icon
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
	
	// 收集所有需要的用户ID和分类ID
	userIds := make(map[uint]bool)
	categoryIds := make(map[uint]bool)
	for _, trans := range transList {
		userIds[trans.UserId] = true
		categoryIds[trans.CategoryId] = true
	}

	// 批量查询用户信息
	userMap := make(map[uint]userModel.User)
	if len(userIds) > 0 {
		var users []userModel.User
		if err := db.Db.Where("id IN ?", getKeys(userIds)).Find(&users).Error; err != nil {
			return err
		}
		for _, user := range users {
			userMap[user.ID] = user
		}
	}

	// 批量查询分类信息
	categoryMap := make(map[uint]categoryModel.Category)
	if len(categoryIds) > 0 {
		var categories []categoryModel.Category
		if err := db.Db.Where("id IN ?", getKeys(categoryIds)).Find(&categories).Error; err != nil {
			return err
		}
		for _, category := range categories {
			categoryMap[category.ID] = category
		}
	}

	// 设置响应数据
	for i, trans := range transList {
		user, userExists := userMap[trans.UserId]
		if !userExists {
			user = userModel.User{Username: "未知用户"}
		}
		
		category, categoryExists := categoryMap[trans.CategoryId]
		if !categoryExists {
			category = categoryModel.Category{
				Icon:          "unknown",
				IncomeExpense: trans.IncomeExpense, // 使用交易中的收支类型作为后备
			}
		}
		
		(*t)[i].Id = trans.ID
		(*t)[i].UserId = trans.UserId
		(*t)[i].UserName = user.Username
		(*t)[i].Amount = trans.Amount
		(*t)[i].CategoryId = trans.CategoryId
		(*t)[i].CategoryIcon = category.Icon
		(*t)[i].IncomeExpense = category.IncomeExpense
		(*t)[i].Remark = trans.Remark
		(*t)[i].RecordType = trans.RecordType
		(*t)[i].TradeTime = trans.TradeTime
		(*t)[i].UpdateTime = trans.UpdatedAt
		(*t)[i].CreateTime = trans.CreatedAt
	}
	return nil
}

// 辅助函数：从map中获取所有key
func getKeys(m map[uint]bool) []uint {
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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
