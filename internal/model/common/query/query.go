package query

import (
	"goAccounting/global"
	commonModel "goAccounting/internal/model/common"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func FirstByPrimaryKey[T commonModel.Model](key interface{}) (T, error) {
	var result T
	// gorm will automatically infer whether the table can match this model
	err := global.GlobalDb.First(&result, key).Error
	return result, err
}

func FirstByField[T commonModel.Model](field string, value interface{}) (T, error) {
	var result T
	dbQuery := global.GlobalDb.Where(map[string]interface{}{field: value})
	// the next line of code is the real statement that executes the query
	err := dbQuery.First(&result).Error
	return result, err

}

func ExistOfResult[T commonModel.Model](query interface{}, args ...interface{}) (bool, error) {
	var result T
	dbQuery := global.GlobalDb.Where(query, args...)
	err := dbQuery.First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil

	} else if err == nil {
		return true, nil
	} else {
		return false, errors.Wrap(err, "existOfQueryResult")
	}
}
