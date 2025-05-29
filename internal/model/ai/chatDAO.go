package aiModel

import (
	"gorm.io/gorm"
)

type ChatDAO struct {
	db *gorm.DB
}

func NewChatDAO() *ChatDAO {
	return &ChatDAO{}
}

func (dao *ChatDAO) Create(record *ChatRecord) error {
	return dao.db.Create(record).Error
}

func (dao *ChatDAO) GetByUserId(userId uint, limit int) ([]ChatRecord, error) {
	var records []ChatRecord
	err := dao.db.Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

func (dao *ChatDAO) DeleteByUserId(userId uint) error {
	return dao.db.Where("user_id = ?", userId).Delete(&ChatRecord{}).Error
}