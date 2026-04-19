package repository

import "practice-8/pkg/modules"

type UserRepository interface {
	CreateUser(username, email, hashedPassword, role string) (*modules.User, error)
	GetUserByUsername(username string) (*modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	PromoteUserToAdmin(id int) error
	UpdateUser(id int, username, email string) error
	DeleteUser(id int) error
	ListUsers() ([]*modules.User, error)
}
