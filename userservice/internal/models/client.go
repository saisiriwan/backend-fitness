package models

import "time"

// Client Struct ที่ตรงกับตาราง clients ใน init.sql ใหม่
type Client struct {
	ID        int `json:"id" db:"id"`
	TrainerID int `json:"trainer_id" db:"trainer_id"`

	// ข้อมูลส่วนตัว
	Name      string  `json:"name" db:"name" binding:"required"`
	Email     *string `json:"email" db:"email"` // อาจเป็น null ได้
	Phone     *string `json:"phone" db:"phone_number"`
	AvatarURL *string `json:"avatar" db:"avatar_url"`

	// ข้อมูลสุขภาพ (Profile เดิม)
	BirthDate         *time.Time `json:"birth_date" db:"birth_date"`
	Gender            *string    `json:"gender" db:"gender"`
	Height            *float64   `json:"height" db:"height_cm"`
	Weight            *float64   `json:"weight" db:"weight_kg"`
	Goal              *string    `json:"goal" db:"goal"`
	Injuries          *string    `json:"injuries" db:"injuries"`
	ActivityLevel     *string    `json:"activity_level" db:"activity_level"`
	MedicalConditions *string    `json:"medical_conditions" db:"medical_conditions"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
