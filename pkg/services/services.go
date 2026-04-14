package services

import (
	"time"

	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
)

//go:generate minimock -i github.com/znmaster911/L2-calendar/pkg/services.* -o ./mocks -s _mock.go
type Services struct {
	Users
	Events
}

type Users interface {
	NewUser(username string) (int64, error)
	LogIn(username string) (int64, error)
}

type Events interface {
	NewEvent(userId int64, event models.Event) error
	UpdateEvent(userID, eventID int64, input models.UpdateEvent) error
	DeleteEvent(userID, eventID int64) error
	GetEvents(dateStart, dateEnd time.Time, userID int) ([]models.Reply, error)
}

func NewService(repo *repositories.Rpository) *Services {
	return &Services{
		Users:  NewUserService(repo.Users),
		Events: NewEventService(repo.Events),
	}
}
