package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"users/internal/config"
	"users/internal/handler"
	"users/internal/middleware"
	"users/internal/repository"
	"users/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	// เชื่อมต่อ Database
	db, err := repository.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	// Health Check ไม่ต้องใช้ Token
	r.GET("/health", func(c *gin.Context) {
		if err := repository.CheckDBConnection(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"detail": "Database connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
	})

	// Protected routes ต้องใช้ Bearer Token
	authRequired := r.Group("/api/v1", middleware.BearerAuth(cfg.APIToken))
	{
		authRequired.GET("/users", userHandler.GetAllUsers)
		authRequired.GET("/users/:id", userHandler.GetUserByID)
		authRequired.POST("/users", userHandler.CreateUser)
		authRequired.PUT("/users/:id", userHandler.UpdateUser)
		authRequired.DELETE("/users/:id", userHandler.DeleteUser)
	}
	r.Run(":80")
}
