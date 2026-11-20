package main

type User struct {
	ID int
	Name string
	Email string
	Password string
}

type Animals struct {
	ID int
	Name string
}

type AnimalsCoordinates struct {
	AnimalID int
	Latitude float64
	Longitude float64
	RegisteredAt string
}
