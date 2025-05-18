package userModel

import (
	"goAccounting/global"
	"time"

	"gorm.io/gorm"
)

// in this code, we directly use globalDb, just because this config only
// related to single table. so it doesn't have to use tx *gorm.DB which
// contains other transaction more information
type TransactionShareConfig struct {
	ID        uint           `gorm:"primarykey"`
	UserId    uint           `gorm:"unique"`
	IsShared  bool           `gorm:"comment:'whether is shared to frineds'"`
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

func (txc *TransactionShareConfig) IsRecordShared(userId uint) (bool, error) {
	var config TransactionShareConfig
	err := global.GlobalDb.Where("user_id = ?", userId).First(&config).Error
	if err != nil {
		return false, err
	}

	return config.IsShared, nil
}

func (txc *TransactionShareConfig) EnableSharing(userId uint) error {
	query := global.GlobalDb.Model(&TransactionShareConfig{}).Where("user_id = ?", userId)
	return query.Update("is_shared", true).Error
}

func (txc *TransactionShareConfig) DisableSharing(userId uint) error {
	query := global.GlobalDb.Model(&TransactionShareConfig{}).Where("user_id = ?", userId)
	return query.Update("is_shared", false).Error
}
