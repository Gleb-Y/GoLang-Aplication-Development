package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"practice-8/internal/repository"
	"practice-8/pkg/modules"
	"practice-8/pkg/utils"
)

//go:generate mockgen -source=user_service.go -destination=mocks/mock_user_service.go -package=mocks UserService

type UserService interface {
	Register(req modules.RegisterRequest) (*modules.User, error)
	Login(req modules.LoginRequest) (string, *modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	PromoteUserToAdmin(id int) error
	UpdateUser(id int, username, email string) error
	DeleteUser(id int) error
	ListUsers() ([]*modules.User, error)
	GetRate(from, to string) (float64, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(req modules.RegisterRequest) (*modules.User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("missing required fields")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := "user"
	if req.Role == "admin" {
		return nil, errors.New("cannot create admin account")
	}

	user, err := s.repo.CreateUser(req.Username, req.Email, hashedPassword, role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) Login(req modules.LoginRequest) (string, *modules.User, error) {
	if req.Username == "" || req.Password == "" {
		return "", nil, errors.New("username and password required")
	}

	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *UserServiceImpl) GetUserByID(id int) (*modules.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetUserByID(id)
}

func (s *UserServiceImpl) PromoteUserToAdmin(id int) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	return s.repo.PromoteUserToAdmin(id)
}

func (s *UserServiceImpl) UpdateUser(id int, username, email string) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	if username == "" || email == "" {
		return errors.New("username and email are required")
	}
	return s.repo.UpdateUser(id, username, email)
}

func (s *UserServiceImpl) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	return s.repo.DeleteUser(id)
}

func (s *UserServiceImpl) ListUsers() ([]*modules.User, error) {
	return s.repo.ListUsers()
}

func (s *UserServiceImpl) GetRate(from, to string) (float64, error) {
	if from == "" || to == "" {
		return 0, errors.New("from and to currencies required")
	}

	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	_ = body
	return 1.0, nil
}
