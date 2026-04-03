package tests

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImportFromCSV_DuplicateKeyDebug(t *testing.T) {
	_, mockCurriculum, mockCourse, mockCategory := setupGraduationService()
	courseService := service.NewCourseService(mockCourse, mockCategory, mockCurriculum)

	// CSV Content
	// Row 1 and Row 2 are identical -> Should be deduped
	csvContent := `code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
101,C1,C1,3,,,Cat1,Curr1,2023
101,C1,C1,3,,,Cat1,Curr1,2023`

	curr1ID := uuid.New()
	cat1ID := uuid.New()

	// Mock Curriculum
	mockCurriculum.On("GetByName", "Curr1").Return(&model.Curriculum{ID: curr1ID, NameEN: "Curr1"}, nil)

	// Mock Category
	mockCategory.On("GetByCurriculumID", curr1ID).Return([]model.Category{{ID: cat1ID, NameTH: "Cat1", NameEN: "Cat1"}}, nil)

	// Mock DeleteAll
	mockCourse.On("DeleteAll").Return(nil)

	// Mock BulkUpsert
	// We expect BulkUpsert to be called with a slice of length 1 (deduped)
	mockCourse.On("BulkUpsert", mock.MatchedBy(func(courses []model.Course) bool {
		return len(courses) == 1 && courses[0].Code == "101"
	})).Return(nil)

	// Mock Find Course (for relationships - none here)
	// No relationships, so no find calls needed for setupCourseRelationships

	// Execute
	err := courseService.ImportFromCSV(strings.NewReader(csvContent))

	// Assert
	assert.NoError(t, err)
}
