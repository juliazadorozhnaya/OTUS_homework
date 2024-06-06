package model

type IUser interface {
	GetID() string
	GetFirstName() string
	GetLastName() string
	GetEmail() string
	GetAge() int64
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       int64  `json:"age"`
}

func (user *User) GetID() string {
	return user.ID
}

func (user *User) GetFirstName() string {
	return user.FirstName
}

func (user *User) GetLastName() string {
	return user.LastName
}

func (user *User) GetEmail() string {
	return user.Email
}

func (user *User) GetAge() int64 {
	return user.Age
}
