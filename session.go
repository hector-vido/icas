package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	_ "modernc.org/sqlite" // sqlite driver
)

type contextKey string

const userContextKey contextKey = "userName"

type sessionStore struct {
	mu       sync.RWMutex
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

	if r.Method == "GET" {
		_, err := getActiveSession(r)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		p := map[string]string{
			"title": "Login",
		}
		t, err := template.ParseFiles("html/login.html")
		if err != nil {
			printError(w, err)
			return
		}
		t.Execute(w, p)
		return
	}
	if r.Method == "POST" {
		var requestUser User
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		db, err := sql.Open("sqlite", "file:icas.db")
		if err != nil {
			fmt.Print(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT id, name, password FROM users WHERE email = ?", requestUser.Email)
		if err != nil {
			fmt.Print(err)
		}
		defer rows.Close()

		found := false
		var u User
		for rows.Next() {
			if err := rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
				fmt.Print(err)
			}
			if verifyPassword(requestUser.Password, u.Password) {
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
			store.sessions[sessionID] = u.Name
			store.mu.Unlock()
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(24 * time.Hour),
			})
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "found")
		} else {
			http.Error(w, "Usuário ou senha inválidos", http.StatusNotFound)
		}
		return
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, err := getActiveSession(r)
		if err != nil {
			//http.Error(w, fmt.Sprintf("Unauthorized: %s", err), http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getActiveSession(r *http.Request) (context.Context, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, errors.New("no session cookie")
	}
	store.mu.RLock()
	name, exists := store.sessions[cookie.Value]
	store.mu.RUnlock()

	if !exists {
		return nil, errors.New("invalid session")
	}

	ctx := context.WithValue(r.Context(), userContextKey, name)
	return ctx, nil
}
