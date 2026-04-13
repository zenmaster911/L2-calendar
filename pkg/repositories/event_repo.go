package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/znmaster911/L2-calendar/internal/models"
)

type EventsPostgres struct {
	db *sqlx.DB
}

func NewEventsPostgres(db *sqlx.DB) *EventsPostgres {
	return &EventsPostgres{
		db: db,
	}
}

func (r *EventsPostgres) NewEvent(userId int64, event models.Event) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %s", err)
	}
	var eventID int64
	eventQuery := "INSERT INTO events ( description, title, starts, deadline) VALUES ($1,$2,$3,$4) RETURNING id"
	if err := tx.QueryRow(eventQuery, event.Description, event.Title, event.Starts, event.Deadline).Scan(&eventID); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction to events failed: %s", err)
	}

	usersEventQuery := "INSERT INTO users_events (event_id, user_id) VALUES($1,$2)"
	if _, err := tx.Exec(usersEventQuery, eventID, userId); err != nil {
		return fmt.Errorf("transaction to users events failed: %s", err)
	}
	tx.Commit()
	return nil
}

func (r *EventsPostgres) UpdateEvent(userID, eventID int64, input models.UpdateEvent) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	if input.Starts != nil {
		setValues = append(setValues, fmt.Sprintf("starts=$%d", argID))
		args = append(args, *input.Starts)
		argID++
	}

	if input.Deadline != nil {
		setValues = append(setValues, fmt.Sprintf("deadline=$%d", argID))
		args = append(args, *input.Deadline)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE events e SET %s FROM users_events ue WHERE e.id=ue.event_id AND ue.event_id=$%d AND ue.user_id=$%d",
		setQuery, argID, argID+1)
	args = append(args, eventID, userID)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *EventsPostgres) DeleteEvent(userID, eventID int64) error {
	query := "DELETE FROM events e USING users_events ue WHERE e.id=ue.event_id AND ue.event_id=$1 AND ue.user_id=$2"
	_, err := r.db.Exec(query, eventID, userID)
	return err
}

func (r *EventsPostgres) GetEvents(dateStart, dateEnd time.Time, userID int) ([]models.Reply, error) {
	var relpies []models.Reply
	query := `SELECT e.title, e.description,starts,deadline FROM events e 
	INNER JOIN users_events ue ON e.id=ue.event_id
	WHERE e.deadline>=$1 AND e.starts<=$2 AND ue.user_id=$3`
	if err := r.db.Select(&relpies, query, dateStart, dateEnd, userID); err != nil {
		return nil, fmt.Errorf("failed to get events: %s", err)
	}
	return relpies, nil
}
