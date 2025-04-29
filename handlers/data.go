package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"project6/cache"
	"project6/sessions"
)

type Data struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func CacheData(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_, ok := sessions.GetUser(cookie.Value)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	data, err := cache.GetData(func() any {
		return Data{
			Timestamp: time.Now().Format(time.RFC3339),
			Message:   "Это сгенерированные данные!",
		}
	})
	if err != nil {
		http.Error(w, "error reading cache", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
