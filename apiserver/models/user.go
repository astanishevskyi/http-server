package models

type User struct {
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint8  `json:"age"`
}

type UserInterface interface {
	GetAll() []User
	Retrieve(id uint32) (User, error)
	Add(name, email string, age uint8) User
	Remove(id uint32) (uint32, error)
	Update(id uint32, name, email string, age uint8) (User, error)
}
