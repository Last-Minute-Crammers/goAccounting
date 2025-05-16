package ctxutil

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// time.Time is also a Struct, so it can have methods
// Time 是普通的时间戳字段用于记录创建时间、更新时间等普通时间信息
// 行为上就是普通的时间字段，存储具体的时间值
type Time time.Time

func (Time) GormDataType() string {
	return "timestamp"
}

func (Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "TIMESTAMP"
}

// 特殊的标记字段，专用于 GORM 的软删除功能
// 在正常状态下为 NUL
// 当记录被"软删除"时，会被设置为删除的时间戳
// GORM 会自动排除已"软删除"的记录（除非明确要求包含）
type DeleteAt gorm.DeletedAt

func (DeleteAt) GormDataType() string {
	return "timestamp"
}

func (DeleteAt) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "TIMESTAMP"
}
