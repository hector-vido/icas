package database

import (
	"database/sql"
	"errors"
	"fmt"
	"icas/pkg/api"
)

func AddCoordinate(animalName string, coordinate api.Coordinate) error {
	db, err := sql.Open("sqlite", "file:icas.db")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	var a api.Animal
	err = db.QueryRow("SELECT id FROM animals WHERE name = ?", animalName).Scan(&a.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("animal not found")
		}
		return err
	}

	_, err = db.Exec("INSERT INTO animals_coordinates (animal_id, latitude, longitude) VALUES (?, ?, ?)", a.ID, coordinate.Lat, coordinate.Lng)
	if err != nil {
		return err
	}
	return nil
}
