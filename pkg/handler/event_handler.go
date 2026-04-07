package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/znmaster911/L2-calendar/internal/models"
)

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	dateStart, dateEnd, userID, err := GetParser(r)
	if err != nil {
		log.Printf("order existance check error: %v", err)
		http.Error(w, "extracting order data error", http.StatusInternalServerError)
		return
	}
	// add existance check

	reply := make([]models.Reply, 0)
	reply, err = h.Services.Events.GetEvents(dateStart, dateEnd, userID)
	if err != nil {
		log.Printf("extracting order data error: %s", err)
		http.Error(w, "extracting order data error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&reply); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(w, r)

	var input models.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("failed to get event data, %s", err), http.StatusBadRequest)
		return
	}

	if err := h.Services.Events.NewEvent(userId, input); err != nil {
		http.Error(w, fmt.Sprintf("failed to write input data to database %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(w, r)
	eventID := GetEventId(w, r)

	var input models.UpdateEvent
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("failed to get event data, %s", err), http.StatusBadRequest)
		return
	}

	if err := h.Services.Events.UpdateEvent(userId, eventID, input); err != nil {
		http.Error(w, fmt.Sprintf("failed to update event %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(w, r)
	eventID := GetEventId(w, r)

	if err := h.Services.Events.DeleteEvent(userId, eventID); err != nil {
		http.Error(w, fmt.Sprintf("failed to write input data to database %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
