package v1

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/global/db"
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	categoryModel "goAccounting/internal/model/category"
	transactionModel "goAccounting/internal/model/transaction"
	userModel "goAccounting/internal/model/user"
	"goAccounting/util/dataTool"
	"goAccounting/util/timeTool"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionApi struct{}

// GetOne
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		id			path		int	true	"Transaction ID"
//	@Success	200			{object}	response.Data{Data=response.TransactionDetail}
//	@Router		/user/transaction/{id} [get]
func (t *TransactionApi) GetOne(ctx *gin.Context) {
	trans, ok := contextFunc.GetTransByParam(ctx)
	if !ok {
		return
	}
	var data response.TransactionDetail
	err := data.SetData(trans)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(data, ctx)
}

// Update
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		accountId	path		int								true	"Account ID"
//	@Param		id			path		int								true	"Transaction ID"
//	@Param		body		body		request.TransactionUpdateOne	true	"Transaction data"
//	@Success	200			{object}	response.Data{Data=response.TransactionDetail}
//	@Router		/account/{accountId}/transaction/{id} [put]
func (t *TransactionApi) Update(ctx *gin.Context) {
	log.Println("[api]: get into TransactionUpdate")
	var requestData request.TransactionUpdateOne
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	oldTrans, ok := contextFunc.GetTransByParam(ctx)
	if !ok {
		return
	}
	trans := oldTrans
	trans.Info = transactionModel.Info{
		UserId:        oldTrans.UserId,
		CategoryId:    requestData.CategoryId,
		IncomeExpense: requestData.IncomeExpense,
		Amount:        requestData.Amount,
		Remark:        requestData.Remark,
		TradeTime:     requestData.TradeTime,
	}
	err := transactionService.Update(oldTrans.ID, trans.Info, ctx)
	if responseError(err, ctx) {
		return
	}

	var responseData response.TransactionDetail
	if err = responseData.SetData(trans); responseError(err, ctx) {
		return
	}
	response.OkWithData(responseData, ctx)
}

// Delete
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		accountId	path		int	true	"Account ID"
//	@Param		id			path		int	true	"Transaction ID"
//	@Success	200			{object}	response.NoContent
//	@Router		/account/{accountId}/transaction/{id} [delete]
func (t *TransactionApi) Delete(ctx *gin.Context) {
	log.Println("[api]: get into TxDelete")
	trans, pass := contextFunc.GetTransByParam(ctx)
	if !pass {
		return
	}
	err := transactionService.Delete(trans.ID, ctx)
	if err != nil {
		response.FailToError(ctx, err)
		return
	}
	response.Ok(ctx)
}

// GetList
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		body		body		request.TransactionGetList	true	"query parameters"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionDetail]{}}
//	@Router		/user/transaction/list [post]
func (t *TransactionApi) GetList(ctx *gin.Context) {
	log.Println("[api]: get into GetList")
	var requestData request.TransactionGetList
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	// 用工具函数获取 userId
	userId := contextFunc.GetUserId(ctx)

	condition := requestData.GetCondition()
	// 只查当前用户
	condition.UserId = userId
	log.Printf("[api]: TxList, userId is %d\n", userId)

	var transactionList []transactionModel.Transaction
	transactionList, err := transactionModel.NewDao().GetListByCondition(
		condition, requestData.Offset, requestData.Limit,
	)
	if responseError(err, ctx) {
		return
	}
	responseData := response.TransactionGetList{List: response.TransactionDetailList{}}
	err = responseData.List.SetData(transactionList)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(responseData, ctx)
}

