package handler

import (
	"net/http"
	"strconv"
	"users/internal/models"
	"users/internal/repository"
	"users/internal/service" // เพิ่ม import service เพื่อเรียก GetUserByID (สำหรับดึงชื่อ Trainer)

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	repo        repository.ClientRepository
	userService service.UserService // เพิ่ม field นี้เพื่อดึงชื่อ Trainer
}

// ต้องแก้ NewClientHandler ให้รับ UserService เข้ามาด้วย
func NewClientHandler(repo repository.ClientRepository, userService service.UserService) *ClientHandler {
	return &ClientHandler{
		repo:        repo,
		userService: userService,
	}
}

// ... (ฟังก์ชันเดิม GetAllClients ... DeleteClient เหมือนเดิม) ...

// --- เพิ่มฟังก์ชันใหม่ต่อท้ายไฟล์ ---

// GET /api/v1/clients/:id/notes
func (h *ClientHandler) GetClientNotes(c *gin.Context) {
	clientID, _ := strconv.Atoi(c.Param("id"))

	notes, err := h.repo.GetNotesByClientID(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

// POST /api/v1/clients/:id/notes
func (h *ClientHandler) CreateClientNote(c *gin.Context) {
	clientID, _ := strconv.Atoi(c.Param("id"))

	var req models.ClientNote
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ดึงชื่อ Trainer จาก Token -> User ID -> Database
	userID, _ := c.Get("user_id")
	trainerID := int(userID.(float64))

	trainer, err := h.userService.GetUserByID(trainerID)
	trainerName := "Unknown Trainer"
	if err == nil {
		trainerName = trainer.Name
	}

	req.ClientID = clientID
	req.CreatedBy = trainerName // บันทึกชื่อคนเขียน

	if err := h.repo.CreateNote(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}
	c.JSON(http.StatusCreated, req)
}
