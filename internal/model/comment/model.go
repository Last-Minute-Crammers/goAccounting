package commentModel

import (
	commonModel "goAccounting/internal/model/common"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	commonModel.BaseModel
	ID            uint           `gorm:"primarykey"`
	UserId        uint           `gorm:"comment:'评论者ID';index"`
	TargetUserId  uint           `gorm:"comment:'目标用户ID';index"`
	TransactionId uint           `gorm:"comment:'交易记录ID';index"`
	Content       string         `gorm:"comment:'评论内容';size:500"`
	CreatedAt     time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

func (c *Comment) TableName() string {
	return "comment"
}
