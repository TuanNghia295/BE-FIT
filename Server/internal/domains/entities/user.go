package entities

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Password  string    `json:"-"` // hashed password
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FitnessProfile struct {
	userId    int64
	gender    string
	age       int
	height    float64
	weight    float64
	level     string // Beginner, Intermediate, Advanced
	createdAt time.Time
	UpdatedAt time.Time
}

type Target struct {
	id           int64
	targetWeight float64
	goal         string //cutting, bulking, maintenance
	calories     float64
	protein      float64
	carbs        float64
	createdAt    time.Time
	updatedAt    time.Time
}
