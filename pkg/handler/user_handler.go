package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.Services.Users.NewUser()
	if err != nil {
		http.Error(w, "user creation error "+fmt.Sprintf("%s", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"user_id": id,
	})

}
