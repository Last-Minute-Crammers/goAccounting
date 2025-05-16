package commonModel

import (
	"goAccounting/global"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BaseModel struct {
}

func (base *BaseModel) modelInterface() {
}

type Model interface {
	modelInterface()
}

func ExistOfModel(model Model, query interface{}, args ...interface{}) (bool, error) {
	err := global.GlobalDb.Where(query, args...).Take(model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil

	} else if err == nil {
		return true, nil

	} else {
		return false, errors.Wrap(err, "existOfModel")
	}
}
