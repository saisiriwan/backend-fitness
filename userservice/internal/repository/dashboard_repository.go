package repository

import (
	"database/sql"
	"users/internal/models"
)

type DashboardRepository interface {
	GetDashboardStats(trainerID int) (*models.DashboardStats, error)
}

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardStats(trainerID int) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	// 1. นับจำนวนลูกเทรน (Clients)
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM client_trainer_links WHERE trainer_id = $1
	`, trainerID).Scan(&stats.TotalClients)
	if err != nil {
		return nil, err
	}

	// 2. นับจำนวนโปรแกรม (Programs)
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM programs WHERE trainer_id = $1
	`, trainerID).Scan(&stats.ActivePrograms)
	if err != nil {
		return nil, err
	}

	// 3. นับนัดหมายที่ยังมาไม่ถึง (Upcoming Schedules)
	// (นับเฉพาะที่มี status='scheduled' และเวลาเริ่มยังมาไม่ถึง)
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM schedules 
		WHERE trainer_id = $1 AND status = 'scheduled' AND start_time > NOW()
	`, trainerID).Scan(&stats.UpcomingSession)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
