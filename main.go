package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", authMiddleware(homeHandler))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/animal/", authMiddleware(animalHandler))
	fmt.Println("Server listening on port 8080...")

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func printError(w http.ResponseWriter, err error) {
	fmt.Printf("internal error: %v\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"title": "Home",
	}
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		printError(w, err)
		return
	}
	name, _ := r.Context().Value(userContextKey).(string)
	fmt.Print(name)
	t.Execute(w, p)
}

func animalHandler(w http.ResponseWriter, r *http.Request) {
	//animal := r.URL.Query().Get("animal")
	animal := strings.TrimPrefix(r.URL.Path, "/animal/")
	fmt.Fprintf(w, "animal: %s", animal)
}
