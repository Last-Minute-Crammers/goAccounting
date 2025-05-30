package service

import (
	categoryService "goAccounting/internal/service/category"
	commonService "goAccounting/internal/service/common"
	thirdpartyService "goAccounting/internal/service/thirdparty"
	transactionService "goAccounting/internal/service/transaction"
	userService "goAccounting/internal/service/user"
)

type Group struct {
	CommonServiceGroup      commonService.Group
	CategoryServiceGroup    categoryService.Group
	TransactionServiceGroup transactionService.Group
	UserServiceGroup        userService.Group
	ThirdpartyServiceGroup  thirdpartyService.Group
}

var GroupApp = new(Group)
