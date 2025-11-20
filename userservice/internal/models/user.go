package models

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	PasswordHash string    `db:"password_hash" json:"-"` // json:"-" คือ ห้ามส่งฟิลด์นี้กลับไปใน JSON
	Role         string    `db:"role" json:"role"`
	AvatarURL    *string   `db:"avatar_url" json:"avatar_url"`
}
