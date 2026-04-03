package main

import (
	"checkjop-be/internal/config"
	"checkjop-be/internal/database"
	"checkjop-be/internal/repository"
	"checkjop-be/internal/routes"
	"checkjop-be/internal/service"
	"checkjop-be/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	curriculumRepo := repository.NewCurriculumRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	setDefaultRepo := repository.NewSetDefaultRepository(db)

	// Initialize services
	curriculumService := service.NewCurriculumService(curriculumRepo)
	categoryService := service.NewCategoryService(categoryRepo, curriculumRepo)
	courseService := service.NewCourseService(courseRepo, categoryRepo, curriculumRepo)
	graduationService := service.NewGraduationService(curriculumRepo, courseRepo, categoryRepo)
	setDefaultService := service.NewSetDefaultService(setDefaultRepo, curriculumRepo, courseRepo)

	// Initialize Gin router
	router := gin.New()

	// Create rate limiter
	rateLimiter := middleware.CreateRateLimiter()

	// Apply middleware
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RateLimitMiddleware(rateLimiter))
	router.Use(middleware.ValidationErrorMiddleware())
	router.Use(gin.Recovery())

	// Add health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "CheckJop API is running"})
	})

	// Setup routes
	routes.SetupRoutes(router, curriculumService, categoryService, courseService, graduationService, setDefaultService)

	log.Printf("🚀 CheckJop Server starting on port %s", cfg.Port)
	log.Printf("📍 Health check: http://localhost:%s/health", cfg.Port)
	log.Printf("📍 API Documentation: http://localhost:%s/api/v1", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
