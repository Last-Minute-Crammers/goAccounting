package userModel

// not necessary, so it's not developed now
import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"time"

	"gorm.io/gorm"
)

type Log struct {
	ID        uint      `gorm:"primarykey"`
	UserId    uint      `gorm:"comment:用户id;not null"`
	Action    string    `gorm:"comment:操作;not null;size:32"` // 简化为普通字符串
	CreatedAt time.Time `gorm:"type:TIMESTAMP"`
}

func (l *Log) TableName() string {
	return "user_log"
}

type LogDao struct {
	db *gorm.DB
}

func NewLogDao(db *gorm.DB) *LogDao {
	if db == nil {
		db = global.GlobalDb
	}
	return &LogDao{db}
}

type LogAddData struct {
	Action constant.UserAction
	Remark string
}

func (l *LogDao) Add(user User, data *LogAddData) (*Log, error) {
	log := &Log{
		UserId: user.ID,
		Action: string(data.Action),
	}
	err := l.db.Create(&log).Error
	return log, err
}
