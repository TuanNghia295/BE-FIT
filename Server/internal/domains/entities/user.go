package entities

import "time"

type Users struct {
	ID        string    `gorm:"column:id" json:"id"`
	Email     string    `gorm:"column:email" json:"email"`
	FullName  string    `gorm:"column:fullName" json:"fullName"`
	Password  string    `gorm:"column:password" json:"-"` // hashed password
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (Users) TableName() string {
	return "Users"
}

type FitnessProfile struct {
	UserID    string    `gorm:"column:userId"`
	Gender    string    `gorm:"column:gender"`
	Age       int       `gorm:"column:age"`
	Height    float64   `gorm:"column:height"`
	Weight    float64   `gorm:"column:weight"`
	Level     string    `gorm:"column:level"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (FitnessProfile) TableName() string {
	return "FitnessProfile"
}

type Target struct {
	ID           string    `gorm:"column:id"`
	TargetWeight float64   `gorm:"column:targetWeight"`
	Goal         string    `gorm:"column:goal"`
	Calories     float64   `gorm:"column:calories"`
	Protein      float64   `gorm:"column:protein"`
	Carbs        float64   `gorm:"column:carbs"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
}

func (Target) TableName() string {
	return "Target"
}
