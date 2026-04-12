package users

import (
	"errors"
	"fmt"
	"prac5/internal/repository/_postgres"
	"prac5/pkg/modules"
	"strings"
)

type Repository struct {
	db *_postgres.Dialect
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{db: db}
}

var allowedOrderBy = map[string]bool{
	"id": true, "name": true, "email": true, "gender": true, "birth_date": true,
}

func (r *Repository) GetPaginatedUsers(page, pageSize int, filter modules.UserFilter) (modules.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	args := []interface{}{}
	conditions := []string{}
	argIdx := 1

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, *filter.ID)
		argIdx++
	}
	if filter.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Name+"%")
		argIdx++
	}
	if filter.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Email+"%")
		argIdx++
	}
	if filter.Gender != nil {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, *filter.Gender)
		argIdx++
	}
	if filter.BirthDate != nil {
		conditions = append(conditions, fmt.Sprintf("birth_date = $%d", argIdx))
		args = append(args, *filter.BirthDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := "id"
	if filter.OrderBy != "" && allowedOrderBy[filter.OrderBy] {
		orderBy = filter.OrderBy
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM users" + whereClause
	if err := r.db.DB.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return modules.PaginatedResponse{}, err
	}

	offset := (page - 1) * pageSize
	dataArgs := append(args, pageSize, offset)
	dataQuery := fmt.Sprintf(
		"SELECT id, name, email, gender, birth_date FROM users%s ORDER BY %s LIMIT $%d OFFSET $%d",
		whereClause, orderBy, argIdx, argIdx+1,
	)

	rows, err := r.db.DB.Query(dataQuery, dataArgs...)
	if err != nil {
		return modules.PaginatedResponse{}, err
	}
	defer rows.Close()

	var userList []modules.User
	for rows.Next() {
		var u modules.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return modules.PaginatedResponse{}, err
		}
		userList = append(userList, u)
	}

	return modules.PaginatedResponse{
		Data:       userList,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT id, name, email, gender, birth_date FROM users WHERE id = $1", id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("user with id %d not found", id))
	}
	return &user, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var id int
	err := r.db.DB.QueryRow(
		"INSERT INTO users (name, email, gender, birth_date) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.Gender, user.BirthDate,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
	result, err := r.db.DB.Exec(
		"UPDATE users SET name=$1, email=$2, gender=$3, birth_date=$4 WHERE id=$5",
		user.Name, user.Email, user.Gender, user.BirthDate, id,
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

func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]modules.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON u.id = uf1.friend_id AND uf1.user_id = $1
		JOIN user_friends uf2 ON u.id = uf2.friend_id AND uf2.user_id = $2
	`
	rows, err := r.db.DB.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []modules.User
	for rows.Next() {
		var u modules.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}
	return userList, nil
}
