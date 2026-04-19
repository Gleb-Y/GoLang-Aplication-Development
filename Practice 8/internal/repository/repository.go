package repository

import (
	"practice-8/internal/repository/_postgres"
	"practice-8/pkg/modules"
)

type UserRepository interface {
	CreateUser(username, email, hashedPassword, role string) (*modules.User, error)
	GetUserByUsername(username string) (*modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	PromoteUserToAdmin(id int) error
	UpdateUser(id int, username, email string) error
	DeleteUser(id int) error
	ListUsers() ([]*modules.User, error)
}

type Repositories struct {
	UserRepository UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: _postgres.NewUserRepository(db),
	}
}
