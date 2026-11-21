package api

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Animal struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	LastCheck     string       `json:"lastCheck"`
	Weight        float64      `json:"weight"`
	EquipmentDate string       `json:"equipmentDate"`
	Coordinates   []Coordinate `json:"coordinates"`
}

type AnimalCoordinate struct {
	AnimalID     int     `json:"animal_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RegisteredAt string  `json:"registered_at"`
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
