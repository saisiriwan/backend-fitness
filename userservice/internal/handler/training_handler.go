package handler

import (
	"net/http"
	"strconv"
	"users/internal/models"
	"users/internal/repository"

	"github.com/gin-gonic/gin"
)

type TrainingHandler struct {
	repo repository.TrainingRepository
}

func NewTrainingHandler(repo repository.TrainingRepository) *TrainingHandler {
	return &TrainingHandler{repo: repo}
}

// GET /api/v1/clients (เปลี่ยนชื่อจาก GetMyTrainees)
func (h *TrainingHandler) GetClients(c *gin.Context) {
	trainerID, _ := c.Get("user_id")
	id := int(trainerID.(float64))

	// เรียก Repo ใหม่ที่คืนค่า []models.Client
	clients, err := h.repo.GetClientsByTrainerID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clients)
}

// POST /api/v1/clients (Create Client)
func (h *TrainingHandler) CreateClient(c *gin.Context) {
	var req models.Client
	// Bind JSON ที่ส่งมาจากหน้าบ้าน
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ดึง ID ของ Trainer (คนที่ Login อยู่)
	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.CreateClient(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GET /api/v1/programs
func (h *TrainingHandler) GetPrograms(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	programs, err := h.repo.GetProgramsByUserID(int(userID.(float64)), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, programs)
}

// GET /api/v1/schedules
func (h *TrainingHandler) GetSchedules(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	schedules, err := h.repo.GetSchedulesByUserID(int(userID.(float64)), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// GET /api/v1/assignments
func (h *TrainingHandler) GetAssignments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	assignments, err := h.repo.GetAssignmentsByUserID(int(userID.(float64)), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// --- ส่วน POST (สร้างข้อมูล) ---

// POST /api/v1/programs
func (h *TrainingHandler) CreateProgram(c *gin.Context) {
	var req models.Program
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ดึง Trainer ID จาก Token (คนที่ Login อยู่คือคนสร้าง)
	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.CreateProgram(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create program"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// POST /api/v1/schedules
func (h *TrainingHandler) CreateSchedule(c *gin.Context) {
	var req models.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ดึง Trainer ID จาก Token
	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.CreateSchedule(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// POST /api/v1/assignments
func (h *TrainingHandler) CreateAssignment(c *gin.Context) {
	var req models.Assignment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ดึง Trainer ID จาก Token (คนที่ Login อยู่คือคนสร้าง)
	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	// ตรวจสอบว่า ClientID ถูกส่งมาหรือไม่
	if req.ClientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client ID is required for assignment"})
		return
	}

	// ตั้ง Status เป็น pending (ถ้า Frontend ไม่ได้ส่งมา)
	if req.Status == "" {
		req.Status = "pending"
	}

	// เรียก Repository เพื่อสร้าง Assignment
	if err := h.repo.CreateAssignment(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// PUT /api/v1/schedules/:id
func (h *TrainingHandler) UpdateSchedule(c *gin.Context) {
	var req models.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	req.ID = id

	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.UpdateSchedule(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// DELETE /api/v1/schedules/:id
func (h *TrainingHandler) DeleteSchedule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	trainerID, _ := c.Get("user_id")

	err := h.repo.DeleteSchedule(id, int(trainerID.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted"})
}

// PUT /api/v1/assignments/:id
func (h *TrainingHandler) UpdateAssignment(c *gin.Context) {
	var req models.Assignment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	req.ID = id

	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.UpdateAssignment(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// DELETE /api/v1/assignments/:id
func (h *TrainingHandler) DeleteAssignment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	trainerID, _ := c.Get("user_id")

	err := h.repo.DeleteAssignment(id, int(trainerID.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted"})
}
