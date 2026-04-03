package handler

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GraduationHandler struct {
	graduationService service.GraduationService
}

func NewGraduationHandler(graduationService service.GraduationService) *GraduationHandler {
	return &GraduationHandler{
		graduationService: graduationService,
	}
}

func (h *GraduationHandler) CheckGraduation(c *gin.Context) {
	var progress model.StudentProgress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.graduationService.CheckGraduation(&progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *GraduationHandler) CheckGraduationByName(c *gin.Context) {
	var request model.StudentProgressByName
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.graduationService.CheckGraduationByName(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *GraduationHandler) CheckCategoryRequirements(c *gin.Context) {
	var progress model.StudentProgress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.graduationService.CheckCategoryRequirements(&progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *GraduationHandler) ValidatePrerequisites(c *gin.Context) {
	var progress model.StudentProgress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violations, err := h.graduationService.ValidatePrerequisites(&progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"violations": violations})
}

func (h *GraduationHandler) ValidateCreditLimits(c *gin.Context) {
	var progress model.StudentProgress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	violations, err := h.graduationService.ValidateCreditLimits(&progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"violations": violations})
}
