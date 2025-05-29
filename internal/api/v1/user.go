package v1

import (
	"errors"
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	transactionModel "goAccounting/internal/model/transaction"
	userModel "goAccounting/internal/model/user"
	"goAccounting/util/timeTool"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserApi struct {
}

type _userPublic interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)

	UpdatePassword(ctx *gin.Context)
}

type _userBase interface {
	Home(ctx *gin.Context)

	responseAndMaskUserInfo(userModel.UserInfo) response.UserInfo
	UpdatePassword(ctx *gin.Context)
	UpdateInfo(ctx *gin.Context)
}

type _userFriend interface {
	GetFriendList(ctx *gin.Context)
	responseUserFriendInvitation(userModel.FriendInvitation) (response.UserFriendInvitation, error)
	CreateFriendInvitation(ctx *gin.Context)
	AcceptFriendInvitation(ctx *gin.Context)
	RefuseFriendInvitation(ctx *gin.Context)

	GetFriendInvitationList(ctx *gin.Context)
	getFriendInvitationByParam(ctx *gin.Context) (userModel.FriendInvitation, bool)
}

// Login
//
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.UserLogin	true	"login data"
//	@Success	200		{object}	response.Data{Data=response.Login}
//	@Router		/public/user/login [post]
func (p *PublicApi) Login(ctx *gin.Context) {
	var requestData request.UserLogin
	var err error
	// handle error
	var loginFailResponseFunc = func() {
		if err != nil {
			key := global.Cache.GetKey(constant.LoginFailCount, requestData.Email)
			count, existCache := global.Cache.GetInt(key)
			if existCache {
				if count > 5 {
					response.FailToError(ctx, errors.New("错误次数过的，请稍后再试"))
					return
				} else {
					_ = global.Cache.Increment(key, 1)
				}
			} else {
				global.Cache.Set(key, 1, time.Hour*12)
			}
			response.FailToError(ctx, err)
			return
		}
	}
	defer loginFailResponseFunc()
	// check
	if err = ctx.ShouldBindJSON(&requestData); err != nil {
		return
	}
	// handler
	var user userModel.User
	var responseData response.Login
	var customClaims jwt.RegisteredClaims
	user, responseData.Token, customClaims, err = userService.Login(
		requestData.Email, requestData.Password, ctx,
	)
	if err != nil {
		err = errors.New("用户名不存在或者密码错误")
		return
	}
	responseData.TokenExpirationTime = customClaims.ExpiresAt.Time
	err = responseData.User.SetData(user)
	if err != nil {
		return
	}
	if responseData.Token == "" {
		err = errors.New("token获取失败")
		return
	}
	if err == nil {
		response.OkWithDetailed(responseData, "登录成功", ctx)
	}
}

// Register
//
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.UserRegister	true	"register data"
//	@Success	200		{object}	response.Data{Data=response.Login}
//	@Router		/public/user/register [post]
func (p *PublicApi) Register(ctx *gin.Context) {
	var requestData request.UserRegister
	var err error
	if err = ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	data := userModel.AddData{Username: requestData.Username, Password: requestData.Password, Email: requestData.Email}

	user, err := userService.Register(data, ctx)
	if responseError(err, ctx) {
		return
	}
	// 注册成功 获取token
	customClaims := commonService.MakeCustomClaims(user.ID)
	token, err := commonService.GenerateJWT(customClaims)
	if responseError(err, ctx) {
		return
	}

	responseData := response.Register{Token: token, TokenExpirationTime: customClaims.ExpiresAt.Time}
	err = responseData.User.SetData(user)
	if responseError(err, ctx) {
		return
	}
	response.OkWithDetailed(responseData, "注册成功", ctx)
}

// Home
//
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.UserHome	true	"data"
//	@Success	200		{object}	response.Data{Data=response.UserHome}
//	@Router		/user/home [get]
func (u *UserApi) Home(ctx *gin.Context) {
	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}
	now := time.Now()
	condition := transactionModel.StatisticCondition{
		ForeignKeyCondition: transactionModel.ForeignKeyCondition{
			UserId: user.ID,
		},
		StartTime: timeTool.GetFirstSecondOfMonth(now),
		EndTime:   timeTool.GetLastSecondOfMonth(now),
	}

	result, err := transactionModel.NewDao().GetIeStatisticByCondition(nil, condition, nil)
	if responseError(err, ctx) {
		return
	}
	responseData := gin.H{
		"month_amount": result.Income.Amount + result.Expense.Amount,
		// "income": result.Income.Amount,
		// "expense": result.Expense.Amount,
	}

	response.OkWithData(responseData, ctx)
}

