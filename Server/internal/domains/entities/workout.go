package entities

import "time"

type TrainingAvailability struct {
	id        int64
	userId    int64
	dayOfWeek int // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	startTime time.Time
	endTime   time.Time
	createdAt time.Time
	updatedAt time.Time
}

type WorkoutPlan struct {
	id                 int64
	userId             int64
	targetId           int64
	version            int // default 1, dùng để track các thay đổi của plan
	durationPerSession int // workout time by minutes
	startDate          time.Time
	endDate            time.Time
	status             string // active, completed, archived
	createdAt          time.Time
	updatedAt          time.Time
}

type WorkoutSchedule struct {
	id            int64
	workoutPlanId int64
	dayOfWeek     int // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	startTime     time.Time
	duration      int
	createdAt     time.Time
	updatedAt     time.Time
}

type Exercise struct {
	id              int64
	name            string
	muscleGroup     string
	instructions    string
	image           string
	video           string
	difficultyLevel string // Beginner, Intermediate, Advanced
	equipment       string // Bodyweight, Dumbbell, Barbell, Machine, Resistance Band, Kettlebell
	createdAt       time.Time
	updatedAt       time.Time
}

type WorkoutExercise struct {
	id            int64
	workoutPLanId int64
	exerciseId    int64
	orderIndex    int64 // the order of the exercise in the workout plan
	createdAt     time.Time
	updatedAt     time.Time
}

type WorkoutSet struct {
	id                int64
	workoutExerciseId int64
	setNumber         int
	targetReps        int
	targetWeight      float64
	restTime          int // in seconds
	createdAt         time.Time
	updatedAt         time.Time
}

type WorkoutSession struct {
	id            int64
	userId        int64
	workoutPlanId int64
	scheduleId    int64
	startedAt     time.Time
	completedAt   time.Time
	status        string // planned, in_progress, completed, skipped
	createdAt     time.Time
	updatedAt     time.Time
}

type ExercisePerformance struct {
	id                int64
	workoutSessionId  int64
	workoutExerciseId int64
	setNumber         int

	// Chỉ hiển thị cho những người có level tập từ INTERMEDIATE
	actualReps   int
	actualWeight float64
	rpe          float64 // Rate of Perceived Exertion
	// -----
	createdAt time.Time
	updatedAt time.Time
}

type BodyProgress struct {
	id         int64
	userId     int64
	weight     float64
	recordedAt time.Time
	createdAt  time.Time
	updatedAt  time.Time
}
