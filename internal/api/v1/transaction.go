package v1

import (
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	transactionModel "goAccounting/internal/model/transaction"
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

// GetList
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionDetail]{}}
//	@Router		/user/transaction/list [get]
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
//	@Param		data		body		request.TransactionMonthStatistic	true	"condition"
//	@Success	200			{object}	response.Data{Data=response.List[response.TransactionStatistic]{}}
//	@Router		/user/transaction/statistic/month [get]
func (t *TransactionApi) GetMonthStatistic(ctx *gin.Context) {
	var requestData request.TransactionMonthStatistic
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	requestData.SetLocal(time.Local)
	// condition
	statisticCondition, extCond := requestData.GetStatisticCondition(), requestData.GetExtensionCondition()
	condition := statisticCondition
	months := timeTool.SplitMonths(statisticCondition.StartTime, statisticCondition.EndTime)
	// select and process
	responseList := make([]response.TransactionStatistic, len(months))
	dao := transactionModel.NewDao()
	for i := len(months) - 1; i >= 0; i-- {
		condition.StartTime = months[i][0]
		condition.EndTime = months[i][1]

		monthStatistic, err := dao.GetIeStatisticByCondition(requestData.IncomeExpense, condition, &extCond)
		if responseError(err, ctx) {
			return
		}
		responseList[i] = response.TransactionStatistic{
			IEStatistic: monthStatistic,
			StartTime:   condition.StartTime,
			EndTime:     condition.EndTime,
		}
	}
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
	
	dao := transactionModel.NewStatisticDao()
	totalStats, err := dao.GetTotalStatistics(userId)
	if responseError(err, ctx) {
		return
	}
	
	// Calculate total assets (income - expense) and convert to int
	totalAssets := int(totalStats.Income.Amount - totalStats.Expense.Amount)
	
	responseData := response.TransactionTotalStatistic{
		IEStatistic: totalStats,
		TotalAssets: totalAssets,
	}
	
	response.OkWithData(responseData, ctx)
}