// UpdateInfo 更新用户信息
//
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.UserUpdateInfo	true	"update data"
//	@Success	200		{object}	response.Data
//	@Router		/user/info [put]
func (u *UserApi) UpdateInfo(ctx *gin.Context) {
	var requestData request.UserUpdateInfo
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}

	// TODO: 实现更新用户信息逻辑
	// err = userService.UpdateInfo(user.ID, requestData)
	// if responseError(err, ctx) {
	//     return
	// }

	response.OkWithMessage("用户信息更新成功", ctx)
}

// UpdatePassword 更新密码
//
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.UserUpdatePassword	true	"password data"
//	@Success	200		{object}	response.Data
//	@Router		/user/password [put]
func (u *UserApi) UpdatePassword(ctx *gin.Context) {
	var requestData request.UserUpdatePassword
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}

	// TODO: 实现更新密码逻辑
	// err = userService.UpdatePassword(user.ID, requestData)
	// if responseError(err, ctx) {
	//     return
	// }

	response.OkWithMessage("密码更新成功", ctx)
}

// GetStats 获取用户统计数据
//
//	@Tags		User
//	@Produce	json
//	@Success	200	{object}	response.Data{Data=response.UserStats}
//	@Router		/user/stats [get]
func (u *UserApi) GetStats(ctx *gin.Context) {
	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}

	// TODO: 实现获取用户统计数据逻辑
	// stats, err := userService.GetUserStats(user.ID)
	// if responseError(err, ctx) {
	//     return
	// }

	// 模拟数据
	mockStats := gin.H{
		"total_income":        25000.00,
		"total_expense":       15600.00,
		"savings_rate":        37.6,
		"monthly_avg_expense": 2600.00,
	}

	response.OkWithData(mockStats, ctx)
}

// GetFriendList
//
//	@Tags		User/Friend
//	@Produce	json
//	@Success	200	{object}	response.Data{Data=response.List[response.UserInfo]{}}
//	@Router		/user/friend/list [get]
func (u *UserApi) GetFriendList(ctx *gin.Context) {
	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}
	var friendList []userModel.Friend
	friendList, err = userModel.NewDao().SelectFriendList(user.ID)
	if err != nil {
		log.Printf("err api:%v\n", err)
	}

	responseData := make([]response.UserInfo, len(friendList))
	var info userModel.UserInfo
	for i := 0; i < len(responseData); i++ {
		info, err = friendList[i].GetFriendInfo()
		if responseError(err, ctx) {
			return
		}
		responseData[i].SetData(info)

	}
	response.OkWithData(response.List[response.UserInfo]{List: responseData}, ctx)
}

func (u *UserApi) responseUserFriendInvitation(data userModel.FriendInvitation) (
	responseData response.UserFriendInvitation, err error,
) {
	var inviterInfo userModel.UserInfo
	var inviteeInfo userModel.UserInfo
	inviterInfo, err = data.GetInviterInfo()
	if err != nil {
		return
	}
	inviteeInfo, err = data.GetInviteeInfo()
	if err != nil {
		return
	}
	responseData = response.UserFriendInvitation{
		Id:         data.ID,
		CreateTime: data.CreatedAt,
	}
	responseData.Inviter.SetData(inviterInfo)
	responseData.Invitee.SetData(inviteeInfo) // Fix: was setting inviter twice
	return
}

