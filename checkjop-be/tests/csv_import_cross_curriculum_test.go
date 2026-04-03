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

func TestImportFromCSV_CrossCurriculumPrerequisite(t *testing.T) {
	_, mockCurriculum, mockCourse, mockCategory := setupGraduationService()
	courseService := service.NewCourseService(mockCourse, mockCategory, mockCurriculum)

	// CSV Content
	// code,nameEN,nameTH,credit,pre,co,category,curriculum,Year
	// 101,C1,C1,3,,,Cat1,Curr1,2023
	// 102,C2,C2,3,101,,Cat1,Curr2,2023
	csvContent := `code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
101,C1,C1,3,,,Cat1,Curr1,2023
102,C2,C2,3,101,,Cat1,Curr2,2023`

	curr1ID := uuid.New()
	curr2ID := uuid.New()
	cat1ID := uuid.New()

	// Mock Curriculum
	mockCurriculum.On("GetByName", "Curr1").Return(&model.Curriculum{ID: curr1ID, NameEN: "Curr1"}, nil)
	mockCurriculum.On("GetByName", "Curr2").Return(&model.Curriculum{ID: curr2ID, NameEN: "Curr2"}, nil)

	// Mock Category
	mockCategory.On("GetByCurriculumID", curr1ID).Return([]model.Category{{ID: cat1ID, NameTH: "Cat1", NameEN: "Cat1"}}, nil)
	mockCategory.On("GetByCurriculumID", curr2ID).Return([]model.Category{{ID: cat1ID, NameTH: "Cat1", NameEN: "Cat1"}}, nil)

	// Mock BulkUpsert (First pass)
	mockCourse.On("DeleteAll").Return(nil)
	mockCourse.On("BulkUpsert", mock.Anything).Return(nil)

	// Mock Find Course (Second pass - setup relationships)
	// 1. Find 101 in Curr1 (for itself) -> Found
	course101 := &model.Course{ID: uuid.New(), Code: "101", Year: 2023, CurriculumID: curr1ID}
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "101", curr1ID, 2023).Return(course101, nil)

	// 2. Find 102 in Curr2 (for itself) -> Found
	course102 := &model.Course{ID: uuid.New(), Code: "102", Year: 2023, CurriculumID: curr2ID}
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "102", curr2ID, 2023).Return(course102, nil)

	// 3. Find 101 (Prereq of 102) in Curr2 -> NOT FOUND (This is the failure point)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "101", curr2ID, 2023).Return((*model.Course)(nil), assert.AnError)

	// 4. Fallback: Find 101 by Code and Year -> Found (This is what we want to implement)
	mockCourse.On("GetByCodeAndYear", "101", 2023).Return(course101, nil)

	// Mock SetCourseRelationshipsWithGroups
	mockCourse.On("SetPrerequisiteGroups", mock.Anything, mock.Anything).Return(nil)
	mockCourse.On("SetCorequisiteGroups", mock.Anything, mock.Anything).Return(nil)

	// Execute
	err := courseService.ImportFromCSV(strings.NewReader(csvContent))

	// Assert
	// If the fix is working, this should be nil.
	// If not, it will fail because GetByCodeAndCurriculumIDAndYear returns error and we don't have fallback yet.
	// Actually, since I haven't implemented the fix yet, I expect this to FAIL if I run it now.
	// But I will implement the fix immediately after this.
	// For now, let's just assert error is nil assuming I will fix it.
	assert.NoError(t, err)
}
