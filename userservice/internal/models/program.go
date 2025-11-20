package models

import "time"

// Program (โปรแกรมการฝึก)
type Program struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	TrainerID   int       `json:"trainer_id" db:"trainer_id"`
	ClientID    *int      `json:"client_id" db:"client_id"` // null = template
	IsTemplate  bool      `json:"is_template" db:"is_template"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// (Optional) อาจจะมี Exercises []ProgramExercise มาด้วยตอน GET Detail
	Exercises []ProgramExercise `json:"exercises,omitempty"`
}

// Program Exercise (Detail)
type ProgramExercise struct {
	ID              int    `json:"id" db:"id"`
	ProgramID       int    `json:"program_id" db:"program_id"`
	ExerciseID      int    `json:"exercise_id" db:"exercise_id" binding:"required"`
	Sets            int    `json:"sets" db:"sets"`
	Reps            int    `json:"reps" db:"reps"`
	DurationSeconds int    `json:"duration_seconds" db:"duration_seconds"`
	RestSeconds     int    `json:"rest_seconds" db:"rest_seconds"`
	Notes           string `json:"notes" db:"notes"`
	Order           int    `json:"order" db:"order"`
}
