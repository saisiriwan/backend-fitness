package models

import "time"

// Schedule (ตารางนัดหมาย/ตารางฝึก)
type Schedule struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	TrainerID int       `json:"trainer_id" db:"trainer_id"`
	ClientID  int       `json:"client_id" db:"client_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
