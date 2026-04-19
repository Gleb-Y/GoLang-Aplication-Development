package usecase

import (
	"os"
	"practice-8/internal/repository"
	"practice-8/pkg/modules"
	"practice-8/pkg/utils"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Register(req modules.RegisterRequest) (*modules.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

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

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) PromoteUserToAdmin(id int) error {
	return u.repo.PromoteUserToAdmin(id)
}
