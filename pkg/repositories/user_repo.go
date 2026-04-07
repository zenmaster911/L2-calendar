package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var UserExistsError = fmt.Errorf("user already exists")

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) NewUser(username string) (int64, error) {
	var id int64
	exists, err := r.UserExists(username)
	if err != nil {
		return -1, fmt.Errorf("failed to create new user, %w", err)
	}
	if exists == true {
		return -1, UserExistsError
	}
	query := " INSERT INTO users (username) VALUES ($1) RETURNING id"
	if err := r.db.QueryRow(query, username).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to add user to DB: %s", err)
	}
	return id, nil
}

func (r *UsersPostgres) LogIn(username string) (int64, error) {
	var id int64
	query := "GET id FROM users WHERE username=$1"
	err := r.db.Get(&id, query, username)
	if err != nil {
		return -1, fmt.Errorf("failed to get user id due to %w", err)
	}
	return id, nil
}

func (r *UsersPostgres) UserExists(username string) (bool, error) {
	exists := false
	query := "SELECT EXISTS(SLECT 1 FROM users WHERE username=$1)"
	if err := r.db.QueryRow(query, username).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check user existance due to %w", err)
	}
	return exists, nil
}
