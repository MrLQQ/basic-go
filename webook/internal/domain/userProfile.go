package domain

type UserProfile struct {
	Id       int64
	User_id  int64
	Nickname string
	Birthday string
	About_me string
}

func (u UserProfile) ValidateEmail() bool {
	// 在这里用正则表达式校验
	return true
}
