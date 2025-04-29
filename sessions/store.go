package sessions

import (
	"math/rand"
	"sync"
	"time"
)

var store map[string]string
var mu sync.RWMutex

func Init() {
	store = make(map[string]string)
	rand.Seed(time.Now().UnixNano())
}

func CreateSession(username string) string {
	sessionID := RandString(32)
	mu.Lock()
	store[sessionID] = username
	mu.Unlock()
	return sessionID
}

func GetUser(sessionID string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	user, ok := store[sessionID]
	return user, ok
}

func DeleteSession(sessionID string) {
	mu.Lock()
	delete(store, sessionID)
	mu.Unlock()
}

func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
