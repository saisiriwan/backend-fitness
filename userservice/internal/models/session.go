package models

import "time"

// Session Log (หัวข้อการบันทึกผล)
type SessionLog struct {
	ID         int       `json:"id" db:"id"`
	ScheduleID int       `json:"schedule_id" db:"schedule_id"`
	ExerciseID *int      `json:"exercise_id" db:"exercise_id"` // อาจจะ null ได้
	Notes      string    `json:"notes" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	// (Optional) อาจจะมี Sets []SessionLogSet มาด้วยตอน GET แต่ตอน DB แยกตาราง
}

// Session Log Set (รายละเอียดแต่ละเซต)
type SessionLogSet struct {
	ID           int     `json:"id" db:"id"`
	SessionLogID int     `json:"session_log_id" db:"session_log_id"`
	SetNumber    int     `json:"set_number" db:"set_number"`
	WeightKg     float64 `json:"weight_kg" db:"weight_kg"`
	Reps         int     `json:"reps" db:"reps"`
	RPE          int     `json:"rpe" db:"rpe"`
}
