package ctxutil

import (
	"context"
	"fmt"
	constant "goAccounting/global/constant"
	"log"

	"gorm.io/gorm"
)

type DbContext struct {
	context.Context
	db *gorm.DB
}

func WithDb(parent context.Context, db *gorm.DB) *DbContext {
	return &DbContext{Context: parent, db: db}
}

func (dc *DbContext) Value(key any) any {
	if key == constant.Db {
		return dc.db
	}
	return dc.Context.Value(key)
}

func (dc *DbContext) GetDb() *gorm.DB {
	return dc.db
}

type TxContext struct {
	context.Context
	txdb *gorm.DB
}

func WithTx(parent context.Context, tx *gorm.DB) *TxContext {
	return &TxContext{Context: parent, txdb: tx}
}

func (tc *TxContext) Value(key any) any {
	if key == constant.Db || key == constant.Tx {
		return tc.txdb
	}
	return tc.Context.Value(key)
}

func (tc *TxContext) GetDb() *gorm.DB {
	return tc.txdb
}

type TxCommitCallback func()
type TxCommitContext struct {
	context.Context
	callbacks []TxCommitCallback
}

func WithTxCommitContext(parent context.Context) *TxCommitContext {
	return &TxCommitContext{Context: parent}
}

func (t *TxCommitContext) Value(key any) any {
	if key == constant.TxCommit {
		return t
	}
	return t.Context.Value(key)
}

// it seems actually don't have errors, so I simplify it into void
func (t *TxCommitContext) AddCallBack(callbacks ...TxCommitCallback) {
	t.callbacks = append(t.callbacks, callbacks...)
}

func (t *TxCommitContext) ExecCallback() {
	if len(t.callbacks) == 0 {
		log.Println("there's no function need to callback")
		return
	}
	parent, ok := t.Context.Value(constant.TxCommit).(*TxCommitContext)
	if ok && parent != nil && parent != t {
		parent.AddCallBack(t.callbacks...)
		fmt.Println("send callbacks to parent TxCommit")
		return
	}

	if !ok {
		log.Println("parent don't have TxCommitContext")
	}
	for _, callback := range t.callbacks {
		callback()
	}
}
