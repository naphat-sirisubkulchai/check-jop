package handler

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) Create(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.courseService.Create(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	course, err := h.courseService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Year is required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	course, err := h.courseService.GetByCodeAndYear(code, year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) GetByCurriculumID(c *gin.Context) {
	curriculumID, err := uuid.Parse(c.Param("curriculum_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid curriculum UUID"})
		return
	}

	courses, err := h.courseService.GetByCurriculumID(curriculumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) GetAll(c *gin.Context) {
	courses, err := h.courseService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.ID = id
	if err := h.courseService.Update(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := h.courseService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func (h *CourseHandler) TestRelationships(c *gin.Context) {
	// Check if both courses exist
	// Assuming year 2560 for testing
	year := 2560
	calcI, err := h.courseService.GetByCodeAndYear("2301107", year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calculus I not found: " + err.Error()})
		return
	}

	calcII, err := h.courseService.GetByCodeAndYear("2301108", year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Calculus II not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Both courses found",
		"calculus_i": map[string]interface{}{
			"id":            calcI.ID,
			"code":          calcI.Code,
			"curriculum_id": calcI.CurriculumID,
		},
		"calculus_ii": map[string]interface{}{
			"id":            calcII.ID,
			"code":          calcII.Code,
			"curriculum_id": calcII.CurriculumID,
		},
	})
}

// CheckCFOption checks if a course has C.F. (Consent of Faculty) option
func (h *CourseHandler) CheckCFOption(c *gin.Context) {
	code := c.Param("code")
	curriculumIDStr := c.Query("curriculum_id")
	yearStr := c.Query("year")

	if curriculumIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "curriculum_id is required"})
		return
	}

	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year is required"})
		return
	}

	curriculumID, err := uuid.Parse(curriculumIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid curriculum_id UUID"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	course, err := h.courseService.GetByCodeAndCurriculumIDAndYear(code, curriculumID, year)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Prepare prerequisite groups info
	prereqGroups := []map[string]interface{}{}
	for _, group := range course.PrerequisiteGroups {
		courseCodes := []string{}
		for _, link := range group.PrerequisiteCourses {
			courseCodes = append(courseCodes, link.PrerequisiteCourse.Code)
		}
		prereqGroups = append(prereqGroups, map[string]interface{}{
			"is_or_group":      group.IsOrGroup,
			"has_cf_condition": group.HasCFCondition,
			"course_codes":     courseCodes,
		})
	}

	// Prepare corequisite groups info
	coreqGroups := []map[string]interface{}{}
	for _, group := range course.CorequisiteGroups {
		courseCodes := []string{}
		for _, link := range group.PrerequisiteCourses {
			courseCodes = append(courseCodes, link.PrerequisiteCourse.Code)
		}
		coreqGroups = append(coreqGroups, map[string]interface{}{
			"is_or_group":      group.IsOrGroup,
			"has_cf_condition": group.HasCFCondition,
			"course_codes":     courseCodes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"course_code":         course.Code,
		"course_name_th":      course.NameTH,
		"course_name_en":      course.NameEN,
		"has_cf_option":       course.HasCFOption,
		"prerequisite_groups": prereqGroups,
		"corequisite_groups":  coreqGroups,
		"message": func() string {
			if course.HasCFOption {
				return "This course allows C.F. exemption"
			}
			return "This course does NOT allow C.F. exemption"
		}(),
	})
}
