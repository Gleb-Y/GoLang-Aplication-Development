package users

import (
	"errors"
	"fmt"
	"prac3/internal/repository/_postgres"
	"prac3/pkg/modules"
	"time"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id, name, email, age, phone FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT id, name, email, age, phone FROM users WHERE id = $1", id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("user with id %d not found", id))
	}
	return &user, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var id int
	err := r.db.DB.QueryRow(
		"INSERT INTO users (name, email, age, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.Age, user.Phone,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
	result, err := r.db.DB.Exec(
		"UPDATE users SET name=$1, email=$2, age=$3, phone=$4 WHERE id=$5",
		user.Name, user.Email, user.Age, user.Phone, id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf("user with id %d does not exist", id))
	}
	return nil
}

func (r *Repository) DeleteUser(id int) error {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf("user with id %d does not exist", id))
	}
	return nil
}
