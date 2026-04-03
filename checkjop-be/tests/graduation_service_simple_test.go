package tests

import (
	"checkjop-be/internal/model"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create complete mocks for courses
func setupCourseMocks(mockCourse *MockCourseRepository) {
	// Base courses with no prerequisites
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{
		Code:               "2301170",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups:  []model.PrerequisiteGroup{},
	}, nil)

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{
		Code:               "2301173",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups:  []model.PrerequisiteGroup{},
	}, nil)

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301220", mock.Anything, 2023).Return(&model.Course{
		Code:               "2301220",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups:  []model.PrerequisiteGroup{},
	}, nil)

	// Course 2301180 requires 2301170 OR 2301173
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301180", mock.Anything, 2023).Return(&model.Course{
		Code: "2301180",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{},
	}, nil)

	// Course 2301230 requires 2301220
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301230", mock.Anything, 2023).Return(&model.Course{
		Code: "2301230",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301220"}},
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{},
	}, nil)

	// Course 2301172 requires 2301170 as corequisite
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2023).Return(&model.Course{
		Code:               "2301172",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}, nil)

	// Transitive prerequisite chain: 2301260 -> nothing
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", mock.Anything, 2023).Return(&model.Course{
		Code:               "2301260",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups:  []model.PrerequisiteGroup{},
	}, nil)

	// 2301263 -> 2301260
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301263", mock.Anything, 2023).Return(&model.Course{
		Code: "2301263",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301260"}},
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{},
	}, nil)

	// 2301365 -> 2301263 (which creates the transitive chain)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301365", mock.Anything, 2023).Return(&model.Course{
		Code: "2301365",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301263"}},
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{},
	}, nil)
}

func TestValidatePrerequisites_Simple_NoPrerequisites(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_Simple_MissingPrerequisite(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301180 without its prerequisites (2301170 OR 2301173)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301180", Year: 2023, Semester: 1, Credits: 2},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301180", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")
	assert.Contains(t, violations[0].MissingPrereqs, "2301173")
}

func TestValidatePrerequisites_Simple_PrerequisiteSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301180 after taking 2301170 (satisfies OR requirement)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301180", Year: 2023, Semester: 2, Credits: 2},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_Simple_StrictTermRequirement(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301230 and its prerequisite 2301220 in same term (should fail strict requirement)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301220", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301230", Year: 2023, Semester: 1, Credits: 2}, // Same term - should fail
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301230", violations[0].CourseCode)
	assert.Empty(t, violations[0].MissingPrereqs)
	assert.Contains(t, violations[0].PrereqsTakenInWrongTerm, "2301220")
}

func TestValidatePrerequisites_Simple_TransitiveChain(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301365 without the transitive prerequisite chain: 2301365 -> 2301263 -> 2301260
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301365", Year: 2023, Semester: 1, Credits: 4},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301365", violations[0].CourseCode)
	// Should show all missing transitive prerequisites
	assert.Contains(t, violations[0].MissingPrereqs, "2301263")
	assert.Contains(t, violations[0].MissingPrereqs, "2301260")
}

func TestValidatePrerequisites_Simple_Corequisites_WrongTerm(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301172 and its corequisite 2301170 in different terms (should fail)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301172", Year: 2023, Semester: 2, Credits: 1}, // Different term - should fail
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301172", violations[0].CourseCode)
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301170")
}

func TestValidatePrerequisites_Simple_Corequisites_SameTerm(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301172 and its corequisite 2301170 in same term (should pass)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301172", Year: 2023, Semester: 1, Credits: 1}, // Same term - should pass
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidateCreditLimits_Simple_RegularSemester(t *testing.T) {
	graduationService, _, _, _ := setupGraduationService()

	// Test regular semester with over 22 credits
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 15},
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 8}, // Total: 23 credits
		},
	}

	violations, err := (*graduationService).ValidateCreditLimits(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, 2023, violations[0].Year)
	assert.Equal(t, 1, violations[0].Semester)
	assert.Equal(t, 23, violations[0].Credits)
	assert.Equal(t, 22, violations[0].MaxCredits)
}

func TestValidateCreditLimits_Simple_SummerSemester(t *testing.T) {
	graduationService, _, _, _ := setupGraduationService()

	// Test summer semester (semester 3) with over 10 credits
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 3, Credits: 8},
			{CourseCode: "2301173", Year: 2023, Semester: 3, Credits: 4}, // Total: 12 credits
		},
	}

	violations, err := (*graduationService).ValidateCreditLimits(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, 2023, violations[0].Year)
	assert.Equal(t, 3, violations[0].Semester) // Summer semester
	assert.Equal(t, 12, violations[0].Credits)
	assert.Equal(t, 10, violations[0].MaxCredits) // Summer limit
}

func TestValidatePrerequisites_Simple_CourseNotFound(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "INVALID", mock.Anything, 2023).Return((*model.Course)(nil), errors.New("course not found"))

	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "INVALID", Year: 2023, Semester: 1, Credits: 3},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err) // Should not error, just skip invalid courses
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_Simple_PartialTransitiveChain(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()
	setupCourseMocks(mockCourse)

	// Taking 2301365 with only part of the prerequisite chain satisfied
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301260", Year: 2023, Semester: 1, Credits: 4}, // Have base prerequisite
			{CourseCode: "2301365", Year: 2023, Semester: 3, Credits: 4}, // But missing 2301263
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301365", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301263")
	// Should not contain 2301260 since it's already taken
	assert.NotContains(t, violations[0].MissingPrereqs, "2301260")
}
