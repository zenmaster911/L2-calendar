package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/znmaster911/L2-calendar/internal/config"
)

type DB struct {
	db         *sqlx.DB
	MaxRetries int
	RetryDelay int
}

func NewPostgresDb(cfg config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %s", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %s", err)
	}

	return &DB{
		db:         db,
		MaxRetries: cfg.MaxRetries,
		RetryDelay: cfg.RetryDelay,
	}, nil
}
