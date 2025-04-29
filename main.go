package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"project6/cache"
	"project6/handlers"
)

func main() {
	sessions.init() // Инициализация сессий

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/profile", handlers.Profile)
	http.HandleFunc("/data", handlers.CacheData)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data, err := cache.GetData(func() any {
			return map[string]any{
				"time": time.Now().Format(time.RFC3339),
				"info": "Это сгенерированные данные",
			}
		})
		if err != nil {
			http.Error(w, "Не удалось получить данные", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(data)
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
