package usecase

import (
	"os"
	"practice-7/internal/repository"
	"practice-7/pkg/modules"
	"practice-7/pkg/utils"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Register creates a new user account
func (u *UserUsecase) Register(req modules.RegisterRequest) (*modules.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Determine role: allow admin registration only if env var is set
	role := "user"
	if req.Role != "" && req.Role == "admin" && os.Getenv("ALLOW_ADMIN_REGISTRATION") == "true" {
		role = "admin"
	}

	user, err := u.repo.CreateUser(req.Username, req.Email, hashedPassword, role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Login authenticates a user and returns a token
func (u *UserUsecase) Login(req modules.LoginRequest) (string, *modules.User, error) {
	user, err := u.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", nil, err
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return "", nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// GetUserByID retrieves a user by ID
func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

// PromoteUserToAdmin promotes a user to admin role
func (u *UserUsecase) PromoteUserToAdmin(id int) error {
	return u.repo.PromoteUserToAdmin(id)
}