// CreateOne
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		body		body		request.TransactionCreateOne	true	"transaction data"
//	@Success	200			{object}	response.Data{Data=response.TransactionDetail}
//	@Router		/user/transaction/ [post]
func (t *TransactionApi) CreateOne(ctx *gin.Context) {
	log.Println("[api]: get into CreateOne")
	var requestData request.TransactionCreateOne
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	// 用工具函数获取 userId
	userId := contextFunc.GetUserId(ctx)
	log.Printf("[txAPI]: userId is %d\n", userId)

	transInfo := transactionModel.Info{
		UserId:        userId, // 关键点
		CategoryId:    requestData.CategoryId,
		IncomeExpense: requestData.IncomeExpense,
		Amount:        requestData.Amount,
		Remark:        requestData.Remark,
		TradeTime:     requestData.TradeTime,
	}

	transaction, err := transactionService.Create(
		transInfo, transactionModel.RecordTypeOfManual, ctx,
	)
	if responseError(err, ctx) {
		return
	}

	var responseData response.TransactionDetail
	if err = responseData.SetData(transaction); responseError(err, ctx) {
		return
	}
	log.Println("[tx] : create sucess")
	response.OkWithData(responseData, ctx)
}

// GetMonthStatistic
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		body		body		request.TransactionMonthStatistic	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionStatistic]{}}
//	@Router		/user/transaction/statistic/month [post]
func (t *TransactionApi) GetMonthStatistic(ctx *gin.Context) {
	var requestData request.TransactionMonthStatistic
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		log.Printf("[api]: GetMonthStatistic bind error: %v", err)
		response.FailToParameter(ctx, err)
		return
	}

	userId := contextFunc.GetUserId(ctx)
	log.Printf("[api]: GetMonthStatistic, userId: %d, request: %+v", userId, requestData)

	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	requestData.SetLocal(time.Local)

	// 获取条件
	statisticCondition := requestData.GetStatisticCondition()
	extCond := requestData.GetExtensionCondition()

	// 确保设置用户ID
	statisticCondition.UserId = userId

	log.Printf("[api]: statisticCondition: %+v", statisticCondition)
	log.Printf("[api]: extCond: %+v", extCond)

	months := timeTool.SplitMonths(statisticCondition.StartTime, statisticCondition.EndTime)
	log.Printf("[api]: split months: %v", months)

	// 如果没有指定IncomeExpense，使用nil表示查询both
	ie := requestData.IncomeExpense
	if ie == nil {
		// 不设置ie，表示查询收入和支出
		log.Printf("[api]: No incomeExpense specified, will query both income and expense")
	}

	responseList := make([]response.TransactionStatistic, len(months))
	dao := transactionModel.NewStatisticDao(db.GetDb(ctx))

	for i := len(months) - 1; i >= 0; i-- {
		// 为每个月设置时间条件
		monthCondition := statisticCondition
		monthCondition.StartTime = months[i][0]
		monthCondition.EndTime = months[i][1]

		log.Printf("[api]: Processing month %d: %v to %v", i, monthCondition.StartTime, monthCondition.EndTime)

		monthStatistic, err := dao.GetIeStatisticByCondition(ie, monthCondition)
		if responseError(err, ctx) {
			return
		}

		log.Printf("[api]: month %d statistic: %+v", i, monthStatistic)

		responseList[i] = response.TransactionStatistic{
			IEStatistic: monthStatistic,
			StartTime:   monthCondition.StartTime,
			EndTime:     monthCondition.EndTime,
		}
	}

	log.Printf("[api]: final month response: %+v", responseList)
	response.OkWithData(response.List[response.TransactionStatistic]{List: responseList}, ctx)
}

