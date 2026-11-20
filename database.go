package main

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Animals struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AnimalsCoordinates struct {
	AnimalID     int     `json:"animal_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RegisteredAt string  `json:"registered_at"`
}
