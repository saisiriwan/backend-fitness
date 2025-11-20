package handler

import (
	"net/http"
	"strconv"
	"users/internal/models"
	"users/internal/repository"

	"github.com/gin-gonic/gin"
)

type ProgramHandler struct {
	repo repository.ProgramRepository
}

func NewProgramHandler(repo repository.ProgramRepository) *ProgramHandler {
	return &ProgramHandler{repo: repo}
}

// GET /api/v1/programs (ดึงรายการโปรแกรม)
func (h *ProgramHandler) GetPrograms(c *gin.Context) {
	trainerID, _ := c.Get("user_id")
	programs, err := h.repo.GetProgramsByTrainerID(int(trainerID.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch programs"})
		return
	}
	c.JSON(http.StatusOK, programs)
}

// GET /api/v1/programs/:id (ดึงรายละเอียด + ท่าฝึก)
func (h *ProgramHandler) GetProgramDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 1. ดึงข้อมูล Program
	program, err := h.repo.GetProgramByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
		return
	}

	// 2. ดึงข้อมูล Exercises
	exercises, err := h.repo.GetExercisesByProgramID(id)
	if err == nil {
		program.Exercises = exercises
	}

	c.JSON(http.StatusOK, program)
}

// PUT /api/v1/programs/:id (แก้ไขโปรแกรม)
func (h *ProgramHandler) UpdateProgram(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.Program
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	trainerID, _ := c.Get("user_id")
	req.ID = id
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.UpdateProgram(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update program"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Program updated successfully"})
}

// DELETE /api/v1/programs/:id (ลบโปรแกรม)
func (h *ProgramHandler) DeleteProgram(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	trainerID, _ := c.Get("user_id")

	if err := h.repo.DeleteProgram(id, int(trainerID.(float64))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete program"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Program deleted successfully"})
}

// POST /api/v1/programs (สร้างโปรแกรมใหม่)
func (h *ProgramHandler) CreateProgram(c *gin.Context) {
	var req models.Program
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.CreateProgram(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create program"})
		return
	}

	// (Optional) ถ้าส่ง Exercises มาด้วยใน JSON ก็วนลูปสร้างเลย

	c.JSON(http.StatusCreated, req)
}

// POST /api/v1/programs/:id/exercises (เพิ่มท่าฝึกในโปรแกรม)
func (h *ProgramHandler) AddExercise(c *gin.Context) {
	programID, _ := strconv.Atoi(c.Param("id"))
	var req models.ProgramExercise
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	req.ProgramID = programID

	if err := h.repo.AddExercise(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add exercise"})
		return
	}
	c.JSON(http.StatusCreated, req)
}
