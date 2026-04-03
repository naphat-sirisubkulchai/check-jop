package handler

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CurriculumHandler struct {
	curriculumService service.CurriculumService
}

func NewCurriculumHandler(curriculumService service.CurriculumService) *CurriculumHandler {
	return &CurriculumHandler{
		curriculumService: curriculumService,
	}
}

func (h *CurriculumHandler) Create(c *gin.Context) {
	var curriculum model.Curriculum
	if err := c.ShouldBindJSON(&curriculum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.curriculumService.Create(&curriculum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, curriculum)
}

func (h *CurriculumHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	curriculum, err := h.curriculumService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curriculum not found"})
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

func (h *CurriculumHandler) GetByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	curriculum, err := h.curriculumService.GetByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curriculum not found"})
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

func (h *CurriculumHandler) GetAll(c *gin.Context) {
	curricula, err := h.curriculumService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curricula)
}
func (h *CurriculumHandler) GetAllWithOutCatAndCourse(c *gin.Context) {
	curricula, err := h.curriculumService.GetAllWithOutCatAndCourse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curricula)
}

func (h *CurriculumHandler) GetActiveByYear(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	curricula, err := h.curriculumService.GetActiveByYear(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curricula)
}

func (h *CurriculumHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var curriculum model.Curriculum
	if err := c.ShouldBindJSON(&curriculum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	curriculum.ID = id
	if err := h.curriculumService.Update(&curriculum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, curriculum)
}

func (h *CurriculumHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := h.curriculumService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curriculum deleted successfully"})
}
