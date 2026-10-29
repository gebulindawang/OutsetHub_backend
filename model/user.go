package model
type User struct{
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`          // 密码哈希不再返回前端
	Nickname string `json:"nickname"`
	Role     int    `json:"role"`
}

func (u *User) TableName() string{
	return  "user"
}