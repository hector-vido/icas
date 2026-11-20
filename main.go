package main

import (
	"fmt"
	"html/template"
	"os"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", homeHandler)
	http.HandleFunc("/animal/", animalHandler)
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func loadPage(fileName string) ([]byte, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return content, err
	}
	return content, nil
}

func printError(w http.ResponseWriter, err error) {
	fmt.Printf("internal error: %v\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//page, err := loadPage("html/index.html")
	p := map[string]string {
		"title": "Home",
	}
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		printError(w, err)
		return
	}
	t.Execute(w, p)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Printf("")
	}
}

func animalHandler(w http.ResponseWriter, r *http.Request) {
	//animal := r.URL.Query().Get("animal")
	animal := strings.TrimPrefix(r.URL.Path, "/animal/")
	fmt.Fprintf(w, "animal: %s", animal)
}
