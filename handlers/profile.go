package handlers

import (
	"fmt"
	"net/http"

	"project6/sessions"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username, ok := sessions.GetUser(cookie.Value)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Простой вывод профиля
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Welcome, %s!"}`, username)
}
