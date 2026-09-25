package entities

import "time"

type TrainingAvailability struct {
	ID        string    `gorm:"column:id"`
	UserID    string    `gorm:"column:userId"`
	DayOfWeek int       `gorm:"column:dayOfWeek"`
	StartTime time.Time `gorm:"column:startTime"`
	EndTime   time.Time `gorm:"column:endTime"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (TrainingAvailability) TableName() string {
	return "TrainingAvailability"
}

type WorkoutPlan struct {
	ID                 string    `gorm:"column:id"`
	UserID             string    `gorm:"column:userId"`
	TargetID           string    `gorm:"column:targetId"`
	Version            int       `gorm:"column:version"`
	DurationPerSession int       `gorm:"column:durationPerSession"`
	StartDate          time.Time `gorm:"column:startDate"`
	EndDate            time.Time `gorm:"column:endDate"`
	Status             string    `gorm:"column:status"`
	CreatedAt          time.Time `gorm:"column:createdAt"`
	UpdatedAt          time.Time `gorm:"column:updatedAt"`
}

func (WorkoutPlan) TableName() string {
	return "WorkoutPlan"
}

type WorkoutSchedule struct {
	ID            string    `gorm:"column:id"`
	WorkoutPlanID string    `gorm:"column:workoutPlanId"`
	DayOfWeek     int       `gorm:"column:dayOfWeek"`
	StartTime     time.Time `gorm:"column:startTime"`
	Duration      int       `gorm:"column:duration"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
}

func (WorkoutSchedule) TableName() string {
	return "WorkoutSchedule"
}

type Exercise struct {
	ID              string    `gorm:"column:id"`
	Name            string    `gorm:"column:name"`
	MuscleGroup     string    `gorm:"column:muscleGroup"`
	Instructions    string    `gorm:"column:instructions"`
	Image           string    `gorm:"column:image"`
	Video           string    `gorm:"column:video"`
	DifficultyLevel string    `gorm:"column:difficultyLevel"`
	Equipment       string    `gorm:"column:equipment"`
	CreatedAt       time.Time `gorm:"column:createdAt"`
	UpdatedAt       time.Time `gorm:"column:updatedAt"`
}

func (Exercise) TableName() string {
	return "Exercise"
}

type WorkoutExercise struct {
	ID            string    `gorm:"column:id"`
	WorkoutPlanID string    `gorm:"column:workoutPlanId"`
	ExerciseID    string    `gorm:"column:exerciseId"`
	OrderIndex    int64     `gorm:"column:orderIndex"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
}

func (WorkoutExercise) TableName() string {
	return "WorkoutExercise"
}

type WorkoutSet struct {
	ID                string    `gorm:"column:id"`
	WorkoutExerciseID string    `gorm:"column:workoutExerciseId"`
	SetNumber         int       `gorm:"column:setNumber"`
	TargetReps        int       `gorm:"column:targetReps"`
	TargetWeight      float64   `gorm:"column:targetWeight"`
	RestTime          int       `gorm:"column:restTime"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
	UpdatedAt         time.Time `gorm:"column:updatedAt"`
}

func (WorkoutSet) TableName() string {
	return "WorkoutSet"
}

type WorkoutSession struct {
	ID            string    `gorm:"column:id"`
	UserID        string    `gorm:"column:userId"`
	WorkoutPlanID string    `gorm:"column:workoutPlanId"`
	ScheduleID    string    `gorm:"column:scheduleId"`
	StartedAt     time.Time `gorm:"column:startedAt"`
	CompletedAt   time.Time `gorm:"column:completedAt"`
	Status        string    `gorm:"column:status"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
}

func (WorkoutSession) TableName() string {
	return "WorkoutSession"
}

type ExercisePerformance struct {
	ID                string `gorm:"column:id"`
	WorkoutSessionID  string `gorm:"column:workoutSessionId"`
	WorkoutExerciseID string `gorm:"column:workoutExerciseId"`
	SetNumber         int    `gorm:"column:setNumber"`

	// Chỉ hiển thị cho những người có level tập từ INTERMEDIATE
	ActualReps   int     `gorm:"column:actualReps"`
	ActualWeight float64 `gorm:"column:actualWeight"`
	RPE          float64 `gorm:"column:rpe"`
	// -----
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (ExercisePerformance) TableName() string {
	return "ExercisePerformance"
}

type BodyProgress struct {
	ID         string    `gorm:"column:id"`
	UserID     string    `gorm:"column:userId"`
	Weight     float64   `gorm:"column:weight"`
	RecordedAt time.Time `gorm:"column:recordedAt"`
	CreatedAt  time.Time `gorm:"column:createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt"`
}

func (BodyProgress) TableName() string {
	return "BodyProgress"
}
