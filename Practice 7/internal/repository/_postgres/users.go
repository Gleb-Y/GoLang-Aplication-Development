package _postgres

import (
	"database/sql"
	"errors"
	"practice-7/pkg/modules"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

// CreateUser creates a new user in the database
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

// GetUserByUsername retrieves a user by username
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

// GetUserByID retrieves a user by ID
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

// PromoteUserToAdmin promotes a user to admin role
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
