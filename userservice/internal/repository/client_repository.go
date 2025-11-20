package repository

import (
	"database/sql"
	"errors"
	"users/internal/models"
)

type ClientRepository interface {
	GetAllClients(trainerID int) ([]models.Client, error)
	CreateClient(client *models.Client) error
	GetClientByID(id int, trainerID int) (*models.Client, error)
	UpdateClient(client *models.Client) error
	DeleteClient(id int, trainerID int) error

	// Note methods
	GetNotesByClientID(clientID int) ([]models.ClientNote, error)
	CreateNote(note *models.ClientNote) error
}

// --- ส่วนที่ขาดหายไป ---
type clientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{db: db}
}

// -----------------------

// 1. Get All Clients
func (r *clientRepository) GetAllClients(trainerID int) ([]models.Client, error) {
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

// 2. Create Client
func (r *clientRepository) CreateClient(client *models.Client) error {
	query := `
		INSERT INTO clients (
			trainer_id, name, email, phone_number, 
			gender, height_cm, weight_kg, goal, birth_date, 
			injuries, activity_level, medical_conditions, avatar_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		client.TrainerID, client.Name, client.Email, client.Phone,
		client.Gender, client.Height, client.Weight, client.Goal, client.BirthDate,
		client.Injuries, client.ActivityLevel, client.MedicalConditions, client.AvatarURL,
	).Scan(&client.ID, &client.CreatedAt)
}

// 3. Get Client By ID
func (r *clientRepository) GetClientByID(id int, trainerID int) (*models.Client, error) {
	query := `
		SELECT id, trainer_id, name, email, phone_number, avatar_url, 
		       birth_date, gender, height_cm, weight_kg, goal, 
		       injuries, activity_level, medical_conditions, created_at
		FROM clients 
		WHERE id = $1 AND trainer_id = $2
	`
	var c models.Client
	err := r.db.QueryRow(query, id, trainerID).Scan(
		&c.ID, &c.TrainerID, &c.Name, &c.Email, &c.Phone, &c.AvatarURL,
		&c.BirthDate, &c.Gender, &c.Height, &c.Weight, &c.Goal,
		&c.Injuries, &c.ActivityLevel, &c.MedicalConditions, &c.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("client not found")
	}
	return &c, err
}

// 4. Update Client
func (r *clientRepository) UpdateClient(client *models.Client) error {
	query := `
		UPDATE clients 
		SET name=$1, email=$2, phone_number=$3, gender=$4, 
		    height_cm=$5, weight_kg=$6, goal=$7, birth_date=$8,
		    injuries=$9, activity_level=$10, medical_conditions=$11, avatar_url=$12
		WHERE id=$13 AND trainer_id=$14
	`
	res, err := r.db.Exec(query,
		client.Name, client.Email, client.Phone, client.Gender,
		client.Height, client.Weight, client.Goal, client.BirthDate,
		client.Injuries, client.ActivityLevel, client.MedicalConditions, client.AvatarURL,
		client.ID, client.TrainerID,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("client not found or unauthorized")
	}
	return nil
}

// 5. Delete Client
func (r *clientRepository) DeleteClient(id int, trainerID int) error {
	query := `DELETE FROM clients WHERE id=$1 AND trainer_id=$2`
	res, err := r.db.Exec(query, id, trainerID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("client not found or unauthorized")
	}
	return nil
}

// --- Notes Implementation ---

func (r *clientRepository) GetNotesByClientID(clientID int) ([]models.ClientNote, error) {
	query := `
		SELECT id, client_id, content, type, created_by, created_at 
		FROM client_notes 
		WHERE client_id = $1 
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.ClientNote
	for rows.Next() {
		var n models.ClientNote
		if err := rows.Scan(&n.ID, &n.ClientID, &n.Content, &n.Type, &n.CreatedBy, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *clientRepository) CreateNote(note *models.ClientNote) error {
	query := `
		INSERT INTO client_notes (client_id, content, type, created_by) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		note.ClientID,
		note.Content,
		note.Type,
		note.CreatedBy,
	).Scan(&note.ID, &note.CreatedAt)
}
