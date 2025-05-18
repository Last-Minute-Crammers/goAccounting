package userModel

// different from dao, this is specific user model's operation
import (
	"goAccounting/global"
	commonModel "goAccounting/internal/model/common"
	"time"

	"gorm.io/gorm"
)

type User struct {
	commonModel.BaseModel
	ID        uint           `gorm:"primarykey"`
	Username  string         `gorm:"type:varchar(128)"`
	Password  string         `gorm:"type:varchar(64)"`
	Email     string         `gorm:"type:varchar(64);unique"`
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

// provide to client, for lighter goal
type UserInfo struct {
	ID       uint
	Username string
	Email    string
}

func (u *User) SelectById(id uint, selects ...any) error {
	// dbQuery means SELECT * FROM users(tablename) WHERE id = <id> in SQL;
	dbQuery := global.GlobalDb.Where("id = ?", id)
	if len(selects) > 0 {
		dbQuery = dbQuery.Select(selects[0], selects[1:]...)
		// equals to
		// fields := make([]any, len(selects))
		// copy(fileds, selects)
		// firstField := fields[0] -- which match parm query
		// remainningFields := fields[1:] -- which match parm args
		// dbQuery = dbQuery.Select(firstField, remainingFields...)
	}
	return dbQuery.First(u).Error
}

// TOUR func is low priority to develop
func (u *User) IsTourist(db *gorm.DB) (bool, error) {
	return true, nil
}
