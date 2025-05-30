package response

import (
	"goAccounting/global/constant"
	categoryModel "goAccounting/internal/model/category"
	"time"
)

type CategoryOne struct {
	ID            uint                       `json:"id"`
	Name          string                     `json:"name"`
	Icon          string                     `json:"icon"`
	IncomeExpense constant.IncomeExpense     `json:"income_expense"`
	UserId        *uint                      `json:"user_id,omitempty"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
}

func (c *CategoryOne) SetData(category categoryModel.Category) {
	c.ID = category.ID
	c.Name = category.Name
	c.Icon = category.Icon
	c.IncomeExpense = category.IncomeExpense
	c.UserId = category.UserId
	c.CreatedAt = category.CreatedAt
	c.UpdatedAt = category.UpdatedAt
}

type CategoryDetailList []CategoryOne

func (c *CategoryDetailList) SetData(categories []categoryModel.Category) error {
	*c = make([]CategoryOne, len(categories))
	for i, category := range categories {
		(*c)[i].SetData(category)
	}
	return nil
}
