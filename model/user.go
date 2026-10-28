package model
type User struct{
	ID   int
	Username string
	Password string
	Nickname string
	Role int
} 

func (u *User) TableName() string{
	return  "user"
}