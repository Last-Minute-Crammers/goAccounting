package response

import (
	"goAccounting/global/constant"
	categoryModel "goAccounting/internal/model/category"
	"goAccounting/util/dataTool"
)

type CategoryOne struct {
	Id            uint
	Name          string
	IncomeExpense constant.IncomeExpense
}

func (co *CategoryOne) SetData(category categoryModel.Category) error {
	co.Id = category.ID
	co.Name = category.Name
	co.IncomeExpense = category.IncomeExpense
	return nil
}

type CategoryDetailList []CategoryOne

func (cdl *CategoryDetailList) SetData(categoryList dataTool.Slice[uint, categoryModel.Category]) error {
	*cdl = make(CategoryDetailList, len(categoryList))
	if len(categoryList) == 0 {
		return nil
	}
	for i, category := range categoryList {
		err := (*cdl)[i].SetData(category)
		if err != nil {
			return err
		}
	}
	return nil
}
