package handler

import (
	"checkjop-be/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CSVImportHandler struct {
	curriculumService service.CurriculumService
	categoryService   service.CategoryService
	courseService     service.CourseService
}

func NewCSVImportHandler(
	curriculumService service.CurriculumService,
	categoryService service.CategoryService,
	courseService service.CourseService,
) *CSVImportHandler {
	return &CSVImportHandler{
		curriculumService: curriculumService,
		categoryService:   categoryService,
		courseService:     courseService,
	}
}

func (h *CSVImportHandler) ImportCurriculumCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file"})
		return
	}
	defer file.Close()

	if err := h.curriculumService.ImportFromCSV(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curriculum CSV imported successfully"})
}

func (h *CSVImportHandler) ImportCategoryCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file"})
		return
	}
	defer file.Close()

	if err := h.categoryService.ImportFromCSV(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category CSV imported successfully"})
}

func (h *CSVImportHandler) ImportCourseCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file"})
		return
	}
	defer file.Close()

	if err := h.courseService.ImportFromCSV(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course CSV imported successfully"})
}

func (h *CSVImportHandler) ImportCourseCSVWithYear(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file"})
		return
	}
	defer file.Close()

	// Get year from request body or query parameter
	var yearParam struct {
		Year int `json:"year" form:"year" binding:"required"`
	}

	if err := c.ShouldBind(&yearParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Year parameter is required"})
		return
	}

	if err := h.courseService.ImportFromCSVWithYear(file, yearParam.Year); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Course CSV imported successfully for year %d", yearParam.Year)})
}

func (h *CSVImportHandler) ResetDatabase(c *gin.Context) {
	if err := h.courseService.ResetDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database reset successfully (Courses, Categories, Curricula cleared)"})
}
