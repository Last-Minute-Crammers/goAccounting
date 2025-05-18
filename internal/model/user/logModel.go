package userModel

// not necessary, so it's not developed now
import (
	"time"
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
