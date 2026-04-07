package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const tformat = "2006-01-02"

func GetParser(r *http.Request) (dateStart, dateEnd time.Time, userID int, err error) {
	rawUrl := r.URL.RawPath
	p, err := url.Parse(rawUrl)
	if err != nil {
		return time.Now(), time.Now(), 0, fmt.Errorf("failed to get data from query string %s", err)
	}
	q := p.Query()
	rawdate := q.Get("date")
	dataslc := strings.Split(rawdate, "-")
	order := len(dataslc)
	switch order {
	case 1:
		dataslc = append(dataslc, "01", "1")
	case 2:
		dataslc = append(dataslc, "1")
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
	userID, err = strconv.Atoi(q.Get("user_id"))
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
		}
	}
	for _, v := range endDateIntSlc {
		indDate := strconv.Itoa(v)
		if len(indDate) == 1 {
			indDate = "0" + indDate
		}
		endDate += indDate + "-"
	}
	return endDate, nil
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

func getUserId(w http.ResponseWriter, r *http.Request) (userID int64, err error) {
	userIDraw := r.Context().Value("userID")
	if userIDraw == nil {
		http.Error(w, "no User ID in context", http.StatusUnauthorized)
		return 0, fmt.Errorf("no user with current ID found")
	}

	userID, ok := userIDraw.(int64)
	if !ok {
		http.Error(w, fmt.Sprintf("wrong type of user ID in context %d", userID), http.StatusInternalServerError)
		return 0, fmt.Errorf("wrong type of user ID in context")
	}

	return userID, nil
}
