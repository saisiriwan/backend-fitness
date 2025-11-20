package repository

import (
	"database/sql"
	"users/internal/models"
)

type ProgramRepository interface {
	// Program CRUD
	CreateProgram(p *models.Program) error
	GetProgramsByTrainerID(trainerID int) ([]models.Program, error)
	GetProgramByID(id int) (*models.Program, error)
	UpdateProgram(p *models.Program) error
	DeleteProgram(id int, trainerID int) error

	// Program Exercises
	AddExercise(pe *models.ProgramExercise) error
	GetExercisesByProgramID(programID int) ([]models.ProgramExercise, error)
	DeleteExerciseFromProgram(id int) error
}

type programRepository struct {
	db *sql.DB
}

func NewProgramRepository(db *sql.DB) ProgramRepository {
	return &programRepository{db: db}
}

// --- Implementation ---

func (r *programRepository) CreateProgram(p *models.Program) error {
	query := `
		INSERT INTO programs (name, description, trainer_id, client_id, is_template)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.db.QueryRow(query, p.Name, p.Description, p.TrainerID, p.ClientID, p.IsTemplate).
		Scan(&p.ID, &p.CreatedAt)
}

func (r *programRepository) GetProgramsByTrainerID(trainerID int) ([]models.Program, error) {
	query := `SELECT id, name, description, trainer_id, client_id, is_template, created_at FROM programs WHERE trainer_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, trainerID)
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

func (r *programRepository) GetProgramByID(id int) (*models.Program, error) {
	query := `SELECT id, name, description, trainer_id, client_id, is_template, created_at FROM programs WHERE id = $1`
	var p models.Program
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.TrainerID, &p.ClientID, &p.IsTemplate, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// --- Program Exercises ---

func (r *programRepository) AddExercise(pe *models.ProgramExercise) error {
	query := `
        INSERT INTO program_exercises (program_id, exercise_id, sets, reps, duration_seconds, rest_seconds, notes, "order")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`
	return r.db.QueryRow(query, pe.ProgramID, pe.ExerciseID, pe.Sets, pe.Reps, pe.DurationSeconds, pe.RestSeconds, pe.Notes, pe.Order).Scan(&pe.ID)
}

func (r *programRepository) GetExercisesByProgramID(programID int) ([]models.ProgramExercise, error) {
	query := `SELECT id, program_id, exercise_id, sets, reps, duration_seconds, rest_seconds, notes, "order" 
              FROM program_exercises WHERE program_id = $1 ORDER BY "order" ASC`
	rows, err := r.db.Query(query, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.ProgramExercise
	for rows.Next() {
		var pe models.ProgramExercise
		rows.Scan(&pe.ID, &pe.ProgramID, &pe.ExerciseID, &pe.Sets, &pe.Reps, &pe.DurationSeconds, &pe.RestSeconds, &pe.Notes, &pe.Order)
		exercises = append(exercises, pe)
	}
	return exercises, nil
}

func (r *programRepository) UpdateProgram(p *models.Program) error {
	query := `
		UPDATE programs 
		SET name=$1, description=$2, is_template=$3 
		WHERE id=$4 AND trainer_id=$5
	`
	// ใช้ Exec เพราะไม่ได้ต้องการ return ค่าอะไรกลับมา นอกจาก error หรือ rows affected
	res, err := r.db.Exec(query, p.Name, p.Description, p.IsTemplate, p.ID, p.TrainerID)
	if err != nil {
		return err
	}

	// เช็คว่ามีแถวที่ถูกแก้ไขจริงหรือไม่ (ถ้าเป็น 0 อาจแปลว่า ID ไม่ถูกต้อง หรือ Trainer ไม่ใช่เจ้าของ)
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// สร้าง error ใหม่ หรือจะ return nil ก็ได้แต่ควรบอกว่าหาไม่เจอ
		return sql.ErrNoRows
	}

	return nil
}

func (r *programRepository) DeleteProgram(id int, trainerID int) error {
	query := `DELETE FROM programs WHERE id=$1 AND trainer_id=$2`

	res, err := r.db.Exec(query, id, trainerID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// --- Program Exercises Delete ---

func (r *programRepository) DeleteExerciseFromProgram(id int) error {
	query := `DELETE FROM program_exercises WHERE id=$1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
