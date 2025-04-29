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

func CachedData(w http.ResponseWriter, r *http.Request) {
	// проверка сессии
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

	// если кэш валиден, отдаем его
	if cache.IsCacheValid() {
		data, err := cache.ReadCache()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}
	}

	// если кэша нет или он устарел
	newData := Data{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "Это сгенерированные данные!",
	}
	response, _ := json.Marshal(newData)
	cache.WriteCache(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
