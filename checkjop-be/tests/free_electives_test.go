package tests

import (
	"checkjop-be/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFreeElective_WithPrerequisites_ShouldCheck(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course 2301999 is a Free Elective but has a prerequisite
	course2301999 := &model.Course{
		Code: "2301999",
		Year: 2023,
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301999", mock.Anything, 2023).Return(course2301999, nil)

	// Student takes 2301999 as a Free Elective, but hasn't taken 2301170
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{
				CourseCode:   "2301999",
				Year:         2023,
				Semester:     1,
				Credits:      3,
				CategoryName: "Free Electives", // Explicitly marked as Free Elective
			},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301999", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")
}

func TestFreeElective_NotInDB_ShouldSkip(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course 9999999 is a Free Elective and NOT in the DB
	// Mock repository returns error for this course
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "9999999", mock.Anything, 2023).Return((*model.Course)(nil), assert.AnError)

	// Student takes 9999999 as a Free Elective
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{
				CourseCode:   "9999999",
				Year:         2023,
				Semester:     1,
				Credits:      3,
				CategoryName: "Free Electives",
			},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations) // Should be skipped, no violations
}

func TestFreeElective_NotInDB_CountsCredits(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	// Mock Curriculum
	mockCurriculum.On("GetByID", curriculumID).Return(&model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 3,
	}, nil)

	// Mock Category "Free Electives"
	categories := []model.Category{
		{
			NameTH:     "Free Electives",
			MinCredits: 3,
			Courses:    []model.Course{}, // No specific courses listed in DB for this category
		},
	}
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	// Course 9999999 is NOT in the DB
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "9999999", mock.Anything, 2023).Return((*model.Course)(nil), assert.AnError)

	// Student takes 9999999 as a Free Elective
	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{
				CourseCode:   "9999999",
				Year:         2023,
				Semester:     1,
				Credits:      3,
				CategoryName: "Free Electives",
			},
		},
	}

	// Check Graduation (which calls CheckCategoryRequirements)
	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.True(t, result.CanGraduate)
	assert.Equal(t, 3, result.TotalCredits)

	// Verify Category Result
	foundCategory := false
	for _, catResult := range result.CategoryResults {
		if catResult.CategoryName == "Free Electives" {
			foundCategory = true
			assert.Equal(t, 3, catResult.EarnedCredits)
			assert.True(t, catResult.IsSatisfied)
		}
	}
	assert.True(t, foundCategory, "Category 'Free Electives' should be present in results")
}
