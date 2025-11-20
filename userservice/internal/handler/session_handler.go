package handler

import (
	"net/http"
	"strconv"
	"users/internal/models"
	"users/internal/repository"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	repo repository.SessionRepository
}

func NewSessionHandler(repo repository.SessionRepository) *SessionHandler {
	return &SessionHandler{repo: repo}
}

// POST /api/v1/sessions (สร้างนัดหมาย)
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req models.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	trainerID, _ := c.Get("user_id")
	req.TrainerID = int(trainerID.(float64))

	if err := h.repo.CreateSchedule(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// GET /api/v1/clients/:id/sessions (ดึงประวัติการนัดของลูกค้าคนนี้)
func (h *SessionHandler) GetClientSessions(c *gin.Context) {
	clientID, _ := strconv.Atoi(c.Param("id"))

	sessions, err := h.repo.GetSchedulesByClientID(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

// POST /api/v1/sessions/:id/logs (บันทึกผลการฝึก)
// (อันนี้ซับซ้อนหน่อย เพราะ Frontend อาจส่งมาเป็น Array ของ Exercises)
// เอาแบบง่ายก่อนคือรับทีละ Log
func (h *SessionHandler) CreateLog(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("id"))
	var req models.SessionLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	req.ScheduleID = scheduleID

	if err := h.repo.CreateSessionLog(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log session"})
		return
	}
	c.JSON(http.StatusCreated, req)
}
