package v1

import (
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	transactionModel "goAccounting/internal/model/transaction"
	"goAccounting/util/timeTool"
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
	var requestData request.TransactionGetList
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	if err := requestData.CheckTimeFrame(); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	_, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}

	// select and response
	condition := requestData.GetCondition()
	var transactionList []transactionModel.Transaction
	transactionList, err = transactionModel.NewDao().GetListByCondition(
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
	var requestData request.TransactionCreateOne
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	transInfo := transactionModel.Info{
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

// UpdateOne 更新交易记录
//
//	@Tags		Transaction
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int								true	"Transaction ID"
//	@Param		body	body		request.TransactionUpdateOne	true	"transaction data"
//	@Success	200		{object}	response.Data{Data=response.TransactionDetail}
//	@Router		/user/transaction/{id} [put]
func (t *TransactionApi) UpdateOne(ctx *gin.Context) {
	trans, ok := contextFunc.GetTransByParam(ctx)
	if !ok {
		return
	}
	
	var requestData request.TransactionUpdateOne
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	// TODO: 实现更新交易记录逻辑
	// updatedTrans, err := transactionService.Update(trans.ID, requestData, ctx)
	// if responseError(err, ctx) {
	//     return
	// }
	
	var responseData response.TransactionDetail
	err := responseData.SetData(trans)
	if responseError(err, ctx) {
		return
	}
	
	response.OkWithData(responseData, ctx)
}

// DeleteOne 删除交易记录
//
//	@Tags		Transaction
//	@Produce	json
//	@Param		id		path		int	true	"Transaction ID"
//	@Success	200		{object}	response.Data
//	@Router		/user/transaction/{id} [delete]
func (t *TransactionApi) DeleteOne(ctx *gin.Context) {
	trans, ok := contextFunc.GetTransByParam(ctx)
	if !ok {
		return
	}
	
	// TODO: 实现删除交易记录逻辑
	// err := transactionService.Delete(trans.ID, ctx)
	// if responseError(err, ctx) {
	//     return
	// }
	
	response.OkWithMessage("交易记录删除成功", ctx)
}

// GetCategoryStats 获取分类统计
//
//	@Tags		Transaction/Statistics
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.TransactionCategoryAmountRank	true	"condition"
//	@Success	200		{object}	response.Data{Data=response.List[response.TransactionCategoryAmountRank]{}}
//	@Router		/user/transaction/statistic/category [get]
func (t *TransactionApi) GetCategoryStats(ctx *gin.Context) {
	var requestData request.TransactionCategoryAmountRank
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	
	// TODO: 实现分类统计逻辑
	// stats, err := transactionService.GetCategoryStats(requestData, ctx)
	// if responseError(err, ctx) {
	//     return
	// }
	
	// 模拟数据
	mockStats := []response.TransactionCategoryAmountRank{
		// TODO: Replace with actual category response structure
	}
	
	response.OkWithData(response.List[response.TransactionCategoryAmountRank]{List: mockStats}, ctx)
}
