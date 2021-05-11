package adapters

import (
	"errors"
	"github.com/astanishevskyi/http-server/apiserver/models"
	"math/rand"
)

type UserStorage struct {
	UsersSlice []models.User `json:" users"`
}

func NewInMemoryUserStorage() *UserStorage {
	storage := UserStorage{UsersSlice: make([]models.User, 0)}
	return &storage
}

func (s *UserStorage) GetAll() []models.User {
	return s.UsersSlice
}

func (s *UserStorage) Retrieve(id uint32) (models.User, error) {
	for _, val := range s.UsersSlice {
		if val.Id == id {
			return val, nil
		}
	}
	return models.User{}, errors.New("no user was found")
}

func (s *UserStorage) Add(name, email string, age uint8) models.User {
	user := models.User{Id: rand.Uint32(), Name: name, Email: email, Age: age}
	s.UsersSlice = append(s.UsersSlice, user)
	return user
}

func (s *UserStorage) Remove(id uint32) (uint32, error) {
	for i, val := range s.UsersSlice {
		if val.Id == id {
			s.UsersSlice = append(s.UsersSlice[:i], s.UsersSlice[i+1:]...)
			return id, nil
		}
	}
	return 0, errors.New("no user was found")
}

func (s *UserStorage) Update(id uint32, name, email string, age uint8) (models.User, error) {
	for i, val := range s.UsersSlice {
		if val.Id == id {
			val.Name = name
			val.Email = email
			val.Age = age
			s.UsersSlice[i] = val
			return val, nil
		}
	}
	return models.User{}, errors.New("no user was found")

}
