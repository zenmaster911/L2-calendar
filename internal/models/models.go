package models

import "time"

type User struct {
	Id       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}

type LogIn struct {
	Username string `json:"username" db:"username"`
	//potential password
}

type Event struct {
	ID          int64     `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	Title       string    `json:"title" db:"title"`
	Starts      time.Time `json:"starts" validate:"required" db:"starts"`
	Deadline    time.Time `json:"deadline" db:"deadline"`
	Created     time.Time `json:"created" validate:"required" db:"created"`
	Deleted     time.Time `json:"deleted" validate:"required" db:"deleted"`
}

type Reply struct {
	Description string    `json:"description" db:"description"`
	Title       string    `json:"title" db:"title"`
	Starts      time.Time `json:"starts" db:"starts"`
	Deadline    time.Time `json:"deadline" db:"deadline"`
}

type UserEvents struct {
	UserId  int64 `json:"user_id" validate:"required" db:"user_id"`
	EventID int64 `json:"event_id" validate:"required" db:"event_id"`
}

type UpdateEvent struct {
	Description *string    `json:"description" db:"description"`
	Title       *string    `json:"title" db:"title"`
	Starts      *time.Time `json:"starts" db:"starts"`
	Deadline    *time.Time `json:"deadline" db:"deadline"`
}
