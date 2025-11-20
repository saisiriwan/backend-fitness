package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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

	// สร้าง Dependencies ใหม่
	trainingRepo := repository.NewTrainingRepository(db)
	trainingHandler := handler.NewTrainingHandler(trainingRepo)

	// --- Init Dashboard Components
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo)

	clientRepo := repository.NewClientRepository(db)
	clientHandler := handler.NewClientHandler(clientRepo, userService)

	sessionRepo := repository.NewSessionRepository(db)
	sessionHandler := handler.NewSessionHandler(sessionRepo)

	programRepo := repository.NewProgramRepository(db)
	programHandler := handler.NewProgramHandler(programRepo)

	r := gin.Default()
	// ----------------------------------------------------
	// 2. ใช้งาน CORS Middleware (ต้องอยู่ก่อน Routes)
	// ----------------------------------------------------
	r.Use(cors.New(cors.Config{
		// อนุญาต Origin (บ้าน) ของ Frontend
		AllowOrigins: []string{"http://localhost:3000"},
		// อนุญาต Methods (ท่า) ที่ Frontend ใช้
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// อนุญาต Headers ที่ Frontend ส่งมา
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// (สำคัญมาก!) อนุญาตให้ส่ง Cookie (JWT Token) ไปด้วย
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		if err := repository.CheckDBConnection(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"detail": "Database connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
	})

	// กลุ่มของ Auth Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userHandler.Register)
		authRoutes.POST("/login", userHandler.Login)
		authRoutes.POST("/logout", userHandler.Logout)
		authRoutes.GET("/google/login", userHandler.GoogleLogin)
		authRoutes.GET("/google/callback", userHandler.GoogleCallback)
	}

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JWTCookieAuth())
	{
		apiV1.DELETE("/users/:id", userHandler.DeleteUser)
		apiV1.PUT("/users/:id", userHandler.UpdateUser)
		apiV1.GET("/auth/me", userHandler.CheckAuth)
		apiV1.GET("/users", userHandler.GetAllUsers)
		apiV1.GET("/users/:id", userHandler.GetUserByID)
		// Training Routes (เพิ่มใหม่)

		apiV1.GET("/schedules", trainingHandler.GetSchedules)
		apiV1.POST("/schedules", trainingHandler.CreateSchedule)
		apiV1.PUT("/schedules/:id", trainingHandler.UpdateSchedule)
		apiV1.DELETE("/schedules/:id", trainingHandler.DeleteSchedule)

		apiV1.GET("/assignments", trainingHandler.GetAssignments)
		apiV1.POST("/assignments", trainingHandler.CreateAssignment)
		apiV1.PUT("/assignments/:id", trainingHandler.UpdateAssignment)
		apiV1.DELETE("/assignments/:id", trainingHandler.DeleteAssignment)

		apiV1.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)

		apiV1.GET("/clients", trainingHandler.GetClients)
		apiV1.POST("/clients", trainingHandler.CreateClient)

		apiV1.GET("/clients/:id/notes", clientHandler.GetClientNotes)
		apiV1.POST("/clients/:id/notes", clientHandler.CreateClientNote)

		apiV1.POST("/sessions", sessionHandler.CreateSession)
		apiV1.GET("/clients/:id/sessions", sessionHandler.GetClientSessions)
		apiV1.POST("/sessions/:id/logs", sessionHandler.CreateLog)

		apiV1.GET("/programs", programHandler.GetPrograms)
		apiV1.POST("/programs", programHandler.CreateProgram)
		apiV1.GET("/programs/:id", programHandler.GetProgramDetail)
		apiV1.POST("/programs/:id/exercises", programHandler.AddExercise)
		apiV1.PUT("/programs/:id", programHandler.UpdateProgram)
		apiV1.DELETE("/programs/:id", programHandler.DeleteProgram)

	}

	r.Run(":8080")
}
