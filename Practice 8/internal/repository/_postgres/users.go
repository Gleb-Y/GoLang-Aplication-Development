package _postgres

import (
	"database/sql"
	"errors"
	"practice-8/pkg/modules"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) CreateUser(username, email, hashedPassword, role string) (*modules.User, error) {
	user := &modules.User{}
	err := r.db.QueryRow(
		"INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, username, email, role",
		username, email, hashedPassword, role,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*modules.User, error) {
	user := &modules.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, password, role FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*modules.User, error) {
	user := &modules.User{}
	err := r.db.QueryRow(
		"SELECT id, username, email, password, role FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) PromoteUserToAdmin(id int) error {
	result, err := r.db.Exec(
		"UPDATE users SET role = $1 WHERE id = $2",
		"admin", id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) UpdateUser(id int, username, email string) error {
	result, err := r.db.Exec(
		"UPDATE users SET username = $1, email = $2 WHERE id = $3",
		username, email, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	result, err := r.db.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) ListUsers() ([]*modules.User, error) {
	var users []*modules.User
	err := r.db.Select(&users, "SELECT id, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}
