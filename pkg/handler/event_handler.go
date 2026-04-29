package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
)

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	dateStart, dateEnd, userID, err := GetParser(r)
	if err != nil {
		log.Printf("order existance check error: %v", err)
		http.Error(w, "extracting order data error", http.StatusInternalServerError)
		return
	}
	// add existance check
	exists, err := h.Services.Users.UserExists(int64(userID))
	if err != nil {
		log.Printf("error occured during user existance check, %s", err)
		http.Error(w, "error occured during user existance check", http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Print("user with this id doesn't exists")
		http.Error(w, "user with this id doesn't exists", http.StatusBadRequest)
		return
	}
	//reply := make([]models.Reply, 0)
	reply, err := h.Services.Events.GetEvents(dateStart, dateEnd, userID)
	if err != nil {
		log.Printf("extracting order data error: %s", err)
		http.Error(w, "extracting order data error", http.StatusInternalServerError)
		return
	}
	if reply == nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("no events were found")
		return
	}
	fmt.Println(reply)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&reply); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(w, r)
	if userId == -1 {
		return
	}
	var input models.Event
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("failed to get event data, %s", err), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(input); err != nil {
		log.Printf("err %v", err)
		sendValidationErrors(w, err)
		return
	}

	if err := h.Services.Events.NewEvent(userId, input); err != nil {
		http.Error(w, fmt.Sprintf("failed to write input data to database %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(w, r)
	if userId == -1 {
		return
	}
	eventID := GetEventId(w, r)
	if exists, err := h.Services.Events.EventExists(userId, eventID); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, "no events were found", http.StatusNotFound)
		return
	}

	var input models.UpdateEvent
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("failed to get event data, %s", err), http.StatusBadRequest)
		return
	}

	err := h.Services.Events.UpdateEvent(userId, eventID, input)
	if errors.Is(err, repositories.EmptyUpdateError) {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update event %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(w, r)
	if userId == -1 {
		return
	}
	eventID := GetEventId(w, r)
	if exists, err := h.Services.Events.EventExists(userId, eventID); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, "no events were found", http.StatusNotFound)
		return
	}

	if err := h.Services.Events.DeleteEvent(userId, eventID); err != nil {
		http.Error(w, fmt.Sprintf("failed to write input data to database %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
