package usecase

import (
	"prac5/internal/repository"
	"prac5/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetPaginatedUsers(page, pageSize int, filter modules.UserFilter) (modules.PaginatedResponse, error) {
	return u.repo.GetPaginatedUsers(page, pageSize, filter)
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUsecase) CreateUser(user modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

func (u *UserUsecase) UpdateUser(id int, user modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *UserUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}

func (u *UserUsecase) GetCommonFriends(userID1, userID2 int) ([]modules.User, error) {
	return u.repo.GetCommonFriends(userID1, userID2)
}
