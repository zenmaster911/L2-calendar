package services

import (
	"time"

	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
)

type EventService struct {
	repo repositories.Events
}

func NewEventService(repo repositories.Events) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) NewEvent(userId int64, event models.Event) error {
	return s.repo.NewEvent(userId, event)
}

func (s *EventService) UpdateEvent(userID, eventID int64, input models.UpdateEvent) error {
	return s.repo.UpdateEvent(userID, eventID, input)
}

func (s *EventService) DeleteEvent(userID, eventID int64) error {
	return s.repo.DeleteEvent(userID, eventID)
}

func (s EventService) GetEvents(dateStart, dateEnd time.Time, userID int) ([]models.Reply, error) {
	return s.repo.GetEvents(dateStart, dateEnd, userID)
}
