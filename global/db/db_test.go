package db

import (
	"context"
	"goAccounting/global/constant"
	"goAccounting/global/ctxutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初始化内存数据库
func initTestDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	return db
}

// 测试 GetDb 方法
func TestGetDb(t *testing.T) {
	// 初始化测试数据库
	testDb := initTestDb()

	// 测试从 context 提取 Db
	ctx := context.WithValue(context.Background(), constant.Db, testDb)
	db := GetDb(ctx)
	assert.Equal(t, testDb, db, "Expected to retrieve the database from context")

	// 测试默认 Db 返回
	db = GetDb(context.Background())
	assert.Equal(t, Db, db, "Expected to retrieve the default database")
}

// 测试 ExecInTransaction 方法的提交
func TestExecInTransaction_Commit(t *testing.T) {
	// 初始化测试数据库
	testDb := initTestDb()
	testDb.AutoMigrate(&TestModel{})

	ctx := ctxutil.WithDb(context.Background(), testDb)
	err := ExecInTransaction(ctx, func(ctx *ctxutil.TxContext) error {
		db := GetDb(ctx)
		return db.Create(&TestModel{Name: "Test"}).Error
	})

	assert.Nil(t, err, "Expected transaction to commit successfully")

	// 验证数据是否写入
	var count int64
	testDb.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(1), count, "Expected one record in the database")
}

// 测试 ExecInTransaction 方法的回滚
func TestExecInTransaction_Rollback(t *testing.T) {
	// 初始化测试数据库
	testDb := initTestDb()
	testDb.AutoMigrate(&TestModel{})

	ctx := ctxutil.WithDb(context.Background(), testDb)
	err := ExecInTransaction(ctx, func(ctx *ctxutil.TxContext) error {
		db := GetDb(ctx)
		if err := db.Create(&TestModel{Name: "Test"}).Error; err != nil {
			return err
		}
		return assert.AnError // 模拟错误，触发回滚
	})

	assert.NotNil(t, err, "Expected transaction to rollback due to error")

	// 验证数据是否回滚
	var count int64
	testDb.Model(&TestModel{}).Count(&count)
	assert.Equal(t, int64(0), count, "Expected no records in the database after rollback")
}

// 测试 AddCommitCallback
func TestAddCommitCallback(t *testing.T) {
	// 初始化测试数据库
	t.Log("Starting TestAddCommitCallback...")

	// 模拟事务上下文
	ctx := ctxutil.WithTxCommitContext(context.Background())
	callbackExecuted := false
	callback := func() {
		callbackExecuted = true
		t.Log("Callback executed")
	}

	err := AddCommitCallback(ctx, callback)
	if err != nil {
		t.Errorf("Failed to add commit callback: %v", err)
	} else {
		t.Log("Successfully added commit callback")
	}

	// 执行回调
	ctx.ExecCallback()
	if callbackExecuted {
		t.Log("Callback executed successfully")
	} else {
		t.Error("Callback was not executed")
	}
}

// 定义测试用的模型
type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}
