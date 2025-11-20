package handler

import (
	"net/http"
	"users/internal/repository"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	repo repository.DashboardRepository
}

func NewDashboardHandler(repo repository.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

// GET /api/v1/dashboard/stats
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	// ดึง Trainer ID จาก Token
	userID, _ := c.Get("user_id")
	trainerID := int(userID.(float64))

	stats, err := h.repo.GetDashboardStats(trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
