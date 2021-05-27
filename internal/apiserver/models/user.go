package models

type User struct {
	ID    uint32 `json:"id"`
	Age   uint8  `json:"age"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
