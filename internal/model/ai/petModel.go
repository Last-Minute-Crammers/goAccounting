package aiModel

import (
	"time"
	"gorm.io/gorm"
)

type PetModel struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserId       uint           `gorm:"index;unique" json:"user_id"`
	Mood         string         `json:"mood"`
	ShortMessage string         `gorm:"type:text" json:"short_message"`
	LastUpdated  time.Time      `json:"last_updated"`
	CreatedAt    time.Time      `gorm:"type:TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type PetDAO struct {
	db *gorm.DB
}

func NewPetDAO() *PetDAO {
	return &PetDAO{}
}

func (dao *PetDAO) GetByUserId(userId uint) (*PetModel, error) {
	var pet PetModel
	err := dao.db.Where("user_id = ?", userId).First(&pet).Error
	if err == gorm.ErrRecordNotFound {
		// 创建默认宠物
		pet = PetModel{
			UserId:       userId,
			Mood:         "普通",
			ShortMessage: "欢迎来到理财世界！",
			LastUpdated:  time.Now(),
		}
		err = dao.db.Create(&pet).Error
	}
	return &pet, err
}

func (dao *PetDAO) UpdateMood(userId uint, mood, message string) error {
	return dao.db.Model(&PetModel{}).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"mood":         mood,
			"short_message": message,
			"last_updated": time.Now(),
		}).Error
}
