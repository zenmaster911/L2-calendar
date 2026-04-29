package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/znmaster911/L2-calendar/internal/logger"
	"github.com/znmaster911/L2-calendar/internal/models"
	"github.com/znmaster911/L2-calendar/pkg/repositories"
	"github.com/znmaster911/L2-calendar/pkg/services"
	"github.com/znmaster911/L2-calendar/pkg/services/mocks"
)

func TestHandler_GetEvent(t *testing.T) {
	type mockBehavior func(e *mocks.EventsMock, u *mocks.UsersMock)
	goodReply := models.Reply{
		Description: "Обсуждение целей на Q3",
		Title:       "просто 3",
		Starts:      time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Deadline:    time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
	}
	year := models.Reply{
		Starts:   time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
		Deadline: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
	}
	month := models.Reply{
		Starts:   time.Date(2024, 07, 01, 0, 0, 0, 0, time.UTC),
		Deadline: time.Date(2024, 8, 01, 0, 0, 0, 0, time.UTC),
	}
	day := models.Reply{
		Starts:   time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Deadline: time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
	}

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		queryParams          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK-year",
			queryParams: "?user_id=1&date=2024",
			mockBehavior: func(e *mocks.EventsMock, u *mocks.UsersMock) {

				u.UserExistsMock.Expect(int64(1)).Return(true, nil)
				e.GetEventsMock.Expect(year.Starts, year.Deadline, 1).Return([]models.Reply{(goodReply)}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"description":"Обсуждение целей на Q3","title":"просто 3","starts":"2024-07-15T00:00:00Z","deadline":"2024-07-15T00:00:00Z"}]` + "\n",
		},
		{
			name:        "OK-month",
			queryParams: "?user_id=1&date=2024-07",
			mockBehavior: func(e *mocks.EventsMock, u *mocks.UsersMock) {

				u.UserExistsMock.Expect(int64(1)).Return(true, nil)
				e.GetEventsMock.Expect(month.Starts, month.Deadline, 1).Return([]models.Reply{(goodReply)}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"description":"Обсуждение целей на Q3","title":"просто 3","starts":"2024-07-15T00:00:00Z","deadline":"2024-07-15T00:00:00Z"}]` + "\n",
		},
		{
			name:        "OK-day",
			queryParams: "?user_id=1&date=2024-07-15",
			mockBehavior: func(e *mocks.EventsMock, u *mocks.UsersMock) {

				u.UserExistsMock.Expect(int64(1)).Return(true, nil)
				e.GetEventsMock.Expect(day.Starts, day.Deadline, 1).Return([]models.Reply{(goodReply)}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"description":"Обсуждение целей на Q3","title":"просто 3","starts":"2024-07-15T00:00:00Z","deadline":"2024-07-15T00:00:00Z"}]` + "\n",
		}, {
			name:        "invalid user",
			queryParams: "?user_id=123&date=2024",
			mockBehavior: func(e *mocks.EventsMock, u *mocks.UsersMock) {
				u.UserExistsMock.Expect(123).Return(false, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "user with this id doesn't exists\n",
		}, {
			name:        "usercheck error",
			queryParams: "?user_id=123&date=2024",
			mockBehavior: func(e *mocks.EventsMock, u *mocks.UsersMock) {
				u.UserExistsMock.Expect(123).Return(false, assert.AnError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "error occured during user existance check\n",
		},
		{
			name: "no events",
			mockBehavior: func(e *mocks.EventsMock, u *mocks.UsersMock) {
				u.UserExistsMock.Expect(1).Return(true, nil)
				e.GetEventsMock.Expect(year.Starts, year.Deadline, 1).Return(nil, nil)
			},
			queryParams:          "?user_id=1&date=2024",
			expectedStatusCode:   204,
			expectedResponseBody: "",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			umockservice := mocks.NewUsersMock(c)
			emockservice := mocks.NewEventsMock(c)
			tt.mockBehavior(emockservice, umockservice)
			logg, _ := logger.NewSLog("/testlogs.txt", slog.LevelDebug)
			services := &services.Services{Users: umockservice, Events: emockservice}
			h := NewHandler(services, logg)
			router := chi.NewRouter()
			router.Get("/", h.GetEvent)
			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/"+tt.queryParams, nil)

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_CreateEvent(t *testing.T) {
	type mockBehavior func(e *mocks.EventsMock, input models.Event)
	goodInput := models.Event{
		Title:       "test",
		Description: "test description",
		Starts:      time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Deadline:    time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Created:     time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
	}
	NoStartInput := models.Event{
		Title:       "test",
		Description: "test description",
		Deadline:    time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Created:     time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
	}
	NoTitleInput := models.Event{
		Description: "test description",
		Starts:      time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Deadline:    time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
		Created:     time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
	}

	testTable := []struct {
		name                 string
		inputArgs            models.Event
		userIdCtx            any
		mockBehavior         mockBehavior
		expectedStatuscode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputArgs: goodInput,
			userIdCtx: int64(1),
			mockBehavior: func(e *mocks.EventsMock, input models.Event) {
				e.NewEventMock.Expect(int64(1), goodInput).Return(nil)
			},
			expectedStatuscode:   201,
			expectedResponseBody: "",
		},
		{
			name:                 "no start time",
			inputArgs:            NoStartInput,
			mockBehavior:         func(e *mocks.EventsMock, input models.Event) {},
			userIdCtx:            int64(1),
			expectedStatuscode:   400,
			expectedResponseBody: `{"error":"Validation failed","fields":["Starts is required"]}` + "\n",
		},
		{
			name:                 "no title",
			inputArgs:            NoTitleInput,
			mockBehavior:         func(e *mocks.EventsMock, input models.Event) {},
			userIdCtx:            int64(1),
			expectedStatuscode:   400,
			expectedResponseBody: `{"error":"Validation failed","fields":["Title is required"]}` + "\n",
		},
		{
			name:                 "unregistred user",
			inputArgs:            goodInput,
			mockBehavior:         func(e *mocks.EventsMock, input models.Event) {},
			userIdCtx:            nil,
			expectedStatuscode:   401,
			expectedResponseBody: `no User ID in context` + "\n",
		},
		{
			name:               "db input err",
			inputArgs:          goodInput,
			userIdCtx:          int64(1),
			expectedStatuscode: 500,
			mockBehavior: func(e *mocks.EventsMock, input models.Event) {
				e.NewEventMock.Expect(1, goodInput).Return(assert.AnError)
			},
			expectedResponseBody: `failed to write input data to database assert.AnError general error for testing` + "\n",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockService := mocks.NewEventsMock(c)
			tt.mockBehavior(mockService, tt.inputArgs)
			services := &services.Services{Events: mockService}
			logg, _ := logger.NewSLog("/testlogs.txt", slog.LevelDebug)
			h := NewHandler(services, logg)
			router := chi.NewRouter()
			router.Post("/events", h.CreateEvent)
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.inputArgs)
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
			req.Header.Set("ContentType", "application/json")
			if tt.userIdCtx != nil {
				ctx := context.WithValue(req.Context(), "userId", tt.userIdCtx)
				req = req.WithContext(ctx)
			}

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatuscode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())

		})
	}
}

func TestHandler_UpdateEvent(t *testing.T) {
	type mockBehavior func(e *mocks.EventsMock, input models.UpdateEvent)
	title, desctiption := "test", "test description"
	starts, edns := time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC)

	goodUpdate := models.UpdateEvent{
		Title:       &title,
		Description: &desctiption,
		Starts:      &starts,
		Deadline:    &edns,
	}

	emptyUpdate := models.UpdateEvent{}

	testTable := []struct {
		name                 string
		inputArgs            models.UpdateEvent
		userIdCtx            any
		eventId              int64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputArgs: goodUpdate,
			userIdCtx: int64(1),
			eventId:   1,
			mockBehavior: func(e *mocks.EventsMock, input models.UpdateEvent) {
				e.EventExistsMock.Expect(1, 1).Return(true, nil)
				e.UpdateEventMock.Expect(1, 1, goodUpdate).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:                 "inavalid user",
			userIdCtx:            nil,
			mockBehavior:         func(e *mocks.EventsMock, input models.UpdateEvent) {},
			expectedStatusCode:   401,
			expectedResponseBody: `no User ID in context` + "\n",
		},
		{
			name:      "empty update",
			userIdCtx: int64(1),
			eventId:   1,
			mockBehavior: func(e *mocks.EventsMock, input models.UpdateEvent) {
				e.EventExistsMock.Expect(1, 1).Return(true, nil)
				e.UpdateEventMock.Expect(1, 1, emptyUpdate).Return(repositories.EmptyUpdateError)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `no parameters entered to be updated` + "\n",
		},
		{
			name:      "no event",
			userIdCtx: int64(1),
			eventId:   404,
			mockBehavior: func(e *mocks.EventsMock, input models.UpdateEvent) {
				e.EventExistsMock.Expect(1, 404).Return(false, nil)
			},
			expectedStatusCode:   404,
			expectedResponseBody: "no events were found\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockService := mocks.NewEventsMock(c)
			tt.mockBehavior(mockService, tt.inputArgs)
			services := &services.Services{Events: mockService}
			logg, _ := logger.NewSLog("/testlogs.txt", slog.LevelDebug)
			h := NewHandler(services, logg)
			router := chi.NewRouter()
			router.Patch("/events/{event_id}", h.UpdateEvent)
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.inputArgs)
			req := httptest.NewRequest(http.MethodPatch, "/events/"+strconv.FormatInt(tt.eventId, 10), bytes.NewReader(body))
			if tt.userIdCtx != nil {
				ctx := context.WithValue(req.Context(), "userId", tt.userIdCtx)
				req = req.WithContext(ctx)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}

}

func TestHandler_DeleteEvent(t *testing.T) {
	type mockBehavior func(e *mocks.EventsMock)
	testTable := []struct {
		name                 string
		userIdCtx            any
		eventId              int64
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			userIdCtx: int64(1),
			eventId:   1,
			mockBehavior: func(e *mocks.EventsMock) {
				e.EventExistsMock.Expect(1, 1).Return(true, nil)
				e.DeleteEventMock.Expect(1, 1).Return(nil)
			},
			expectedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name:                 "inavalid user",
			userIdCtx:            nil,
			mockBehavior:         func(e *mocks.EventsMock) {},
			expectedStatusCode:   401,
			expectedResponseBody: `no User ID in context` + "\n",
		},
		{
			name:      "no event",
			userIdCtx: int64(1),
			eventId:   404,
			mockBehavior: func(e *mocks.EventsMock) {
				e.EventExistsMock.Expect(1, 404).Return(false, nil)
			},
			expectedStatusCode:   404,
			expectedResponseBody: "no events were found\n",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := minimock.NewController(t)
			mockservice := mocks.NewEventsMock(c)
			tt.mockBehavior(mockservice)
			service := services.Services{Events: mockservice}
			logg, _ := logger.NewSLog("/testlog.txt", slog.LevelDebug)
			h := NewHandler(&service, logg)
			router := chi.NewRouter()
			router.Delete("/events/{event_id}", h.DeleteEvent)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/events/"+strconv.FormatInt(tt.eventId, 10), nil)
			if tt.userIdCtx != nil {
				ctx := context.WithValue(req.Context(), "userId", tt.userIdCtx)
				req = req.WithContext(ctx)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())

		})
	}
}
