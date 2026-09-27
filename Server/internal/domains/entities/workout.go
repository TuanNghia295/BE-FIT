package entities

import (
	"time"

	"github.com/google/uuid"
)

type TrainingAvailability struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"column:userId;type:uuid;not null"`
	DayOfWeek int       `gorm:"column:dayOfWeek"`
	StartTime time.Time `gorm:"column:startTime;type:time"`
	EndTime   time.Time `gorm:"column:endTime;type:time"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (TrainingAvailability) TableName() string { return "training_availability" }

type WorkoutPlan struct {
	ID                 uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID             uuid.UUID  `gorm:"column:userId;type:uuid;not null"`
	TargetID           uuid.UUID  `gorm:"column:targetId;type:uuid;not null"`
	Version            int        `gorm:"column:version;default:1"`
	DurationPerSession int        `gorm:"column:durationPerSession"`
	StartDate          *time.Time `gorm:"column:startDate;type:date"`
	EndDate            *time.Time `gorm:"column:endDate;type:date"`
	Status             string     `gorm:"column:status"`
	CreatedAt          time.Time  `gorm:"column:createdAt"`
	UpdatedAt          time.Time  `gorm:"column:updatedAt"`
}

func (WorkoutPlan) TableName() string { return "workout_plan" }

type WorkoutSchedule struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkoutPlanID uuid.UUID `gorm:"column:workoutPlanId;type:uuid;not null"`
	DayOfWeek     int       `gorm:"column:dayOfWeek"`
	StartTime     time.Time `gorm:"column:startTime;type:time"`
	Duration      int       `gorm:"column:duration"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
}

func (WorkoutSchedule) TableName() string { return "workout_schedule" }

type Exercise struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `gorm:"column:name"`
	MuscleGroup  string    `gorm:"column:muscleGroup"`
	Instructions string    `gorm:"column:instructions;type:text"`
	Image        string    `gorm:"column:image"`
	Video        string    `gorm:"column:video"`
	Difficulty   string    `gorm:"column:difficulty"`
	Equipment    string    `gorm:"column:equipment"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
}

func (Exercise) TableName() string { return "exercise" }

type WorkoutExercise struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkoutPlanID uuid.UUID `gorm:"column:workoutPlanId;type:uuid;not null"`
	ExerciseID    uuid.UUID `gorm:"column:exerciseId;type:uuid;not null"`
	OrderIndex    int       `gorm:"column:orderIndex"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
}

func (WorkoutExercise) TableName() string { return "workout_exercise" }

type WorkoutSet struct {
	ID                uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkoutExerciseID uuid.UUID `gorm:"column:workoutExerciseId;type:uuid;not null"`
	SetNumber         int       `gorm:"column:setNumber"`
	TargetReps        int       `gorm:"column:targetReps"`
	TargetWeight      float64   `gorm:"column:targetWeight"`
	RestTime          int       `gorm:"column:restTime"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
	UpdatedAt         time.Time `gorm:"column:updatedAt"`
}

func (WorkoutSet) TableName() string { return "workout_set" }

type WorkoutSession struct {
	ID            uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID        uuid.UUID  `gorm:"column:userId;type:uuid;not null"`
	WorkoutPlanID uuid.UUID  `gorm:"column:workoutPlanId;type:uuid;not null"`
	ScheduleID    uuid.UUID  `gorm:"column:scheduleId;type:uuid;not null"`
	StartedAt     *time.Time `gorm:"column:startedAt"`
	CompletedAt   *time.Time `gorm:"column:completedAt"`
	Status        string     `gorm:"column:status"`
	CreatedAt     time.Time  `gorm:"column:createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt"`
}

func (WorkoutSession) TableName() string { return "workout_session" }

type ExercisePerformance struct {
	ID                uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkoutSessionID  uuid.UUID `gorm:"column:workoutSessionId;type:uuid;not null"`
	WorkoutExerciseID uuid.UUID `gorm:"column:workoutExerciseId;type:uuid;not null"`
	SetNumber         int       `gorm:"column:setNumber"`
	ActualReps        int       `gorm:"column:actualReps"`
	ActualWeight      float64   `gorm:"column:actualWeight"`
	RPE               float64   `gorm:"column:rpe"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
	UpdatedAt         time.Time `gorm:"column:updatedAt"`
}

func (ExercisePerformance) TableName() string { return "exercise_performance" }

type BodyProgress struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"column:userId;type:uuid;not null"`
	Weight     float64   `gorm:"column:weight"`
	RecordedAt time.Time `gorm:"column:recordedAt"`
	CreatedAt  time.Time `gorm:"column:createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt"`
}

func (BodyProgress) TableName() string { return "body_progress" }
