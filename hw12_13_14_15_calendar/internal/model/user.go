package model

// IUser интерфейс для структуры User, предоставляющий методы доступа к полям.
type IUser interface {
	GetID() string
	GetFirstName() string
	GetLastName() string
	GetEmail() string
	GetAge() int64
}

// User структура, представляющая пользователя.
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       int64  `json:"age"`
}

// GetID возвращает ID пользователя.
func (user *User) GetID() string {
	return user.ID
}

// GetFirstName возвращает имя пользователя.
func (user *User) GetFirstName() string {
	return user.FirstName
}

// GetLastName возвращает фамилию пользователя.
func (user *User) GetLastName() string {
	return user.LastName
}

// GetEmail возвращает email пользователя.
func (user *User) GetEmail() string {
	return user.Email
}

// GetAge возвращает возраст пользователя.
func (user *User) GetAge() int64 {
	return user.Age
}
