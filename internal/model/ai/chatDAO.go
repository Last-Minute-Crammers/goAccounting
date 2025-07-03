package aiModel

import (
	"goAccounting/global"
	"time"
	"gorm.io/gorm"
)

type ChatDAO struct {
	db *gorm.DB
}

func NewChatDAO(db ...*gorm.DB) *ChatDAO {
	if len(db) > 0 {
		return &ChatDAO{db: db[0]}
	}
	return &ChatDAO{db: global.GlobalDb}
}

// Create 创建聊天记录
func (dao *ChatDAO) Create(record *ChatRecord) error {
	return dao.db.Create(record).Error
}

// GetBySessionId 根据会话ID获取聊天记录
func (dao *ChatDAO) GetBySessionId(sessionId string, limit int) ([]ChatRecord, error) {
	var records []ChatRecord
	err := dao.db.Where("session_id = ?", sessionId).
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetByUserId 根据用户ID获取聊天记录
func (dao *ChatDAO) GetByUserId(userId uint, offset, limit int) ([]ChatRecord, error) {
	var records []ChatRecord
	err := dao.db.Where("user_id = ?", userId).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetByRequestId 根据问答ID获取特定记录
func (dao *ChatDAO) GetByRequestId(requestId string) (*ChatRecord, error) {
	var record ChatRecord
	err := dao.db.Where("request_id = ?", requestId).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// DeleteBySessionId 删除指定会话的所有记录
func (dao *ChatDAO) DeleteBySessionId(sessionId string) error {
	return dao.db.Where("session_id = ?", sessionId).Delete(&ChatRecord{}).Error
}

// GetRecentByUserId 获取用户最近的聊天记录
func (dao *ChatDAO) GetRecentByUserId(userId uint, days int, limit int) ([]ChatRecord, error) {
	var records []ChatRecord
	startTime := time.Now().AddDate(0, 0, -days)
	err := dao.db.Where("user_id = ? AND created_at >= ?", userId, startTime).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}
