package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
	"database/sql"
	_ "modernc.org/sqlite" // sqlite driver
)


type contextKey string

const userContextKey contextKey = "user"

type sessionStore struct {
	mu sync.RWMutex
	sessions map[string]string // sessionID -> email
}

var store = &sessionStore{
	sessions: make(map[string]string),
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	db, err := sql.Open("sqlite", "file:icas.db")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, password FROM users WHERE email = ?", email)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
			fmt.Print(err)
		}
		if verifyPassword(password, u.Password) {
			found = true
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Print(err)
	}
	if found {
		sessionID, err := generateSessionID()
		if err != nil {
			fmt.Print(err)
		}

		store.mu.Lock()
		store.sessions[sessionID] = email
		store.mu.Unlock()
		http.SetCookie(w, &http.Cookie{
			Name: "session_id",
			Value: sessionID,
			Path: "/",
			HttpOnly: true,
			Secure: true,
			SameSite: http.SameSiteLaxMode,
			Expires: time.Now().Add(24 * time.Hour),
		})
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "found")
	} else {
		fmt.Fprint(w, "not found")
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized: no session cookie", http.StatusUnauthorized)
			return
		}

		store.mu.RLock()
		email, exists := store.sessions[cookie.Value]
		store.mu.RUnlock()

		if !exists {
			http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), userContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
