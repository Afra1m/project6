package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"project6/sessions"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // hashed password
}

func loadUsers() ([]User, error) {
	data, err := os.ReadFile("users.json")
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(data, &users)
	return users, err
}

func saveUsers(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("users.json", data, 0644)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds User
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &creds)

	users, _ := loadUsers()
	for _, u := range users {
		if u.Username == creds.Username {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	users = append(users, User{Username: creds.Username, Password: string(hash)})
	saveUsers(users)

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds User
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &creds)

	users, _ := loadUsers()
	for _, u := range users {
		if u.Username == creds.Username {
			if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(creds.Password)); err == nil {
				sessionID := sessions.CreateSession(u.Username)
				http.SetCookie(w, &http.Cookie{
					Name:     "session",
					Value:    sessionID,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Path:     "/",
				})
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sessions.DeleteSession(cookie.Value)

	// удалить cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
