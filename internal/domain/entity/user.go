package entity

type User struct {
	UserName string
	Password string
}

func NewUser(userName, password string) *User {
	return &User{
		UserName: userName,
		Password: password,
	}
}
