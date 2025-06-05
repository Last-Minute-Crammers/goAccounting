package util

import (
	"goAccounting/global/constant"
	"goAccounting/global/db"
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	transactionModel "goAccounting/internal/model/transaction"
	userModel "goAccounting/internal/model/user"
	"log"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ContextFunc = new(contextFunc)

type contextFunc struct{}

func (cf *contextFunc) GetToken(ctx *gin.Context) string {
	return ctx.Request.Header.Get("authorization")
}

func (cf *contextFunc) SetClaims(claims jwt.RegisteredClaims, ctx *gin.Context) {
	ctx.Set(string(constant.Claims), claims)
}

func (cf *contextFunc) GetClaims(ctx *gin.Context) jwt.RegisteredClaims {
	return ctx.Value(string(constant.Claims)).(jwt.RegisteredClaims)
}

// insert into other ctx
func (cf *contextFunc) SetUserId(id uint, ctx *gin.Context) {
	ctx.Set(string(constant.UserId), id)
}

func (cf *contextFunc) GetUserId(ctx *gin.Context) uint {
	return cf.GetUint(constant.UserId, ctx)
}

func (cf *contextFunc) GetUser(ctx *gin.Context) (userModel.User, error) {
	value, exits := ctx.Get(string(constant.User))
	if exits {
		return value.(userModel.User), nil
	}
	var user userModel.User
	err := db.Db.First(&user, cf.GetUserId(ctx)).Error
	ctx.Set(string(constant.User), user)
	return user, err
}

func (cf *contextFunc) GetId(ctx *gin.Context) uint {
	return cf.GetUint("id", ctx)
}

func (cf *contextFunc) GetUintParamByKey(key string, ctx *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(ctx.Param(key))
	if err != nil {
		response.FailToParameter(ctx, err)
		return 0, false
	}
	return uint(id), true
}

func (cf *contextFunc) GetInfoTypeFormParam(ctx *gin.Context) request.InfoType {
	return request.InfoType(ctx.Param("type"))
}

func (cf *contextFunc) GetParamId(ctx *gin.Context) (uint, bool) {
	return cf.GetUintParamByKey("id", ctx)
}

func (cf *contextFunc) GetUint(key constant.Key, ctx *gin.Context) uint {
	param := ctx.Param(string(key))
	if len(param) != 0 {
		id, err := strconv.Atoi(param)
		if err != nil {
			log.Printf("[GetUint] ctx.Param(%v) parse error: %v", key, err)
			panic(err)
		}
		return uint(id)
	}
	value := ctx.Value(string(key))
	log.Printf("[GetUint] ctx.Value(%v) = %#v, type: %T", key, value, value)
	switch v := value.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("[GetUint] string to int error: %v", err)
			panic(err)
		}
		return uint(id)
	case nil:
		log.Printf("[GetUint] value is nil for key: %v", key)
		panic("userId not found in context")
	default:
		log.Printf("[GetUint] unknown value type for key %v: %v (%v)", key, v, reflect.TypeOf(v))
		panic("userId is not a recognized type")
	}
}

func (cf *contextFunc) GetTransByParam(ctx *gin.Context) (transactionModel.Transaction, bool) {
	id, ok := cf.GetParamId(ctx)
	if !ok {
		return transactionModel.Transaction{}, false
	}
	tx, err := transactionModel.NewDao().SelectById(id, false)
	if err != nil {
		response.FailToError(ctx, err)
		return transactionModel.Transaction{}, false
	}
	return tx, true
}
