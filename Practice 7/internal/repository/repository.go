package repository

import (
	"practice-7/internal/repository/_postgres"
	"practice-7/pkg/modules"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	CreateUser(username, email, hashedPassword, role string) (*modules.User, error)
	GetUserByUsername(username string) (*modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	PromoteUserToAdmin(id int) error
}

type Repositories struct {
	UserRepository UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: _postgres.NewUserRepository(db),
	}
}
