package handler

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SetDefaultHandler struct {
	setDefaultService service.SetDefaultService
}

func NewSetDefaultHandler(setDefaultService service.SetDefaultService) *SetDefaultHandler {
	return &SetDefaultHandler{
		setDefaultService: setDefaultService,
	}
}

func (h *SetDefaultHandler) Create(c *gin.Context) {
	var setDefault model.SetDefault
	if err := c.ShouldBindJSON(&setDefault); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.setDefaultService.Create(&setDefault); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, setDefault)
}

func (h *SetDefaultHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	setDefault, err := h.setDefaultService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SetDefault not found"})
		return
	}

	c.JSON(http.StatusOK, setDefault)
}

func (h *SetDefaultHandler) GetByCurriculumName(c *gin.Context) {
	curriculumName := c.Param("name")

	setDefaults, err := h.setDefaultService.GetByCurriculumName(curriculumName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setDefaults)
}

func (h *SetDefaultHandler) GetByCurriculumID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("curriculum_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	setDefaults, err := h.setDefaultService.GetByCurriculumID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setDefaults)
}

func (h *SetDefaultHandler) GetAll(c *gin.Context) {
	setDefaults, err := h.setDefaultService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setDefaults)
}

func (h *SetDefaultHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var setDefault model.SetDefault
	if err := c.ShouldBindJSON(&setDefault); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setDefault.ID = id
	if err := h.setDefaultService.Update(&setDefault); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setDefault)
}

func (h *SetDefaultHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := h.setDefaultService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SetDefault deleted successfully"})
}

func (h *SetDefaultHandler) ImportCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded or invalid file"})
		return
	}
	defer file.Close()

	if err := h.setDefaultService.ImportFromCSV(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SetDefault CSV imported successfully"})
}
