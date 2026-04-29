package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
)

var validate = validator.New()

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf("failed to get user data due to %s", err), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		log.Printf("err %v", err)
		sendValidationErrors(w, err)
		return
	}

	id, err := h.Services.Users.NewUser(user.Username)
	if errors.Is(err, repositories.UserExistsError) {
		http.Error(w, repositories.UserExistsError.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "user creation error "+fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"user_id": id,
	})

}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	var input models.LogIn
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode login data due to %s", err), http.StatusBadRequest)
		return
	}
	id, err := h.Services.Users.LogIn(input.Username)
	if err != nil {
		http.Error(w, "login error "+fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("User_id", fmt.Sprintf("%d", id))
	w.WriteHeader(http.StatusOK)
}
