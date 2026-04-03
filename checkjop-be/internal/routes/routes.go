package routes

import (
	"checkjop-be/internal/handler"
	"checkjop-be/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	curriculumService service.CurriculumService,
	categoryService service.CategoryService,
	courseService service.CourseService,
	graduationService service.GraduationService,
	setDefaultService service.SetDefaultService,
) {
	// API v1 group - all endpoints are prefixed with /api/v1
	api := r.Group("/api/v1")

	// Initialize handlers
	curriculumHandler := handler.NewCurriculumHandler(curriculumService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	courseHandler := handler.NewCourseHandler(courseService)
	graduationHandler := handler.NewGraduationHandler(graduationService)
	setDefaultHandler := handler.NewSetDefaultHandler(setDefaultService)
	csvImportHandler := handler.NewCSVImportHandler(curriculumService, categoryService, courseService)

	// CURRICULUM MANAGEMENT ROUTES
	curriculum := api.Group("/curricula")
	{
		// CRUD operations for curricula
		curriculum.POST("/", curriculumHandler.Create)                             // Create new curriculum
		curriculum.GET("/", curriculumHandler.GetAll)                              // Get all curricula
		curriculum.GET("/:id", curriculumHandler.GetByID)                          // Get curriculum by ID
		curriculum.GET("/name/:name", curriculumHandler.GetByName)                 // Get curriculum by name
		curriculum.GET("/year/:year", curriculumHandler.GetActiveByYear)           // Get active curricula by year
		curriculum.PUT("/:id", curriculumHandler.Update)                           // Update curriculum
		curriculum.DELETE("/:id", curriculumHandler.Delete)                        // Delete curriculum
		curriculum.GET("/allwithout", curriculumHandler.GetAllWithOutCatAndCourse) // Get all curricula without categories and courses
	}

	// CATEGORY MANAGEMENT ROUTES
	category := api.Group("/categories")
	{
		// CRUD operations for categories (course categories within curricula)
		category.POST("/", categoryHandler.Create)                                    // Create new category
		category.GET("/", categoryHandler.GetAll)                                     // Get all categories
		category.GET("/:id", categoryHandler.GetByID)                                 // Get category by ID
		category.GET("/curriculum/:curriculum_id", categoryHandler.GetByCurriculumID) // Get categories by curriculum
		category.PUT("/:id", categoryHandler.Update)                                  // Update category
		category.DELETE("/:id", categoryHandler.Delete)                               // Delete category
	}

	// COURSE MANAGEMENT ROUTES
	course := api.Group("/courses")
	{
		// CRUD operations for courses
		course.POST("/", courseHandler.Create)                                    // Create new course
		course.GET("/", courseHandler.GetAll)                                     // Get all courses
		course.GET("/test-relationships", courseHandler.TestRelationships)        // Test relationship functionality
		course.GET("/:id", courseHandler.GetByID)                                 // Get course by ID
		course.GET("/code/:code", courseHandler.GetByCode)                        // Get course by course code
		course.GET("/code/:code/cf-option", courseHandler.CheckCFOption)          // Check if course has C.F. option
		course.GET("/curriculum/:curriculum_id", courseHandler.GetByCurriculumID) // Get courses by curriculum
		course.PUT("/:id", courseHandler.Update)                                  // Update course
		course.DELETE("/:id", courseHandler.Delete)                               // Delete course
	}

	// GRADUATION CHECKING ROUTES
	graduation := api.Group("/graduation")
	{
		// Graduation eligibility and requirement checking
		graduation.POST("/check", graduationHandler.CheckGraduation)                // Check overall graduation eligibility
		graduation.POST("/check/name", graduationHandler.CheckGraduationByName)     // Check graduation by curriculum name
		graduation.POST("/categories", graduationHandler.CheckCategoryRequirements) // Check category credit requirements
		graduation.POST("/prerequisites", graduationHandler.ValidatePrerequisites)  // Validate prerequisite/corequisite rules
		graduation.POST("/credit-limits", graduationHandler.ValidateCreditLimits)   // Validate credit limits per term
	}
	// SET DEFAULT MANAGEMENT ROUTES
	setDefault := api.Group("/set-defaults")
	{
		// CRUD operations for set defaults
		setDefault.POST("/", setDefaultHandler.Create)                                    // Create new set default
		setDefault.GET("/", setDefaultHandler.GetAll)                                     // Get all set defaults
		setDefault.GET("/:id", setDefaultHandler.GetByID)                                 // Get set default by ID
		setDefault.GET("/curriculum/name/:name", setDefaultHandler.GetByCurriculumName)   // Get set defaults by curriculum name
		setDefault.GET("/curriculum/:curriculum_id", setDefaultHandler.GetByCurriculumID) // Get set defaults by curriculum ID
		setDefault.PUT("/:id", setDefaultHandler.Update)                                  // Update set default
		setDefault.DELETE("/:id", setDefaultHandler.Delete)                               // Delete set default
	}

	// BULK CSV IMPORT ROUTES
	csvImport := api.Group("/import")
	{
		// Bulk data import via CSV files (all support curriculum names)
		csvImport.POST("/curriculum-csv", csvImportHandler.ImportCurriculumCSV)       // Import curricula from CSV
		csvImport.POST("/category-csv", csvImportHandler.ImportCategoryCSV)           // Import categories from CSV
		csvImport.POST("/course-csv", csvImportHandler.ImportCourseCSV)               // Import courses from CSV (with upsert, version 3)
		csvImport.POST("/course-csv-with-year", csvImportHandler.ImportCourseCSVWithYear) // Import courses from CSV with year parameter (version 4)
		csvImport.POST("/set-default-csv", setDefaultHandler.ImportCSV)               // Import set defaults from CSV
		csvImport.DELETE("/reset", csvImportHandler.ResetDatabase)                    // Reset database (clear all data)
	}
}
