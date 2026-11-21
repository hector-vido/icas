package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"html/template"
	"icas/pkg/api"
	"icas/pkg/database"
	"icas/pkg/helpers"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func main() {

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	http.HandleFunc("/", authMiddleware(homeHandler))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/animal/", authMiddleware(animalHandler))
	http.HandleFunc("/gps/", authMiddleware(gpsHandler))

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	var err error
	var port string
	if os.Getenv("TLS_CERT") != "" && os.Getenv("TLS_KEY") != "" {
		port = cmp.Or(os.Getenv("PORT"), "8443")
		slog.Info("server listening", "port", port)
		err = http.ListenAndServeTLS(fmt.Sprintf(":%s", port), os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY"), nil)
	} else {
		port = cmp.Or(os.Getenv("PORT"), "8080")
		slog.Info("server listening", "port", port)
		err = http.ListenAndServe(fmt.Sprintf(":%s", cmp.Or(os.Getenv("PORT"), "8080")), nil)
	}
	if err != nil {
		slog.Error("error starting server", "error", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(sessionContextKey).(api.User)
	p := map[string]string{
		"title":    "Home",
		"userId":   fmt.Sprintf("%d", user.ID),
		"userName": user.Name,
		"userRole": user.Role,
	}
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		helpers.PrintError(w, err)
		return
	}
	t.Execute(w, p)
}

func animalHandler(w http.ResponseWriter, r *http.Request) {
	animalName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/animal/"), "/")
	animal, err := database.GetAnimal(animalName)
	if err != nil {
		helpers.PrintError(w, err)
	}
	helpers.PrintJSON(w, animal)
}

func gpsHandler(w http.ResponseWriter, r *http.Request) {
	animalName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/gps/"), "/")
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var coordinate api.Coordinate
	if err := decoder.Decode(&coordinate); err != nil {
		helpers.PrintError(w, err)
		return
	}
	if err := database.AddCoordinate(animalName, coordinate); err != nil {
		helpers.PrintError(w, err)
	}
}
