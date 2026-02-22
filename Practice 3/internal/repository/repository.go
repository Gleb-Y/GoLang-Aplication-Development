package repository

import (
	"prac3/internal/repository/_postgres"
	"prac3/internal/repository/_postgres/users"
	"prac3/pkg/modules"
)

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user modules.User) (int, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) error
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
