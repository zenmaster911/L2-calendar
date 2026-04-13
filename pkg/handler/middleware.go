package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const tformat = "2006-01-02"

func GetParser(r *http.Request) (dateStart, dateEnd time.Time, userID int, err error) {
	rawQuery := r.URL.RawQuery
	p, err := url.ParseQuery(rawQuery)
	if err != nil {
		return time.Now(), time.Now(), 0, fmt.Errorf("failed to get data from query string %s", err)
	}
	rawdate := p.Get("date")
	dataslc := strings.Split(rawdate, "-")
	order := len(dataslc)
	switch order {
	case 1:
		dataslc = append(dataslc, "01", "01")
	case 2:
		dataslc = append(dataslc, "01")
	case 3:
	default:
		return time.Now(), time.Now(), -1, fmt.Errorf("invalid data format %s", err)
	}

	enddata, err := EndDateReceiver(dataslc, order)
	if err != nil {
		return time.Now(), time.Now(), -1, fmt.Errorf("failed to get date %s", err)
	}
	dateEnd, err = time.Parse(tformat, enddata)
	if err != nil {
		return time.Now(), time.Now(), -1, fmt.Errorf("failed to parse starting date %s", err)
	}

	dateStart, err = time.Parse(tformat, strings.Join(dataslc, "-"))

	if err != nil {
		return time.Now(), time.Now(), -1, fmt.Errorf("failed to parse closing date %s", err)
	}
	userID, err = strconv.Atoi(p.Get("user_id"))
	if err != nil {
		return time.Now(), time.Now(), -1, fmt.Errorf("failed to get user_id from query string  %s", err)
	}
	return dateStart, dateEnd, userID, nil
}

func EndDateReceiver(date []string, order int) (string, error) {
	endDateIntSlc := make([]int, 0, len(date))
	var endDate string
	for _, v := range date {
		indDate, err := strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("failed to convert data due to %s", err)
		}
		endDateIntSlc = append(endDateIntSlc, indDate)
	}
	switch order {
	case 1:
		endDateIntSlc[0]++
	case 2:
		if endDateIntSlc[1] == 12 {
			endDateIntSlc[0]++
			endDateIntSlc[1] = 1
		} else {
			endDateIntSlc[1]++
		}
	}
	for _, v := range endDateIntSlc {
		indDate := strconv.Itoa(v)
		if len(indDate) == 1 {
			indDate = "0" + indDate
		}
		endDate += indDate + "-"
	}
	return endDate[:10], nil
}

func getUserId(w http.ResponseWriter, r *http.Request) (userID int64) {
	userIDraw := r.Context().Value("userId")
	if userIDraw == nil {
		http.Error(w, "no User ID in context", http.StatusUnauthorized)
		return -1
	}

	userID, ok := userIDraw.(int64)
	if !ok {
		http.Error(w, fmt.Sprintf("wrong type of user ID in context %d", userID), http.StatusInternalServerError)
		return -1
	}

	return userID
}

func GetEventId(w http.ResponseWriter, r *http.Request) (eventID int64) {
	param := chi.URLParam(r, "event_id")
	eventID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get eventID %s", err), http.StatusBadRequest)
		return -1
	}
	return eventID
}

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("User_id")
		if header == "" {
			http.Error(w, "empty authorization header", http.StatusUnauthorized)
			return
		}
		id, err := strconv.ParseInt(header, 10, 64)
		if err != nil {
			http.Error(w, "failed to get user identity", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) rLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		start := time.Now()
		next.ServeHTTP(w, r)
		h.Logger.Logger.Info("http request",
			slog.String("req_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
