package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/znmaster911/L2-calendar/internal/config"
)

type DB struct {
	DB         *sqlx.DB
	MaxRetries int
	RetryDelay int
}

func NewPostgresDb(cfg config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Print(dsn)
		return nil, fmt.Errorf("failed to connect to DB: %s", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %s", err)
	}

	return &DB{
		DB:         db,
		MaxRetries: cfg.MaxRetries,
		RetryDelay: cfg.RetryDelay,
	}, nil
}

func (d *DB) DbInit() {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Fatalf("failed to start transaction %s", err)
	}
	_, err = tx.Exec(` CREATE TABLE IF NOT EXISTS users(
	id BIGSERIAL PRIMARY KEY,
	username VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatalf("something went totally wrong")
		}
		log.Fatalf("failed to create or read users table %s", err)
	}

	_, err = tx.Exec(` CREATE TABLE IF NOT EXISTS events(
	id BIGSERIAL PRIMARY KEY,
	description TEXT NOT NULL,
	title VARCHAR(255) NOT NULL,
	starts DATE NOT NULL,
	deadline DATE,
	created TIMESTAMP NOT NULL DEFAULT now(),
	deleted TIMESTAMP DEFAULT NULL
	);`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatalf("something went totally wrong")
		}
		log.Fatalf("failed to create or read events table %s", err)
	}

	_, err = tx.Exec(` CREATE TABLE IF NOT EXISTS users_events(
	user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
	event_id  BIGINT REFERENCES events(id) ON DELETE CASCADE,
	PRIMARY KEY(user_id,event_id)
	);`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatalf("something went totally wrong")
		}
		log.Fatalf("failed to create or read users_events table %s", err)
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatalf("something went totally wrong")
		}
		log.Fatalf("failed to commit transaction %s", err)
	}
}
