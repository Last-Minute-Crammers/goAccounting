package userModel

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Match your application configuration
		},
	})
	assert.NoError(t, err)

	return gormDB, mock
}

// CheckTableExists: 检查表是否存在的函数
func CheckTableExists(db *gorm.DB, tableName string) (bool, error) {
	var result string
	err := db.Raw("SHOW TABLES LIKE ?", tableName).Scan(&result).Error
	if err != nil {
		return false, err
	}
	return result == tableName, nil
}

// TestCheckTableExists: 测试 CheckTableExists 函数
func TestCheckTableExists(t *testing.T) {
	db, mock := setupMockDB(t)

	// 模拟查询表是否存在
	mock.ExpectQuery(regexp.QuoteMeta("SHOW TABLES LIKE ?")).
		WithArgs("user").
		WillReturnRows(sqlmock.NewRows([]string{"Tables_in_test"}).AddRow("user"))

	// 执行测试
	exists, err := CheckTableExists(db, "user")

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, exists)

	// 验证所有预期已满足
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSelectByIdWithMock(t *testing.T) {
	db, mock := setupMockDB(t)
	db.Logger = logger.Default.LogMode(logger.Info) // 测试前添加
	dao := NewDao(db)

	// 创建时间值
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	updatedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	// 设置预期查询和结果
	rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "test_user", "hashed_password", "test@example.com", createdAt, updatedAt, nil)

	fmt.Println("start to mockExecQuery")
	// Changed from `users` to `user` to match SingularTable: true configuration
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE Id = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT ?")).
		WithArgs(1, 1). // 第二个参数为 LIMIT 的值
		WillReturnRows(rows)
	fmt.Println("end mockExecQuery")
	// 执行测试

	fmt.Println("start to Exec dao's func")
	user, err := dao.SelectById(1)
	if err != nil {
		fmt.Printf("we got err: %v\n", err)
	} else {
		fmt.Println("err is nil")
	}

	fmt.Printf("according to dao, user.Id is %d\n", user.ID)
	assert.NoError(t, err)
	fmt.Println("end  Exec dao's func")
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "test_user", user.Username)

	// 验证所有预期已满足
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
