package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"users/internal/config"
	"users/internal/models"

	_ "github.com/lib/pq"
)

// UserRepository interface ตามหลัก DIP
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Create(name, email string) (*models.User, error)
	Update(id int, name, email string) (*models.User, error)
	Delete(id int) error

	// (เพิ่มฟังก์ชันสำหรับ Auth)
	CreateUser(user models.User, hashedPassword string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	// (DSN) Data Source Name คือ String ที่ใช้ระบุข้อมูลการเชื่อมต่อกับฐานข้อมูล
	// การใช้งาน DSN จะขึ้นอยู่กับ Library หรือ Framework ที่ใช้งาน
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// ตั้งค่า Connection Pool
	db.SetMaxOpenConns(25)                 // จำนวน Connection สูงสุดที่สามารถเปิดได้
	db.SetMaxIdleConns(10)                 // จำนวน Connection สูงสุดที่สามารถอยู่ใน Idle State
	db.SetConnMaxLifetime(5 * time.Minute) // อายุการใช้งานสูงสุดของ Connection

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func CheckDBConnection(db *sql.DB) error {
	return db.Ping()
}

func (r *userRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(name, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at, updated_at",
		name, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(id int, name, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"UPDATE users SET name=$1, email=$2, updated_at=now() WHERE id=$3 RETURNING id, name, email, created_at, updated_at",
		name, email, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// CreateUser (สำหรับ Register)
func (r *userRepository) CreateUser(user models.User, hashedPassword string) (*models.User, error) {
	var u models.User

	// (แก้ไข SQL ให้ INSERT ลงคอลัมน์ใหม่ด้วย)
	err := r.db.QueryRow(
		"INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, name, email, role, created_at, updated_at",
		user.Name, user.Email, hashedPassword, user.Role,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByEmail (สำหรับ Login)
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User

	// (แก้ไข SQL ให้ SELECT คอลัมน์ใหม่มาด้วย)
	err := r.db.QueryRow("SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE email=$1", email).
		Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}
