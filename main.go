package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	user, _ := r.Context().Value(sessionContextKey).(User)
	p := map[string]string{
		"title":    "Home",
		"userId":   fmt.Sprintf("%d", user.ID),
		"userName": user.Name,
		"userRole": user.Role,
	}
	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		printError(w, err)
		return
	}
	t.Execute(w, p)
}

func animalHandler(w http.ResponseWriter, r *http.Request) {
	animal := strings.TrimPrefix(r.URL.Path, "/animal/")
	db, err := sql.Open("sqlite", "file:icas.db")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	var a Animal
	err = db.QueryRow("SELECT id, name, last_check, weight, equipment_date FROM animals WHERE name = ?", animal).Scan(&a.ID, &a.Name, &a.LastCheck, &a.Weight, &a.EquipmentDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Animal não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT latitude, longitude FROM animals_coordinates WHERE animal_id = ?", a.ID)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Coordinate
		if err := rows.Scan(&c.Lat, &c.Lng); err != nil {
			fmt.Print(err)
		}
		a.Coordinates = append(a.Coordinates, c)
	}
	if err := rows.Err(); err != nil {
		fmt.Print(err)
	}
	animalJson, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(animalJson))
}
