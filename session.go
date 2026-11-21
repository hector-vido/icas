package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"icas/pkg/api"
	"icas/pkg/database"
	"icas/pkg/helpers"
	"net/http"
	"sync"
	"time"

	_ "modernc.org/sqlite" // sqlite driver
)

type contextKey string

const sessionContextKey contextKey = "sessionKey"

type sessionStore struct {
	mu       sync.RWMutex
	sessions map[string]api.User // sessionID -> email
}

var store = &sessionStore{
	sessions: make(map[string]api.User),
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
			helpers.PrintError(w, err)
			return
		}
		t.Execute(w, p)
		return
	}
	if r.Method == "POST" {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var requestUser api.User
		if err := decoder.Decode(&requestUser); err != nil {
			helpers.PrintError(w, err)
			return
		}
		user, err := database.UserLogin(requestUser.Email, requestUser.Password)
		if err != nil {
			http.Error(w, "Usuário ou senha inválidos", http.StatusNotFound)
		} else {
			sessionID, err := generateSessionID()
			if err != nil {
				helpers.PrintError(w, err)
				return
			}

			store.mu.Lock()
			store.sessions[sessionID] = user
			store.mu.Unlock()
			setCookie(w, sessionID)
			http.Error(w, "Usuário encontrado", http.StatusOK)
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
	sid, exists := store.sessions[cookie.Value]
	store.mu.RUnlock()

	if !exists {
		return nil, errors.New("invalid session")
	}

	ctx := context.WithValue(r.Context(), sessionContextKey, sid)
	return ctx, nil
}

func setCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
