package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/handler/middleware"
)

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	dateStart, dateEnd, userID, err := middleware.GetParser(r)
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
