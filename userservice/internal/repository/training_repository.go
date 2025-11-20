// userservice/internal/repository/training_repository.go
package repository

import (
	"database/sql"
	"users/internal/models"
)

type TrainingRepository interface {
	GetClientsByTrainerID(trainerID int) ([]models.Client, error)
	CreateClient(client *models.Client) error
	GetProgramsByUserID(userID int, role string) ([]models.Program, error)
	GetSchedulesByUserID(userID int, role string) ([]models.Schedule, error)
	GetAssignmentsByUserID(userID int, role string) ([]models.Assignment, error)

	CreateAssignment(assignment *models.Assignment) error
	CreateProgram(program *models.Program) error
	CreateSchedule(schedule *models.Schedule) error

	UpdateSchedule(schedule *models.Schedule) error
	DeleteSchedule(id int, trainerID int) error

	UpdateAssignment(a *models.Assignment) error
	DeleteAssignment(id int, trainerID int) error
}

type trainingRepository struct {
	db *sql.DB
}

func NewTrainingRepository(db *sql.DB) TrainingRepository {
	return &trainingRepository{db: db}
}

// 1. ดึงรายชื่อลูกเทรน (Trainees) ของเทรนเนอร์คนนั้น
func (r *trainingRepository) GetClientsByTrainerID(trainerID int) ([]models.Client, error) {
	query := `
        SELECT id, trainer_id, name, email, phone_number, avatar_url, 
               birth_date, gender, height_cm, weight_kg, goal, 
               injuries, activity_level, medical_conditions, created_at
        FROM clients 
        WHERE trainer_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(
			&c.ID, &c.TrainerID, &c.Name, &c.Email, &c.Phone, &c.AvatarURL,
			&c.BirthDate, &c.Gender, &c.Height, &c.Weight, &c.Goal,
			&c.Injuries, &c.ActivityLevel, &c.MedicalConditions, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// 2. ดึง Program (ถ้าเป็น Trainer เห็นของที่ตัวเองสร้าง, Client เห็นของตัวเอง)
func (r *trainingRepository) GetProgramsByUserID(userID int, role string) ([]models.Program, error) {
	var query string
	if role == "trainer" {
		query = `SELECT id, name, description, trainer_id, client_id, is_template, created_at FROM programs WHERE trainer_id = $1`
	} else {
		query = `SELECT id, name, description, trainer_id, client_id, is_template, created_at FROM programs WHERE client_id = $1`
	}

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.TrainerID, &p.ClientID, &p.IsTemplate, &p.CreatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, nil
}

// 3. ดึง Schedules
func (r *trainingRepository) GetSchedulesByUserID(userID int, role string) ([]models.Schedule, error) {
	var query string
	if role == "trainer" {
		query = `SELECT id, title, trainer_id, client_id, start_time, end_time, status FROM schedules WHERE trainer_id = $1`
	} else {
		query = `SELECT id, title, trainer_id, client_id, start_time, end_time, status FROM schedules WHERE client_id = $1`
	}

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.Title, &s.TrainerID, &s.ClientID, &s.StartTime, &s.EndTime, &s.Status); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

// 4. ดึง Assignments
func (r *trainingRepository) GetAssignmentsByUserID(userID int, role string) ([]models.Assignment, error) {
	var query string
	if role == "trainer" {
		query = `SELECT id, title, description, client_id, trainer_id, due_date, status FROM assignments WHERE trainer_id = $1`
	} else {
		query = `SELECT id, title, description, client_id, trainer_id, due_date, status FROM assignments WHERE client_id = $1`
	}

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.Assignment
	for rows.Next() {
		var a models.Assignment
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.ClientID, &a.TrainerID, &a.DueDate, &a.Status); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}

// 5. สร้างโปรแกรมการฝึกใหม่
func (r *trainingRepository) CreateProgram(program *models.Program) error {
	query := `
		INSERT INTO programs (name, description, trainer_id, client_id, is_template)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at,updated_at
	`
	// ถ้า ClientID เป็น 0 หรือ nil ให้ส่ง nil เข้า DB
	return r.db.QueryRow(
		query,
		program.Name,
		program.Description,
		program.TrainerID,
		program.ClientID,
		program.IsTemplate,
	).Scan(&program.ID, &program.CreatedAt, &program.UpdatedAt)
}

// 6. สร้างตารางนัดหมายใหม่
func (r *trainingRepository) CreateSchedule(schedule *models.Schedule) error {
	query := `
		INSERT INTO schedules (title, trainer_id, client_id, start_time, end_time, status)
		VALUES ($1, $2, $3, $4, $5, 'scheduled')
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		schedule.Title,
		schedule.TrainerID,
		schedule.ClientID,
		schedule.StartTime,
		schedule.EndTime,
	).Scan(&schedule.ID, &schedule.CreatedAt)
}

// 7. สร้างงานมอบหมายใหม่ (Create Assignment)
func (r *trainingRepository) CreateAssignment(assignment *models.Assignment) error {
	query := `
		INSERT INTO assignments (title, description, client_id, trainer_id, due_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	// (แก้ไข) เปลี่ยนค่าที่ส่งไปให้ตรงกับ Assignment struct
	return r.db.QueryRow(
		query,
		assignment.Title,
		assignment.Description,
		assignment.ClientID,
		assignment.TrainerID,
		assignment.DueDate,
		assignment.Status,
	).Scan(&assignment.ID, &assignment.CreatedAt)
}

// Update Schedule
func (r *trainingRepository) UpdateSchedule(schedule *models.Schedule) error {
	query := `
		UPDATE schedules
		SET title=$1, client_id=$2, start_time=$3, end_time=$4, status=$5, updated_at=NOW()
		WHERE id=$6 AND trainer_id=$7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		schedule.Title,
		schedule.ClientID,
		schedule.StartTime,
		schedule.EndTime,
		schedule.Status,
		schedule.ID,
		schedule.TrainerID,
	).Scan(&schedule.UpdatedAt)
}

// Delete Schedule
func (r *trainingRepository) DeleteSchedule(id int, trainerID int) error {
	query := `DELETE FROM schedules WHERE id=$1 AND trainer_id=$2`
	res, err := r.db.Exec(query, id, trainerID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Update Assignment
func (r *trainingRepository) UpdateAssignment(a *models.Assignment) error {
	query := `
		UPDATE assignments
		SET title=$1, description=$2, client_id=$3, due_date=$4, status=$5, updated_at=NOW()
		WHERE id=$6 AND trainer_id=$7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		a.Title,
		a.Description,
		a.ClientID,
		a.DueDate,
		a.Status,
		a.ID,
		a.TrainerID,
	).Scan(&a.UpdatedAt)
}

// Delete Assignment
func (r *trainingRepository) DeleteAssignment(id int, trainerID int) error {
	query := `DELETE FROM assignments WHERE id=$1 AND trainer_id=$2`
	res, err := r.db.Exec(query, id, trainerID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// 8. สร้างลูกค้าใหม่ (Create Client)
func (r *trainingRepository) CreateClient(client *models.Client) error {
	query := `
		INSERT INTO clients (
            trainer_id, name, email, phone_number, 
            gender, height_cm, weight_kg, goal, birth_date
        )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		client.TrainerID, client.Name, client.Email, client.Phone,
		client.Gender, client.Height, client.Weight, client.Goal, client.BirthDate,
	).Scan(&client.ID, &client.CreatedAt)
}
