package tests

import (
	"checkjop-be/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)



func TestCourseYearVersioning_DifferentVersions(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course 2301170 version 2023: No prerequisites
	course2023 := &model.Course{
		Code:               "2301170",
		Year:               2023,
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	// Course 2301170 version 2024: Requires 2301172
	course2024 := &model.Course{
		Code: "2301170",
		Year: 2024,
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301172"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2023, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2024).Return(course2024, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2024).Return(&model.Course{Code: "2301172"}, nil)

	// Student admitted in 2023 taking 2301170 -> Should use 2023 version (no prereq)
	progress2023 := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
		},
	}

	violations2023, err := (*graduationService).ValidatePrerequisites(progress2023)
	assert.NoError(t, err)
	assert.Empty(t, violations2023)

	// Student admitted in 2024 taking 2301170 -> Should use 2024 version (requires 2301172)
	progress2024 := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2024,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2024, Semester: 1, Credits: 3},
		},
	}

	violations2024, err := (*graduationService).ValidatePrerequisites(progress2024)
	assert.NoError(t, err)
	assert.Len(t, violations2024, 1)
	assert.Equal(t, "2301170", violations2024[0].CourseCode)
	assert.Contains(t, violations2024[0].MissingPrereqs, "2301172")
}
