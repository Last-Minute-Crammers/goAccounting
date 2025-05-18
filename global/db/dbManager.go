package db

import (
	"context"
	"fmt"
	"goAccounting/global/constant"
	"goAccounting/global/ctxutil"
	"goAccounting/initialize"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db = initialize.Database
	// we don't want some operation influence the log, so it hvae a silent
	// charater compared with GlobalDb. so we use it for init.
	InitDb  = Db.Session(&gorm.Session{Logger: Db.Logger.LogMode(logger.Silent)})
	Context *ctxutil.DbContext
)

func init() {
	Context = ctxutil.WithDb(context.Background(), Db)
}

func GetDb(ctx context.Context) *gorm.DB {
	value := ctx.Value(constant.Db)
	if value == nil {
		return Db
	}
	return value.(*gorm.DB)
}

type TxFunc func(ctx *ctxutil.TxContext) error

func ExecInTransaction(parent context.Context, fc TxFunc) error {
	ctx := ctxutil.WithTxCommitContext(parent)
	err := GetDb(ctx).Transaction(
		func(tx *gorm.DB) error {
			return fc(ctxutil.WithTx(ctx, tx))
		},
	)
	if err == nil {
		ctx.ExecCallback()
	}
	return err
}

func AddCommitCallback(parent context.Context, callbacks ...ctxutil.TxCommitCallback) error {
	value := parent.Value(constant.TxCommit)
	txCommitCtx, ok := value.(*ctxutil.TxCommitContext)
	if !ok {
		return fmt.Errorf("context doer not contain a valid TxCommitContext")
	}
	txCommitCtx.AddCallBack(callbacks...)
	return nil
}
