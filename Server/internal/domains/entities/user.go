package entities

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"column:email;unique" json:"email"`
	FullName  string    `gorm:"column:fullName" json:"fullName"`
	Password  string    `gorm:"column:password;type:text" json:"-"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (Users) TableName() string { return "users" }

type FitnessProfile struct {
	UserID    uuid.UUID `gorm:"column:userId;type:uuid;primaryKey;not null"`
	Gender    string    `gorm:"column:gender"`
	Age       int       `gorm:"column:age"`
	Height    float64   `gorm:"column:height"`
	Weight    float64   `gorm:"column:weight"`
	Level     string    `gorm:"column:level"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (FitnessProfile) TableName() string { return "fitness_profile" }

type Target struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"column:userId;type:uuid;not null"`
	TargetWeight float64   `gorm:"column:targetWeight"`
	Goal         string    `gorm:"column:goal"`
	Calories     float64   `gorm:"column:calories"`
	Protein      float64   `gorm:"column:protein"`
	Carbs        float64   `gorm:"column:carbs"`
	Fat          float64   `gorm:"column:fat"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
}

func (Target) TableName() string { return "target" }
