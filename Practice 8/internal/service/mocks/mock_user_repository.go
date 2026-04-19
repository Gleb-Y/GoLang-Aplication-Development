package mocks

import (
	"practice-8/pkg/modules"
	"practice-8/pkg/utils"

	"github.com/golang/mock/gomock"
)

type MockUserRepository struct {
	ctrl *gomock.Controller
}

func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	return &MockUserRepository{
		ctrl: ctrl,
	}
}

func (m *MockUserRepository) CreateUser(username, email, hashedPassword, role string) (*modules.User, error) {
	return &modules.User{
		ID:       1,
		Username: username,
		Email:    email,
		Role:     role,
	}, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*modules.User, error) {
	hashedPassword, _ := utils.HashPassword("password123")
	return &modules.User{
		ID:       1,
		Username: username,
		Email:    "user@example.com",
		Password: hashedPassword,
		Role:     "user",
	}, nil
}

func (m *MockUserRepository) GetUserByID(id int) (*modules.User, error) {
	return &modules.User{
		ID:       id,
		Username: "testuser",
		Email:    "user@example.com",
		Role:     "user",
	}, nil
}

func (m *MockUserRepository) PromoteUserToAdmin(id int) error {
	return nil
}

func (m *MockUserRepository) UpdateUser(id int, username, email string) error {
	return nil
}

func (m *MockUserRepository) DeleteUser(id int) error {
	return nil
}

func (m *MockUserRepository) ListUsers() ([]*modules.User, error) {
	return []*modules.User{}, nil
}
