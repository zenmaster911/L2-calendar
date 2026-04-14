package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/znmaster911/L2-calendar/internal/models"
)

//go:generate minimock -i github.com/znmaster911/L2-calendar/pkg/repositories.* -o ./mocks -s _mock.go
type Rpository struct {
	Users
	Events
}

type Users interface {
	NewUser(string) (int64, error)
	LogIn(username string) (int64, error)
	UserExists(username string) (bool, error)
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
