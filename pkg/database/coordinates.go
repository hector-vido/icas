package database

import (
	"database/sql"
	"errors"
	"icas/pkg/api"
)

func AddCoordinate(animalName string, coordinate api.Coordinate) error {
	db, err := sql.Open("sqlite", "file:icas.db")
	if err != nil {
		return err
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

func GetCoordinates(animalName string, limit int) ([]api.Coordinate, error) {
	db, err := Open()
	if err != nil {
		return []api.Coordinate{}, err
	}
	defer db.Close()

	var a api.Animal
	err = db.QueryRow("SELECT id FROM animals WHERE name = ?", animalName).Scan(&a.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []api.Coordinate{}, errors.New("animal not found")
		}
		return []api.Coordinate{}, err
	}

	rows, err := db.Query("SELECT latitude, longitude, registered_at FROM animals_coordinates WHERE animal_id = ? ORDER BY registered_at DESC LIMIT ?", a.ID, limit)
	if err != nil {
		return []api.Coordinate{}, err
	}
	defer rows.Close()

	var coordinates []api.Coordinate
	for rows.Next() {
		var c api.Coordinate
		if err := rows.Scan(&c.Lat, &c.Lng, &c.RegisteredAt); err != nil {
			return []api.Coordinate{}, err
		}
		coordinates = append(coordinates, c)
	}
	if err := rows.Err(); err != nil {
		return []api.Coordinate{}, err
	}
	return coordinates, nil
}