// GetYearStatistic
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		body		body		request.TransactionYearStatistic	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionStatistic]{}}
//	@Router		/user/transaction/statistic/year [post]
func (t *TransactionApi) GetYearStatistic(ctx *gin.Context) {
	var requestData request.TransactionYearStatistic
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		log.Printf("[api]: GetYearStatistic bind error: %v", err)
		response.FailToParameter(ctx, err)
		return
	}

	userId := contextFunc.GetUserId(ctx)
	log.Printf("[api]: GetYearStatistic, userId: %d, request: %+v", userId, requestData)

	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	requestData.SetLocal(time.Local)

	// 获取条件
	statisticCondition := requestData.GetStatisticCondition()
	extCond := requestData.GetExtensionCondition()

	// 确保设置用户ID
	statisticCondition.UserId = userId

	log.Printf("[api]: statisticCondition: %+v", statisticCondition)
	log.Printf("[api]: extCond: %+v", extCond)

	years := timeTool.SplitYears(statisticCondition.StartTime, statisticCondition.EndTime)
	log.Printf("[api]: split years: %v", years)

	// 如果没有指定IncomeExpense，使用nil表示查询both
	ie := requestData.IncomeExpense
	if ie == nil {
		// 不设置ie，表示查询收入和支出
		log.Printf("[api]: No incomeExpense specified, will query both income and expense")
	}

	responseList := make([]response.TransactionStatistic, len(years))
	dao := transactionModel.NewStatisticDao(db.GetDb(ctx))

	for i := len(years) - 1; i >= 0; i-- {
		// 为每年设置时间条件
		yearCondition := statisticCondition
		yearCondition.StartTime = years[i][0]
		yearCondition.EndTime = years[i][1]

		log.Printf("[api]: Processing year %d: %v to %v", i, yearCondition.StartTime, yearCondition.EndTime)

		yearStatistic, err := dao.GetIeStatisticByCondition(ie, yearCondition)
		if responseError(err, ctx) {
			return
		}

		log.Printf("[api]: year %d statistic: %+v", i, yearStatistic)

		responseList[i] = response.TransactionStatistic{
			IEStatistic: yearStatistic,
			StartTime:   yearCondition.StartTime,
			EndTime:     yearCondition.EndTime,
		}
	}

	log.Printf("[api]: final year response: %+v", responseList)
	response.OkWithData(response.List[response.TransactionStatistic]{List: responseList}, ctx)
}

// GetWeekStatistic
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		body		body		request.TransactionWeekStatistic	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionStatistic]{}}
//	@Router		/user/transaction/statistic/week [post]
func (t *TransactionApi) GetWeekStatistic(ctx *gin.Context) {
	var requestData request.TransactionWeekStatistic
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		log.Printf("[api]: GetWeekStatistic bind error: %v", err)
		response.FailToParameter(ctx, err)
		return
	}

	userId := contextFunc.GetUserId(ctx)
	log.Printf("[api]: GetWeekStatistic, userId: %d, request: %+v", userId, requestData)

	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	requestData.SetLocal(time.Local)

	// 获取条件
	statisticCondition := requestData.GetStatisticCondition()
	extCond := requestData.GetExtensionCondition()

	// 确保设置用户ID
	statisticCondition.UserId = userId

	log.Printf("[api]: statisticCondition: %+v", statisticCondition)
	log.Printf("[api]: extCond: %+v", extCond)

	weeks := timeTool.SplitWeeks(statisticCondition.StartTime, statisticCondition.EndTime)
	log.Printf("[api]: split weeks: %v", weeks)

	// 如果没有指定IncomeExpense，使用nil表示查询both
	ie := requestData.IncomeExpense
	if ie == nil {
		// 不设置ie，表示查询收入和支出
		log.Printf("[api]: No incomeExpense specified, will query both income and expense")
	}

	responseList := make([]response.TransactionStatistic, len(weeks))
	dao := transactionModel.NewStatisticDao(db.GetDb(ctx))

	for i := len(weeks) - 1; i >= 0; i-- {
		// 为每周设置时间条件
		weekCondition := statisticCondition
		weekCondition.StartTime = weeks[i][0]
		weekCondition.EndTime = weeks[i][1]

		log.Printf("[api]: Processing week %d: %v to %v", i, weekCondition.StartTime, weekCondition.EndTime)

		weekStatistic, err := dao.GetIeStatisticByCondition(ie, weekCondition)
		if responseError(err, ctx) {
			return
		}

		log.Printf("[api]: week %d statistic: %+v", i, weekStatistic)

		responseList[i] = response.TransactionStatistic{
			IEStatistic: weekStatistic,
			StartTime:   weekCondition.StartTime,
			EndTime:     weekCondition.EndTime,
		}
	}

	log.Printf("[api]: final week response: %+v", responseList)
	response.OkWithData(response.List[response.TransactionStatistic]{List: responseList}, ctx)
}

