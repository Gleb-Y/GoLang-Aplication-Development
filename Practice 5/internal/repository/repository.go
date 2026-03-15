package repository

import (
	"prac5/internal/repository/_postgres"
	"prac5/internal/repository/_postgres/users"
	"prac5/pkg/modules"
)

type UserRepository interface {
	GetPaginatedUsers(page, pageSize int, filter modules.UserFilter) (modules.PaginatedResponse, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user modules.User) (int, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) error
	GetCommonFriends(userID1, userID2 int) ([]modules.User, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
