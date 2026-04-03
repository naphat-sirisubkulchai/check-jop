package tests

import (
	"checkjop-be/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetDatabase(t *testing.T) {
	_, mockCurriculum, mockCourse, mockCategory := setupGraduationService()
	courseService := service.NewCourseService(mockCourse, mockCategory, mockCurriculum)

	// Mock DeleteAll calls
	// Expect calls in reverse dependency order: Course -> Category -> Curriculum
	mockCourse.On("DeleteAll").Return(nil)
	mockCategory.On("DeleteAll").Return(nil)
	mockCurriculum.On("DeleteAll").Return(nil)

	// Execute
	err := courseService.ResetDatabase()

	// Assert
	assert.NoError(t, err)
	mockCourse.AssertExpectations(t)
	mockCategory.AssertExpectations(t)
	mockCurriculum.AssertExpectations(t)
}
