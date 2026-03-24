package repositories

import "github.com/jmoiron/sqlx"

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db sqlx.DB) *UsersPostgres {
	return &UsersPostgres{
		db: &db,
	}
}