// CreateFriendInvitation
//
//	@Tags		User/Friend/Invitation
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.UserCreateFriendInvitation	true	"data"
//	@Success	200		{object}	response.Data{Data=response.UserFriendInvitation}
//	@Router		/user/friend/invitation [post]
func (u *UserApi) CreateFriendInvitation(ctx *gin.Context) {
	var requestData request.UserCreateFriendInvitation
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}
	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}
	// 处理
	var invitation userModel.FriendInvitation
	var invitee userModel.User
	invitee, err = userModel.NewDao().SelectById(requestData.Invitee)
	if responseError(err, ctx) {
		return
	}
	invitation, err = userService.Friend.CreateInvitation(user, invitee, ctx)
	if responseError(err, ctx) {
		return
	}
	// 响应
	var responseData response.UserFriendInvitation
	responseData, err = u.responseUserFriendInvitation(invitation)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(responseData, ctx)
}

func (u *UserApi) getFriendInvitationByParam(ctx *gin.Context) (result userModel.FriendInvitation, isSuccess bool) {
	id, pass := contextFunc.GetUintParamByKey("id", ctx)
	if !pass {
		return
	}
	if pass, result = checkFunc.FriendInvitationBelongAndGet(id, ctx); !pass {
		return
	}
	isSuccess = true
	return
}

// AcceptFriendInvitation
//
//	@Tags		User/Friend/Invitation
//	@Produce	json
//	@Param		id	path		int	true	"Invitation ID"
//	@Success	200	{object}	response.Data{Data=response.UserFriendInvitation}
//	@Router		/user/friend/invitation/{id}/accept [put]
func (u *UserApi) AcceptFriendInvitation(ctx *gin.Context) {
	invitation, pass := u.getFriendInvitationByParam(ctx)
	if !pass {
		return
	}
	if invitation.Invitee != contextFunc.GetUserId(ctx) {
		response.FailToError(ctx, errors.New("非被邀请者！"))
		return
	}
	_, _, err := userService.Friend.AcceptInvitation(&invitation, ctx)
	if responseError(err, ctx) {
		return
	}
	var responseData response.UserFriendInvitation
	responseData, err = u.responseUserFriendInvitation(invitation)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(responseData, ctx)
}

// RefuseFriendInvitation
//
//	@Tags		User/Friend/Invitation
//	@Produce	json
//	@Param		id	path		int	true	"Invitation ID"
//	@Success	200	{object}	response.Data{Data=response.UserFriendInvitation}
//	@Router		/user/friend/invitation/{id}/refuse [put]
func (u *UserApi) RefuseFriendInvitation(ctx *gin.Context) {
	invitation, pass := u.getFriendInvitationByParam(ctx)
	if !pass {
		return
	}
	if invitation.Invitee != contextFunc.GetUserId(ctx) {
		response.FailToError(ctx, errors.New("非被邀请者！"))
		return
	}
	err := userService.Friend.RefuseInvitation(&invitation, ctx)
	if responseError(err, ctx) {
		return
	}
	var responseData response.UserFriendInvitation
	responseData, err = u.responseUserFriendInvitation(invitation)
	if responseError(err, ctx) {
		return
	}
	response.OkWithData(responseData, ctx)
}

// GetFriendInvitationList
//
//	@Tags		User/Friend/Invitation
//	@Produce	json
//	@Success	200	{object}	response.Data{Data=response.List[response.UserFriendInvitation]{}}
//	@Router		/user/friend/invitation [get]
func (u *UserApi) GetFriendInvitationList(ctx *gin.Context) {
	var requestData request.UserGetFriendInvitation
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response.FailToParameter(ctx, err)
		return
	}

	user, err := contextFunc.GetUser(ctx)
	if responseError(err, ctx) {
		return
	}
	var list []userModel.FriendInvitation
	if requestData.IsInvite {
		list, err = userModel.NewDao().SelectFriendInvitationList(&user.ID, nil)
		if err != nil {
			log.Printf("err api:%v\n", err)
		}
	} else {
		list, err = userModel.NewDao().SelectFriendInvitationList(nil, &user.ID)
		if err != nil {
			log.Printf("err api:%v\n", err)
		}

	}

	responseData := make([]response.UserFriendInvitation, len(list))
	for i := 0; i < len(responseData); i++ {
		responseData[i], err = u.responseUserFriendInvitation(list[i])
		if responseError(err, ctx) {
			return
		}
	}
	response.OkWithData(response.List[response.UserFriendInvitation]{List: responseData}, ctx)
}