// GetDayStatistic
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		body		body		request.TransactionDayStatistic	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionStatistic]{}}
//	@Router		/user/transaction/statistic/day [post]
func (t *TransactionApi) GetDayStatistic(ctx *gin.Context) {
	var requestData request.TransactionDayStatistic
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		log.Printf("[api]: GetDayStatistic bind error: %v", err)
		response.FailToParameter(ctx, err)
		return
	}

	userId := contextFunc.GetUserId(ctx)
	log.Printf("[api]: GetDayStatistic, userId: %d, request: %+v", userId, requestData)

	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	requestData.SetLocal(time.Local)

	// 获取条件
	statisticCondition := requestData.GetStatisticCondition()
	statisticCondition.UserId = userId

	days := timeTool.SplitDays(statisticCondition.StartTime, statisticCondition.EndTime)
	log.Printf("[api]: split days: %d days", len(days))

	responseList := make([]response.TransactionStatistic, len(days))
	dao := transactionModel.NewDao(db.GetDb(ctx))

	for i := len(days) - 1; i >= 0; i-- {
		dayStart := timeTool.GetFirstSecondOfDay(days[i])
		dayEnd := timeTool.GetLastSecondOfDay(days[i])

		log.Printf("[api]: Processing day %d: %v to %v", i, dayStart, dayEnd)

		// 创建查询条件
		condition := transactionModel.Condition{
			ForeignKeyCondition: transactionModel.ForeignKeyCondition{
				UserId: userId,
			},
			TimeCondition: transactionModel.TimeCondition{
				TradeStartTime: &dayStart,
				TradeEndTime:   &dayEnd,
			},
		}

		// 查询收入统计
		var incomeStats global.AmountCount
		if requestData.IncomeExpense == nil || *requestData.IncomeExpense == "income" {
			incomeCondition := condition
			incomeType := constant.Income
			incomeCondition.IncomeExpense = &incomeType

			incomeTransactions, err := dao.GetListByCondition(incomeCondition, 0, 1000)
			if err != nil {
				log.Printf("[api]: Error getting income transactions: %v", err)
			} else {
				for _, tx := range incomeTransactions {
					incomeStats.Amount += int64(tx.Amount) // 👈 类型转换
					incomeStats.Count++
				}
			}
		}

		// 查询支出统计
		var expenseStats global.AmountCount
		if requestData.IncomeExpense == nil || *requestData.IncomeExpense == "expense" {
			expenseCondition := condition
			expenseType := constant.Expense
			expenseCondition.IncomeExpense = &expenseType

			expenseTransactions, err := dao.GetListByCondition(expenseCondition, 0, 1000)
			if err != nil {
				log.Printf("[api]: Error getting expense transactions: %v", err)
			} else {
				for _, tx := range expenseTransactions {
					expenseStats.Amount += int64(tx.Amount) // 👈 类型转换
					expenseStats.Count++
				}
			}
		}

		log.Printf("[api]: day %d statistic - Income: %+v, Expense: %+v", i, incomeStats, expenseStats)

		responseList[i] = response.TransactionStatistic{
			IEStatistic: global.IEStatistic{
				Income:  incomeStats,
				Expense: expenseStats,
			},
			StartTime: dayStart,
			EndTime:   dayEnd,
		}
	}

	log.Printf("[api]: final day response: %d entries", len(responseList))
	response.OkWithData(response.List[response.TransactionStatistic]{List: responseList}, ctx)
}

// GetTotalStatistic
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Success	200			{object}	response.Data{Data=response.TransactionTotalStatistic}
//	@Router		/user/transaction/statistic/total [get]
func (t *TransactionApi) GetTotalStatistic(ctx *gin.Context) {
	userId := contextFunc.GetUserId(ctx)
	log.Printf("[api]: GetTotalStatistic, userId is %d", userId)

	dao := transactionModel.NewStatisticDao()
	totalStats, err := dao.GetTotalStatistics(userId)
	if responseError(err, ctx) {
		return
	}

	log.Printf("[api]: totalStats from DB: %+v", totalStats)

	// Calculate total assets (income - expense) and convert to int
	totalAssets := int(totalStats.Income.Amount - totalStats.Expense.Amount)
	log.Printf("[api]: calculated totalAssets: %d", totalAssets)

	responseData := response.TransactionTotalStatistic{
		IEStatistic: totalStats,
		TotalAssets: totalAssets,
	}

	log.Printf("[api]: final response data: %+v", responseData)
	response.OkWithData(responseData, ctx)
}

