package database

import (
	"database/sql"
	"errors"
	"icas/pkg/api"
)

func GetAnimal(animalName string) (api.Animal, error) {
	db, err := Open()
	if err != nil {
		return api.Animal{}, err
	}
	defer db.Close()

	var a api.Animal
	err = db.QueryRow("SELECT id, name, last_check, weight, equipment_date FROM animals WHERE name = ?", animalName).Scan(&a.ID, &a.Name, &a.LastCheck, &a.Weight, &a.EquipmentDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.Animal{}, errors.New("animal not found")
		}
		return api.Animal{}, err
	}

	rows, err := db.Query("SELECT latitude, longitude, registered_at FROM animals_coordinates WHERE animal_id = ?", a.ID)
	if err != nil {
		return api.Animal{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var c api.Coordinate
		if err := rows.Scan(&c.Lat, &c.Lng, &c.RegisteredAt); err != nil {
			return api.Animal{}, err
		}
		a.Coordinates = append(a.Coordinates, c)
	}
	if err := rows.Err(); err != nil {
		return api.Animal{}, err
	}
	return a, nil
}
