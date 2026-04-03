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

func TestImportFromCSV_ComplexPrerequisites(t *testing.T) {
	_, mockCurriculum, mockCourse, mockCategory := setupGraduationService()
	courseService := service.NewCourseService(mockCourse, mockCategory, mockCurriculum)

	// CSV Content with (A AND B) OR C
	// A=101, B=102, C=103
	// Target=200
	csvContent := `code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
101,C1,C1,3,,,Cat1,Curr1,2023
102,C2,C2,3,,,Cat1,Curr1,2023
103,C3,C3,3,,,Cat1,Curr1,2023
200,Target,Target,3,(101 AND 102) OR 103,,Cat1,Curr1,2023`

	curr1ID := uuid.New()
	cat1ID := uuid.New()

	// Mock Curriculum
	mockCurriculum.On("GetByName", "Curr1").Return(&model.Curriculum{ID: curr1ID, NameEN: "Curr1"}, nil)

	// Mock Category
	mockCategory.On("GetByCurriculumID", curr1ID).Return([]model.Category{{ID: cat1ID, NameTH: "Cat1", NameEN: "Cat1"}}, nil)

	// Mock DeleteAll
	mockCourse.On("DeleteAll").Return(nil)

	// Mock BulkUpsert
	mockCourse.On("BulkUpsert", mock.Anything).Return(nil)

	// Mock Find Course
	course101 := &model.Course{ID: uuid.New(), Code: "101", Year: 2023, CurriculumID: curr1ID}
	course102 := &model.Course{ID: uuid.New(), Code: "102", Year: 2023, CurriculumID: curr1ID}
	course103 := &model.Course{ID: uuid.New(), Code: "103", Year: 2023, CurriculumID: curr1ID}
	course200 := &model.Course{ID: uuid.New(), Code: "200", Year: 2023, CurriculumID: curr1ID}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "101", curr1ID, 2023).Return(course101, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "102", curr1ID, 2023).Return(course102, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "103", curr1ID, 2023).Return(course103, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "200", curr1ID, 2023).Return(course200, nil)

	// Mock SetCorequisiteGroups (empty)
	mockCourse.On("SetCorequisiteGroups", mock.Anything, mock.Anything).Return(nil)

	// Verify Prerequisite Groups
	// Expected: (101 OR 103) AND (102 OR 103)
	// Two groups. Both IsOrGroup=true.
	// Group 1: [101, 103]
	// Group 2: [102, 103]
	// Order of groups and order within groups might vary, so we need flexible matching or sorting.

	mockCourse.On("SetPrerequisiteGroups", course200.ID, mock.MatchedBy(func(groups []model.PrerequisiteGroup) bool {
		if len(groups) != 2 {
			return false
		}

		// Check if we have the expected groups
		hasGroup1 := false // 101, 103
		hasGroup2 := false // 102, 103

		for _, g := range groups {
			if !g.IsOrGroup {
				return false
			}

			for _, link := range g.PrerequisiteCourses {
				// We can't easily check codes here because PrerequisiteCourses has IDs, not codes.
				// But we can check the count.
				// Wait, in the service, we append PrerequisiteCourseLink with ID.
				// We don't have the code in the link struct passed to repo?
				// model.PrerequisiteCourseLink usually has PrerequisiteCourseID.
				// We need to map IDs back to codes or check IDs.
				_ = link
			}
			// Since checking IDs is hard without a map, let's just check the structure size for now
			// or assume the IDs match the mocked objects.

			// Let's verify by ID
			has101 := false
			has102 := false
			has103 := false

			for _, link := range g.PrerequisiteCourses {
				if link.PrerequisiteCourseID == course101.ID {
					has101 = true
				}
				if link.PrerequisiteCourseID == course102.ID {
					has102 = true
				}
				if link.PrerequisiteCourseID == course103.ID {
					has103 = true
				}
			}

			if has101 && has103 && !has102 {
				hasGroup1 = true
			}
			if has102 && has103 && !has101 {
				hasGroup2 = true
			}
		}

		return hasGroup1 && hasGroup2
	})).Return(nil)

	// Execute
	err := courseService.ImportFromCSV(strings.NewReader(csvContent))

	// Assert
	assert.NoError(t, err)
}
