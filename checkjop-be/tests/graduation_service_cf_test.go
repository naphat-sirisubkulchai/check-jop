package tests

import (
	"checkjop-be/internal/model"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidatePrerequisites_CFCondition_NoPermission(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR C.F."
	course2301199 := &model.Course{
		Code: "2301199",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup:      true,
				HasCFCondition: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)

	// Taking course without prerequisite and without C.F. permission
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301199", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301199", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")
}

func TestValidatePrerequisites_CFCondition_WithPermission(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR C.F."
	course2301199 := &model.Course{
		Code:        "2301199",
		HasCFOption: true, // Course allows C.F. exemption
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup:      true,
				HasCFCondition: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)

	// Taking course without prerequisite BUT WITH C.F. permission
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301199", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{"2301199"}, // Has permission for this course
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_CFCondition_WithPrerequisite(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR C.F."
	course2301199 := &model.Course{
		Code: "2301199",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup:      true,
				HasCFCondition: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)

	// Taking course WITH prerequisite (no need for C.F.)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301199", Year: 2023, Semester: 2, Credits: 3},
		},
		Exemptions: []string{},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_CFCondition_OnlyCF_NoPermission(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring ONLY "C.F." (no other course codes)
	course2301199 := &model.Course{
		Code: "2301199",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup:           false,
				HasCFCondition:      true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{}, // Empty - only C.F.
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)

	// Taking course without C.F. permission
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301199", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301199", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "C.F.")
}

func TestValidatePrerequisites_CFCondition_OnlyCF_WithPermission(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring ONLY "C.F." (no other course codes)
	course2301199 := &model.Course{
		Code:        "2301199",
		HasCFOption: true, // Course allows C.F. exemption
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup:           false,
				HasCFCondition:      true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{}, // Empty - only C.F.
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)

	// Taking course WITH C.F. permission
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301199", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{"2301199"}, // Has permission
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_CFCondition_WrongCourse(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course 2301170 has NO C.F. requirement (HasCFOption = false)
	course2301170 := &model.Course{
		Code:               "2301170",
		HasCFOption:        false,
		PrerequisiteGroups: []model.PrerequisiteGroup{}, // No prerequisites at all
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)

	// Taking course that doesn't allow C.F., but with C.F. exemption for it
	// This should cause an error during validation
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{"2301170"}, // Has exemption but course doesn't allow it
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not allow C.F. exemption")
	assert.Nil(t, violations)
}

func TestValidatePrerequisites_CFCondition_ValidExemption(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR C.F." (HasCFOption = true)
	course2301199 := &model.Course{
		Code:        "2301199",
		HasCFOption: true,
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup:      true,
				HasCFCondition: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)

	// Taking course WITH valid C.F. permission
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301199", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{"2301199"}, // Valid exemption
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_CFCondition_InvalidExemptionNonExistentCourse(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	course2301199 := &model.Course{
		Code:        "2301199",
		HasCFOption: false,
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301199", mock.Anything, 2023).Return(course2301199, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "NONEXISTENT", mock.Anything, 2023).Return(nil, fmt.Errorf("course not found"))

	// Exemption for a course that doesn't exist
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301199", Year: 2023, Semester: 1, Credits: 3},
		},
		Exemptions: []string{"NONEXISTENT"},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in curriculum")
	assert.Nil(t, violations)
}
