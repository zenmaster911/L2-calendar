package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/znmaster911/L2-calendar/internal/logger"
	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
	"github.com/znmaster911/L2-calendar/pkg/services"
	"github.com/znmaster911/L2-calendar/pkg/services/mocks"
)

func TestHandler_NewUser(t *testing.T) {
	type mockBehavior func(u *mocks.UsersMock, input models.User)

	testTable := []struct {
		name                 string
		inputArgs            models.User
		inputErr             error
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{{
		name:      "OK",
		inputArgs: models.User{Username: "username1"},
		inputErr:  nil,
		mockBehavior: func(u *mocks.UsersMock, input models.User) {
			u.NewUserMock.Expect("username1").Return(1, nil)
		},
		expectedStatusCode:   201,
		expectedResponseBody: `{"user_id":1}` + "\n",
	}, {
		name:      "usernameExists",
		inputArgs: models.User{Username: "username1"},
		inputErr:  nil,
		mockBehavior: func(u *mocks.UsersMock, input models.User) {
			u.NewUserMock.Expect("username1").Return(-1, repositories.UserExistsError)
		},
		expectedStatusCode:   409,
		expectedResponseBody: `user already exists` + "\n",
	}, {
		name:                 "emptyUsername",
		inputArgs:            models.User{Username: ""},
		inputErr:             nil,
		mockBehavior:         func(u *mocks.UsersMock, input models.User) {},
		expectedStatusCode:   http.StatusBadRequest,
		expectedResponseBody: `{"error":"Validation failed","fields":["Username is required"]}` + "\n",
	},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockService := mocks.NewUsersMock(c)
			tt.mockBehavior(mockService, tt.inputArgs)
			services := &services.Services{Users: mockService}
			logg, _ := logger.NewSLog("/testlogs.txt", slog.LevelDebug)
			h := NewHandler(services, logg)
			router := chi.NewRouter()
			router.Post("/users", h.NewUser)
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.inputArgs)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("ContentType", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())

		})
	}
}

func TestHandler_LogIn(t *testing.T) {
	type mockBehavior func(u *mocks.UsersMock, input models.LogIn)

	testTable := []struct {
		name                 string
		inputArgs            models.LogIn
		inputError           error
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK",
			inputArgs:  models.LogIn{Username: "username1"},
			inputError: nil,
			mockBehavior: func(u *mocks.UsersMock, input models.LogIn) {
				u.LogInMock.Expect("username1").Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		}, {
			name:       "wrong username",
			inputArgs:  models.LogIn{Username: "username1"},
			inputError: nil,
			mockBehavior: func(u *mocks.UsersMock, input models.LogIn) {
				u.LogInMock.Expect("username1").Return(-1, fmt.Errorf("login error failed to get user id due to sql: no rows in result set"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "login error login error failed to get user id due to sql: no rows in result set\n",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockService := mocks.NewUsersMock(c)
			tt.mockBehavior(mockService, tt.inputArgs)
			services := &services.Services{Users: mockService}
			logg, _ := logger.NewSLog("/testlogs.txt", slog.LevelDebug)
			h := NewHandler(services, logg)
			router := chi.NewRouter()
			router.Post("/auth", h.LogIn)
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.inputArgs)
			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
			req.Header.Set("ContentType", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())

		})
	}

}
