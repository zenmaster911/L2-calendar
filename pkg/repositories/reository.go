package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/znmaster911/L2-calendar/internal/models"
)

type Rpository struct {
	Users
	Events
}

type Users interface {
	NewUser() (int, error)
}

type Events interface {
	NewEvent(userId int64, event models.Event) error
	UpdateEvent(userID, eventID int64, input models.UpdateEvent) error
	DeleteEvent(userID, eventID int64) error
	GetEvents(dateStart, dateEnd time.Time, userID int) ([]models.Reply, error)
}

func NewRepo(db *sqlx.DB) *Rpository {
	return &Rpository{
		Users:  NewUsersPostgres(db),
		Events: NewEventsPostgres(db),
	}
}
