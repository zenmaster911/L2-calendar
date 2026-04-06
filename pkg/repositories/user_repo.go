package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{
		db: db,
	}
}

func (r *UsersPostgres) NewUser() (int64, error) {
	var id int64
	query := " INSERT INTO users DDEFAULT VALUES"
	if err := r.db.QueryRow(query).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to add user to DB: %s", err)
	}
	return id, nil
}
