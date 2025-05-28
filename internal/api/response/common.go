package response

import (
	"time"
)

type CommonCaptcha struct {
	CaptchaId     string
	PicBase64     string
	CaptchaLength int
	OpenCaptcha   bool
}

type Id struct {
	Id uint
}

type CreateResponse struct {
	Id        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Token struct {
	Token               string
	TokenExpirationTime time.Time
}

type TwoLevelTree struct {
	Tree []Father
}

type Father struct {
	NameId
	Children []NameId
}

type NameId struct {
	Id   uint
	Name string
}

type NameValue struct {
	Name  string
	Value int
}

type PageData struct {
	page  int
	limit int
	count int
}

type ExpirationTime struct {
	ExpirationTime int
}

type List[T any] struct {
	List []T
}