// GetCategoryAmountRank
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		accountId	path		int										true	"Account ID"
//	@Param		data		body		request.TransactionCategoryAmountRank	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionCategoryAmountRank]{}}
//	@Router		/account/{accountId}/transaction/category/amount/rank [get]
func (t *TransactionApi) GetCategoryAmountRank(ctx *gin.Context) {
	log.Println("[api]: get into func GetCategoryAmountRank")
	var requestData request.TransactionCategoryAmountRank
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	if err := requestData.CheckTimeFrame(); responseError(err, ctx) {
		return
	}
	UserId := contextFunc.GetUserId(ctx)
	log.Printf("[GetCategoryAmountRank]: userId = %d\n", UserId)
	// fetch ranking List
	var startTime, endTime = requestData.FormatDayTime()
	condition := transactionModel.CategoryAmountRankCondition{
		User:      userModel.User{ID: UserId},
		StartTime: startTime,
		EndTime:   endTime,
	}
	var err error
	var rankingList dataTool.Slice[uint, transactionModel.CategoryAmountRank]
	rankingList, err = transactionModel.NewStatisticDao().GetCategoryAmountRank(
		requestData.IncomeExpense, condition, requestData.Limit,
	)

	if responseError(err, ctx) {
		return
	}
	categoryIds := rankingList.ExtractValues(
		func(rank transactionModel.CategoryAmountRank) uint {
			return rank.CategoryId
		},
	)
	// fetch category
	var categoryList dataTool.Slice[uint, categoryModel.Category]
	err = db.Db.Where("id IN (?)", categoryIds).Find(&categoryList).Error
	if responseError(err, ctx) {
		return
	}
	categoryMap := categoryList.ToMap(
		func(category categoryModel.Category) uint {
			return category.ID
		},
	)
	// response
	responseData := make([]response.TransactionCategoryAmountRank, len(rankingList), len(rankingList))
	for i, rank := range rankingList {
		responseData[i].Amount = rank.Amount
		responseData[i].Count = rank.Count
		err = responseData[i].Category.SetData(categoryMap[rank.CategoryId])
		if responseError(err, ctx) {
			return
		}
	}

	categoryList = []categoryModel.Category{}
	query := db.Db.Where("user_id = ?", UserId)
	query = query.Where("income_expense = ?", requestData.IncomeExpense)
	if len(categoryIds) > 0 {
		query = query.Where("id NOT IN (?)", categoryIds)
	}
	err = query.Find(&categoryList).Error
	if responseError(err, ctx) {
		return
	}
	for _, category := range categoryList {
		responseCategory := response.TransactionCategoryAmountRank{}
		err = responseCategory.Category.SetData(category)
		if responseError(err, ctx) {
			return
		}
		responseData = append(responseData, responseCategory)
	}
	response.OkWithData(response.List[response.TransactionCategoryAmountRank]{List: responseData}, ctx)
}

// GetAmountRank
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		accountId	path		int								true	"Account ID"
//	@Param		data		body		request.TransactionAmountRank	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionDetailList]{}}
//	@Router		/account/{accountId}/transaction/amount/rank [get]
func (t *TransactionApi) GetAmountRank(ctx *gin.Context) {
	log.Println("[api]: get into func GetAmountRank")
	var requestData request.TransactionAmountRank
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	userId := contextFunc.GetUserId(ctx)
	log.Printf("[GetAmountRank]: userId = %d\n", userId)
	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	// fetch
	timeCond := transactionModel.NewTimeCondition()
	timeCond.SetTradeTimes(requestData.StartTime, requestData.EndTime)
	rankingList, err := transactionModel.NewDao().GetAmountRank(
		userId, requestData.IncomeExpense, *timeCond,
	)
	if responseError(err, ctx) {
		return
	}
	// response
	var responseList response.TransactionDetailList
	err = responseList.SetData(rankingList)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(response.List[response.TransactionDetail]{List: responseList}, ctx)
}
