package repository

import (
	"database/sql"
	"users/internal/models"
)

type SessionRepository interface {
	// Schedule
	CreateSchedule(s *models.Schedule) error
	GetSchedulesByClientID(clientID int) ([]models.Schedule, error)
	GetScheduleByID(id int) (*models.Schedule, error)
	UpdateScheduleStatus(id int, status string) error

	// Session Log
	CreateSessionLog(log *models.SessionLog) error
	CreateSessionLogSet(set *models.SessionLogSet) error
	GetLogsByScheduleID(scheduleID int) ([]models.SessionLog, error)
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

// --- Implementation ---

func (r *sessionRepository) CreateSchedule(s *models.Schedule) error {
	query := `
		INSERT INTO schedules (title, trainer_id, client_id, start_time, end_time, status)
		VALUES ($1, $2, $3, $4, $5, 'scheduled')
		RETURNING id, created_at`
	return r.db.QueryRow(query, s.Title, s.TrainerID, s.ClientID, s.StartTime, s.EndTime).
		Scan(&s.ID, &s.CreatedAt)
}

func (r *sessionRepository) GetSchedulesByClientID(clientID int) ([]models.Schedule, error) {
	query := `SELECT id, title, trainer_id, client_id, start_time, end_time, status, created_at 
              FROM schedules WHERE client_id = $1 ORDER BY start_time ASC`
	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.Title, &s.TrainerID, &s.ClientID, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// (ฟังก์ชัน GetScheduleByID, UpdateScheduleStatus เขียนคล้ายๆ กัน)
func (r *sessionRepository) UpdateScheduleStatus(id int, status string) error {
	query := `UPDATE schedules SET status=$1 WHERE id=$2`
	_, err := r.db.Exec(query, status, id)
	return err
}
func (r *sessionRepository) GetScheduleByID(id int) (*models.Schedule, error) {
	query := `SELECT id, title, trainer_id, client_id, start_time, end_time, status, created_at 
              FROM schedules WHERE id = $1`
	var s models.Schedule
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.TrainerID, &s.ClientID, &s.StartTime, &s.EndTime, &s.Status, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// --- Logs ---

func (r *sessionRepository) CreateSessionLog(log *models.SessionLog) error {
	query := `INSERT INTO session_logs (schedule_id, exercise_id, notes) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, log.ScheduleID, log.ExerciseID, log.Notes).Scan(&log.ID, &log.CreatedAt)
}

func (r *sessionRepository) CreateSessionLogSet(set *models.SessionLogSet) error {
	query := `INSERT INTO session_log_sets (session_log_id, set_number, weight_kg, reps, rpe) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, set.SessionLogID, set.SetNumber, set.WeightKg, set.Reps, set.RPE).Scan(&set.ID)
}

func (r *sessionRepository) GetLogsByScheduleID(scheduleID int) ([]models.SessionLog, error) {
	// (Logic ดึง Log และอาจต้อง Join เอา Sets มาด้วย ถ้าทำแบบละเอียด)
	// เบื้องต้นดึงแค่ Log header ก่อน
	query := `SELECT id, schedule_id, exercise_id, notes, created_at FROM session_logs WHERE schedule_id = $1`
	rows, err := r.db.Query(query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.SessionLog
	for rows.Next() {
		var l models.SessionLog
		rows.Scan(&l.ID, &l.ScheduleID, &l.ExerciseID, &l.Notes, &l.CreatedAt)
		logs = append(logs, l)
	}
	return logs, nil
}
