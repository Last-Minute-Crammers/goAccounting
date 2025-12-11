package request

type Pagination struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type UserQueryFilter struct {
	Pagination
	Email    string `form:"email" json:"email"`
	Username string `form:"username" json:"username"`
	IsAdmin  *bool  `form:"isAdmin" json:"isAdmin"` // 使用指针区分未设置和 false
}
