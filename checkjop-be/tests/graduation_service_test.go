package tests

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGraduationService() (*service.GraduationService, *MockCurriculumRepository, *MockCourseRepository, *MockCategoryRepository) {
	mockCurriculum := &MockCurriculumRepository{}
	mockCourse := &MockCourseRepository{}
	mockCategory := &MockCategoryRepository{}

	graduationService := service.NewGraduationService(mockCurriculum, mockCourse, mockCategory)

	return &graduationService, mockCurriculum, mockCourse, mockCategory
}

func TestValidatePrerequisites_NoPrerequisites(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course 2301170 (Computer and Programming) has no prerequisites based on PDF
	course2301170 := &model.Course{
		Code:               "2301170",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups:  []model.PrerequisiteGroup{},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)

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

func TestValidatePrerequisites_BasicPrerequisite_2301180(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// From PDF: 2301180 requires "2301170 OR 2301173" as prerequisite
	course2301180 := &model.Course{
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
	}

	// Mock the prerequisite courses too (for transitive checking)
	course2301170 := &model.Course{
		Code:               "2301170",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	course2301173 := &model.Course{
		Code:               "2301173",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301180", mock.Anything, 2023).Return(course2301180, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(course2301173, nil)

	// Test case: Taking 2301180 without any prerequisites
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
	assert.Empty(t, violations[0].PrereqsTakenInWrongTerm)
}

func TestValidatePrerequisites_PrerequisiteSatisfied_2301180(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// 2301180 requires 2301170 OR 2301173
	course2301180 := &model.Course{
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
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301180", mock.Anything, 2023).Return(course2301180, nil)

	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Prerequisite satisfied
			{CourseCode: "2301180", Year: 2023, Semester: 2, Credits: 2}, // Course taken after
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_StrictTermRequirement_2301230(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// From PDF: 2301230 (Discrete Math for CS) requires 2301220 as prerequisite
	course2301230 := &model.Course{
		Code: "2301230",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301220"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301220", mock.Anything, 2023).Return(&model.Course{Code: "2301220"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301230", mock.Anything, 2023).Return(course2301230, nil)

	// Test case: Taking both courses in same term (should violate strict requirement)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301220", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301230", Year: 2023, Semester: 1, Credits: 2}, // Same term violation
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301230", violations[0].CourseCode)
	assert.Empty(t, violations[0].MissingPrereqs)
	assert.Contains(t, violations[0].PrereqsTakenInWrongTerm, "2301220")
}

func TestValidatePrerequisites_TransitiveChain_DataStructures(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Set up transitive chain based on PDF:
	// 2301365 (Algorithm Design & Analysis) -> 2301263 (Data Structures & Algorithms) -> 2301260 (Programming Techniques)
	course2301260 := &model.Course{
		Code:               "2301260",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	course2301263 := &model.Course{
		Code: "2301263",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301260"}},
				},
			},
		},
	}

	course2301365 := &model.Course{
		Code: "2301365",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301263"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", mock.Anything, 2023).Return(course2301260, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301263", mock.Anything, 2023).Return(course2301263, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301365", mock.Anything, 2023).Return(course2301365, nil)

	// Test: Taking 2301365 without any prerequisites
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

func TestValidatePrerequisites_ComplexPrerequisites_2301367(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// From PDF: 2301367 (Software Engineering Methods) requires 2301375
	course2301367 := &model.Course{
		Code: "2301367",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301375"}},
				},
			},
		},
	}

	course2301375 := &model.Course{
		Code: "2301375",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301263"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301263", mock.Anything, 2023).Return(&model.Course{Code: "2301263"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301375", mock.Anything, 2023).Return(course2301375, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301367", mock.Anything, 2023).Return(course2301367, nil)

	// Test partial prerequisite chain
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301263", Year: 2023, Semester: 1, Credits: 4}, // Have this prerequisite
			{CourseCode: "2301367", Year: 2023, Semester: 3, Credits: 3}, // But missing direct prerequisite 2301375
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301367", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301375")
}

func TestValidateCreditLimits_RegularSemester(t *testing.T) {
	graduationService, _, _, _ := setupGraduationService()

	// Test regular semester with over 22 credits
	progress := &model.StudentProgress{
		CurriculumID: uuid.New(),
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 10},
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 8},
			{CourseCode: "2301180", Year: 2023, Semester: 1, Credits: 5}, // Total: 23 credits
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

func TestValidateCreditLimits_SummerSemester(t *testing.T) {
	graduationService, _, _, _ := setupGraduationService()

	// Test summer semester (semester 3) with over 10 credits
	progress := &model.StudentProgress{
		CurriculumID: uuid.New(),
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 3, Credits: 6},
			{CourseCode: "2301173", Year: 2023, Semester: 3, Credits: 6}, // Total: 12 credits (over summer limit of 10)
		},
	}

	violations, err := (*graduationService).ValidateCreditLimits(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, 2023, violations[0].Year)
	assert.Equal(t, 3, violations[0].Semester)
	assert.Equal(t, 12, violations[0].Credits)
	assert.Equal(t, 10, violations[0].MaxCredits) // Summer semester limit
}

func TestValidatePrerequisites_Corequisites_2301172(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// From PDF: 2301172 requires 2301170 as corequisite (same term)
	course2301172 := &model.Course{
		Code: "2301172",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2023).Return(course2301172, nil)

	// Test: Taking corequisites in different terms (should violate)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301172", Year: 2023, Semester: 2, Credits: 1}, // Different term
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301172", violations[0].CourseCode)
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301170")
}

func TestValidatePrerequisites_Corequisites_SameTerm_Valid(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// 2301172 requires 2301170 as corequisite
	course2301172 := &model.Course{
		Code: "2301172",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2023).Return(course2301172, nil)

	// Test: Taking corequisites in same term (should be valid)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301172", Year: 2023, Semester: 1, Credits: 1}, // Same term
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestCheckCategoryRequirements(t *testing.T) {
	graduationService, _, _, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			NameEN:     "General Education",
			MinCredits: 30,
			Courses: []model.Course{
				{Code: "2301170", Credits: 3},
				{Code: "2301173", Credits: 4},
			},
		},
		{
			NameTH:     "วิชาเฉพาะ",
			NameEN:     "Major Courses",
			MinCredits: 90,
			Courses: []model.Course{
				{Code: "2301260", Credits: 4},
				{Code: "2301263", Credits: 4},
			},
		},
	}

	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID: curriculumID,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 4},
			{CourseCode: "2301260", Year: 2023, Semester: 2, Credits: 4},
		},
	}

	results, err := (*graduationService).CheckCategoryRequirements(progress)

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// General Education category
	genEd := results[0]
	assert.Equal(t, "วิชาศึกษาทั่วไป", genEd.CategoryName)
	assert.Equal(t, 7, genEd.EarnedCredits) // 3 + 4
	assert.Equal(t, 30, genEd.RequiredCredits)
	assert.False(t, genEd.IsSatisfied) // 7 < 30

	// Major courses category
	major := results[1]
	assert.Equal(t, "วิชาเฉพาะ", major.CategoryName)
	assert.Equal(t, 4, major.EarnedCredits) // Only 2301260
	assert.Equal(t, 90, major.RequiredCredits)
	assert.False(t, major.IsSatisfied) // 4 < 90
}

func TestCheckGraduation_Complete(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 120,
	}

	// Simple courses with no prerequisites
	course2301170 := &model.Course{Code: "2301170", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}
	course2301173 := &model.Course{Code: "2301173", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}
	course2301260 := &model.Course{Code: "2301260", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}

	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			MinCredits: 30,
			Courses:    []model.Course{{Code: "2301170", Credits: 30}, {Code: "2301173", Credits: 30}},
		},
		{
			NameTH:     "วิชาเฉพาะ",
			MinCredits: 90,
			Courses:    []model.Course{{Code: "2301260", Credits: 90}},
		},
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(course2301173, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", mock.Anything, 2023).Return(course2301260, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 15},
			{CourseCode: "2301173", Year: 2023, Semester: 2, Credits: 15},
			{CourseCode: "2301260", Year: 2024, Semester: 1, Credits: 22}, // Major category part 1
			{CourseCode: "2301260", Year: 2024, Semester: 2, Credits: 22}, // Major category part 2
			{CourseCode: "2301260", Year: 2025, Semester: 1, Credits: 22}, // Major category part 3
			{CourseCode: "2301260", Year: 2025, Semester: 2, Credits: 22}, // Major category part 4
			{CourseCode: "2301260", Year: 2026, Semester: 1, Credits: 2},  // Complete major category (90 total)
		},
		ManualCredits: map[string]int{
			"วิชาศึกษาทั่วไป": 15, // Add 15 more to reach 30 required for general education
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.True(t, result.CanGraduate)
	assert.Equal(t, 120, result.TotalCredits)
	assert.Equal(t, 120, result.RequiredCredits)
	assert.Empty(t, result.PrerequisiteViolations)
	assert.Empty(t, result.CreditLimitViolations)
}

func TestValidatePrerequisites_CourseNotFound(t *testing.T) {
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

func TestValidatePrerequisites_MultipleViolationTypes(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course with both missing prerequisites and wrong-term corequisites
	course2301260 := &model.Course{
		Code: "2301260",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301172"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2023).Return(&model.Course{Code: "2301172"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", mock.Anything, 2023).Return(course2301260, nil)

	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301172", Year: 2023, Semester: 1, Credits: 1}, // Corequisite in different term
			{CourseCode: "2301260", Year: 2023, Semester: 2, Credits: 4}, // Missing prerequisite 2301170
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301260", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")         // Missing prerequisite
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301172") // Wrong term corequisite
	assert.True(t, violations[0].TakenInWrongTerm)
}

// ==================== COMPREHENSIVE COREQUISITE TESTS ====================

func TestValidatePrerequisites_Corequisites_MissingCorequisite(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course that requires a corequisite
	course2301185 := &model.Course{
		Code: "2301185",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301185", mock.Anything, 2023).Return(course2301185, nil)

	// Taking course without its required corequisite
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301185", Year: 2023, Semester: 1, Credits: 3}, // Missing corequisite 2301170
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301185", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingCoreqs, "2301170")
	assert.Empty(t, violations[0].MissingPrereqs)
}

func TestValidatePrerequisites_Corequisites_OrGroup_OneSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR 2301173" as corequisite
	course2301190 := &model.Course{
		Code: "2301190",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301190", mock.Anything, 2023).Return(course2301190, nil)

	// Taking course with one of the OR corequisites (should pass)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Satisfies OR requirement
			{CourseCode: "2301190", Year: 2023, Semester: 1, Credits: 2}, // Same term
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_Corequisites_OrGroup_NoneSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR 2301173" as corequisite
	course2301190 := &model.Course{
		Code: "2301190",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301190", mock.Anything, 2023).Return(course2301190, nil)

	// Taking course without any of the OR corequisites
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301190", Year: 2023, Semester: 1, Credits: 2}, // Missing both OR corequisites
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301190", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingCoreqs, "2301170")
	assert.Contains(t, violations[0].MissingCoreqs, "2301173")
}

func TestValidatePrerequisites_Corequisites_OrGroup_WrongTerm(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "2301170 OR 2301173" as corequisite
	course2301190 := &model.Course{
		Code: "2301190",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301190", mock.Anything, 2023).Return(course2301190, nil)

	// Taking course with corequisites in different terms
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Different term
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 3}, // Different term
			{CourseCode: "2301190", Year: 2023, Semester: 2, Credits: 2}, // Different from both
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301190", violations[0].CourseCode)
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301170")
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301173")
}

func TestValidatePrerequisites_Corequisites_AndGroup_AllSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring both "2301170 AND 2301173" as corequisites
	course2301195 := &model.Course{
		Code: "2301195",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301195", mock.Anything, 2023).Return(course2301195, nil)

	// Taking course with both required corequisites in same term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301195", Year: 2023, Semester: 1, Credits: 2}, // Same term as both
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_Corequisites_AndGroup_OneMissing(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring both "2301170 AND 2301173" as corequisites
	course2301195 := &model.Course{
		Code: "2301195",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301195", mock.Anything, 2023).Return(course2301195, nil)

	// Taking course with only one of the required AND corequisites
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Have this one
			{CourseCode: "2301195", Year: 2023, Semester: 1, Credits: 2}, // Missing 2301173
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301195", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingCoreqs, "2301173")
	assert.NotContains(t, violations[0].MissingCoreqs, "2301170") // This one is taken
}

func TestValidatePrerequisites_Corequisites_TransitiveCorequisites(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Chain: 2301195 requires 2301172, and 2301172 requires 2301170
	course2301170 := &model.Course{Code: "2301170"}
	course2301172 := &model.Course{
		Code: "2301172",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}
	course2301195 := &model.Course{
		Code: "2301195",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301172"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2023).Return(course2301172, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301195", mock.Anything, 2023).Return(course2301195, nil)

	// Taking 2301195 with complete transitive corequisite chain (should pass)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Base corequisite
			{CourseCode: "2301172", Year: 2023, Semester: 1, Credits: 1}, // Intermediate corequisite
			{CourseCode: "2301195", Year: 2023, Semester: 1, Credits: 2}, // Target course
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_Corequisites_TransitiveMissing(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Chain: 2301195 requires 2301172, and 2301172 requires 2301170
	course2301170 := &model.Course{Code: "2301170"}
	course2301172 := &model.Course{
		Code: "2301172",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
				},
			},
		},
	}
	course2301195 := &model.Course{
		Code: "2301195",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301172"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2023).Return(course2301172, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301195", mock.Anything, 2023).Return(course2301195, nil)

	// Taking 2301195 and 2301172 but missing the transitive corequisite 2301170
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301172", Year: 2023, Semester: 1, Credits: 1}, // Missing its corequisite 2301170
			{CourseCode: "2301195", Year: 2023, Semester: 1, Credits: 2}, // Has direct corequisite but missing transitive
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 2) // Both courses violate corequisite requirements

	// Find violations by course code
	var violation2301172, violation2301195 *model.PrerequisiteViolation
	for i := range violations {
		if violations[i].CourseCode == "2301172" {
			violation2301172 = &violations[i]
		} else if violations[i].CourseCode == "2301195" {
			violation2301195 = &violations[i]
		}
	}

	// 2301172 violates its direct corequisite requirement
	assert.NotNil(t, violation2301172)
	assert.Contains(t, violation2301172.MissingCoreqs, "2301170")

	// 2301195 also violates because its corequisite 2301172 doesn't have its corequisite
	assert.NotNil(t, violation2301195)
	assert.Contains(t, violation2301195.MissingCoreqs, "2301170")
}

func TestValidatePrerequisites_Corequisites_MixedPrereqAndCoreq(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course with both prerequisites and corequisites
	course2301200 := &model.Course{
		Code: "2301200",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}}, // Must be taken before
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301173"}}, // Must be taken same term
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301200", mock.Anything, 2023).Return(course2301200, nil)

	// Proper scenario: prerequisite taken before, corequisite taken same term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Prerequisite (before)
			{CourseCode: "2301173", Year: 2023, Semester: 2, Credits: 3}, // Corequisite (same term)
			{CourseCode: "2301200", Year: 2023, Semester: 2, Credits: 4}, // Target course
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

// ==================== COMPLEX PREREQUISITE TESTS ====================

func TestValidatePrerequisites_ComplexOrGroups_BothGroupsSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) - both groups required
	course2301400 := &model.Course{
		Code: "2301400",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock prerequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301400", mock.Anything, 2023).Return(course2301400, nil)

	// Taking course with one from each OR group (should pass)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301265", Year: 2023, Semester: 1, Credits: 3}, // Satisfies first OR group
			{CourseCode: "2301369", Year: 2023, Semester: 2, Credits: 3}, // Satisfies second OR group
			{CourseCode: "2301400", Year: 2024, Semester: 1, Credits: 4}, // Target course
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_ComplexOrGroups_OnlyFirstGroupSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) - both groups required
	course2301400 := &model.Course{
		Code: "2301400",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock prerequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301400", mock.Anything, 2023).Return(course2301400, nil)

	// Taking course with only first OR group satisfied
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301274", Year: 2023, Semester: 1, Credits: 3}, // Satisfies first OR group
			{CourseCode: "2301400", Year: 2024, Semester: 1, Credits: 4}, // Missing second OR group
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301400", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301279")
	assert.Contains(t, violations[0].MissingPrereqs, "2301369")
	assert.NotContains(t, violations[0].MissingPrereqs, "2301265") // First group satisfied
	assert.NotContains(t, violations[0].MissingPrereqs, "2301274") // First group satisfied
}

func TestValidatePrerequisites_ComplexOrGroups_NoGroupsSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) - both groups required
	course2301400 := &model.Course{
		Code: "2301400",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock prerequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301400", mock.Anything, 2023).Return(course2301400, nil)

	// Taking course with no prerequisites satisfied
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301400", Year: 2024, Semester: 1, Credits: 4}, // Missing all prerequisites
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301400", violations[0].CourseCode)
	// Should show all missing courses from both OR groups
	assert.Contains(t, violations[0].MissingPrereqs, "2301265")
	assert.Contains(t, violations[0].MissingPrereqs, "2301274")
	assert.Contains(t, violations[0].MissingPrereqs, "2301279")
	assert.Contains(t, violations[0].MissingPrereqs, "2301369")
}

func TestValidatePrerequisites_MixedOrAndRequirements_BothSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), 2301279 - OR group and single requirement
	course2301500 := &model.Course{
		Code: "2301500",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
				},
			},
		},
	}

	// Mock prerequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301500", mock.Anything, 2023).Return(course2301500, nil)

	// Taking course with OR group satisfied and single requirement satisfied
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301274", Year: 2023, Semester: 1, Credits: 3}, // Satisfies OR group
			{CourseCode: "2301279", Year: 2023, Semester: 2, Credits: 3}, // Satisfies single requirement
			{CourseCode: "2301500", Year: 2024, Semester: 1, Credits: 4}, // Target course
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_MixedOrAndRequirements_OrGroupMissing(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), 2301279 - OR group and single requirement
	course2301500 := &model.Course{
		Code: "2301500",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
				},
			},
		},
	}

	// Mock prerequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301500", mock.Anything, 2023).Return(course2301500, nil)

	// Taking course with single requirement satisfied but OR group missing
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301279", Year: 2023, Semester: 2, Credits: 3}, // Single requirement satisfied
			{CourseCode: "2301500", Year: 2024, Semester: 1, Credits: 4}, // OR group missing
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301500", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301265")
	assert.Contains(t, violations[0].MissingPrereqs, "2301274")
	assert.NotContains(t, violations[0].MissingPrereqs, "2301279") // Single requirement satisfied
}

func TestValidatePrerequisites_MixedOrAndRequirements_SingleRequirementMissing(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), 2301279 - OR group and single requirement
	course2301500 := &model.Course{
		Code: "2301500",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
				},
			},
		},
	}

	// Mock prerequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301500", mock.Anything, 2023).Return(course2301500, nil)

	// Taking course with OR group satisfied but single requirement missing
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301265", Year: 2023, Semester: 1, Credits: 3}, // OR group satisfied
			{CourseCode: "2301500", Year: 2024, Semester: 1, Credits: 4}, // Single requirement missing
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301500", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301279")
	assert.NotContains(t, violations[0].MissingPrereqs, "2301265") // OR group satisfied
	assert.NotContains(t, violations[0].MissingPrereqs, "2301274") // OR group satisfied
}

// ==================== COMPLEX COREQUISITE TESTS ====================

func TestValidatePrerequisites_ComplexCorequisiteOrGroups_BothGroupsSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) as corequisites - both groups required
	course2301600 := &model.Course{
		Code: "2301600",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301600", mock.Anything, 2023).Return(course2301600, nil)

	// Taking course with one from each OR group in same term (should pass)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301265", Year: 2023, Semester: 1, Credits: 3}, // Satisfies first OR group
			{CourseCode: "2301369", Year: 2023, Semester: 1, Credits: 3}, // Satisfies second OR group
			{CourseCode: "2301600", Year: 2023, Semester: 1, Credits: 4}, // Target course (same term)
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_ComplexCorequisiteOrGroups_OnlyFirstGroupSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) as corequisites - both groups required
	course2301600 := &model.Course{
		Code: "2301600",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301600", mock.Anything, 2023).Return(course2301600, nil)

	// Taking course with only first OR group satisfied in same term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301274", Year: 2023, Semester: 1, Credits: 3}, // Satisfies first OR group
			{CourseCode: "2301600", Year: 2023, Semester: 1, Credits: 4}, // Missing second OR group
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301600", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingCoreqs, "2301279")
	assert.Contains(t, violations[0].MissingCoreqs, "2301369")
	assert.NotContains(t, violations[0].MissingCoreqs, "2301265") // First group satisfied
	assert.NotContains(t, violations[0].MissingCoreqs, "2301274") // First group satisfied
}

func TestValidatePrerequisites_MixedCorequisiteOrAndRequirements_BothSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), 2301279 as corequisites - OR group and single requirement
	course2301700 := &model.Course{
		Code: "2301700",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301700", mock.Anything, 2023).Return(course2301700, nil)

	// Taking course with OR group satisfied and single requirement satisfied in same term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301274", Year: 2023, Semester: 1, Credits: 3}, // Satisfies OR group
			{CourseCode: "2301279", Year: 2023, Semester: 1, Credits: 3}, // Satisfies single requirement
			{CourseCode: "2301700", Year: 2023, Semester: 1, Credits: 4}, // Target course (same term)
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_MixedCorequisiteOrAndRequirements_OrGroupMissing(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), 2301279 as corequisites - OR group and single requirement
	course2301700 := &model.Course{
		Code: "2301700",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301700", mock.Anything, 2023).Return(course2301700, nil)

	// Taking course with single requirement satisfied but OR group missing in same term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301279", Year: 2023, Semester: 1, Credits: 3}, // Single requirement satisfied
			{CourseCode: "2301700", Year: 2023, Semester: 1, Credits: 4}, // OR group missing
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301700", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingCoreqs, "2301265")
	assert.Contains(t, violations[0].MissingCoreqs, "2301274")
	assert.NotContains(t, violations[0].MissingCoreqs, "2301279") // Single requirement satisfied
}

func TestValidatePrerequisites_MixedCorequisiteOrAndRequirements_WrongTerm(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), 2301279 as corequisites - OR group and single requirement
	course2301700 := &model.Course{
		Code: "2301700",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301700", mock.Anything, 2023).Return(course2301700, nil)

	// Taking course with all corequisites taken in different terms
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301274", Year: 2023, Semester: 1, Credits: 3}, // Different term
			{CourseCode: "2301279", Year: 2023, Semester: 1, Credits: 3}, // Different term
			{CourseCode: "2301700", Year: 2023, Semester: 2, Credits: 4}, // Target course (different term)
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301700", violations[0].CourseCode)
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301274")
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301279")
}

func TestValidatePrerequisites_ComplexCorequisiteOrGroups_BothGroupsWrongTerm(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) as corequisites - both groups required
	course2301600 := &model.Course{
		Code: "2301600",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301600", mock.Anything, 2023).Return(course2301600, nil)

	// Taking course with both OR groups satisfied but in different terms
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301265", Year: 2023, Semester: 1, Credits: 3}, // Different term
			{CourseCode: "2301369", Year: 2023, Semester: 1, Credits: 3}, // Different term
			{CourseCode: "2301600", Year: 2023, Semester: 2, Credits: 4}, // Target course (different term)
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301600", violations[0].CourseCode)
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301265")
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301369")
}

func TestValidatePrerequisites_ComplexCorequisiteOrGroups_NoGroupsSatisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring (2301265 OR 2301274), (2301279 OR 2301369) as corequisites - both groups required
	course2301600 := &model.Course{
		Code: "2301600",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301265"}},
					{PrerequisiteCourse: model.Course{Code: "2301274"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	// Mock corequisite courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301265", mock.Anything, 2023).Return(&model.Course{Code: "2301265"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301274", mock.Anything, 2023).Return(&model.Course{Code: "2301274"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2023).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2023).Return(&model.Course{Code: "2301369"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301600", mock.Anything, 2023).Return(course2301600, nil)

	// Taking course without any corequisites
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301600", Year: 2023, Semester: 1, Credits: 4}, // No corequisites at all
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301600", violations[0].CourseCode)
	// Should show all missing courses from both OR groups
	assert.Contains(t, violations[0].MissingCoreqs, "2301265")
	assert.Contains(t, violations[0].MissingCoreqs, "2301274")
	assert.Contains(t, violations[0].MissingCoreqs, "2301279")
	assert.Contains(t, violations[0].MissingCoreqs, "2301369")
}

func TestValidatePrerequisites_StrictTermRequirement_Satisfied(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// From PDF: 2301230 (Discrete Math for CS) requires 2301220 as prerequisite
	course2301230 := &model.Course{
		Code: "2301230",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301220"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301220", mock.Anything, 2023).Return(&model.Course{Code: "2301220"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301230", mock.Anything, 2023).Return(course2301230, nil)

	// Test case: Taking prerequisite in earlier term (should pass)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301220", Year: 2023, Semester: 1, Credits: 3},
			{CourseCode: "2301230", Year: 2023, Semester: 2, Credits: 2}, // Different term - should pass
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidateCreditLimits_RegularSemester_WithinLimit(t *testing.T) {
	graduationService, _, _, _ := setupGraduationService()

	// Test regular semester with exactly 22 credits (at limit)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 12},
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 10}, // Total: 22 credits (at limit)
		},
	}

	violations, err := (*graduationService).ValidateCreditLimits(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidateCreditLimits_SummerSemester_WithinLimit(t *testing.T) {
	graduationService, _, _, _ := setupGraduationService()

	// Test summer semester (semester 3) with exactly 10 credits (at limit)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 3, Credits: 6},
			{CourseCode: "2301173", Year: 2023, Semester: 3, Credits: 4}, // Total: 10 credits (at limit)
		},
	}

	violations, err := (*graduationService).ValidateCreditLimits(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestValidatePrerequisites_MixedPrereqCoreq_MissingPrerequisite(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course with both prerequisites and corequisites
	course2301200 := &model.Course{
		Code: "2301200",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}}, // Must be taken before
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301173"}}, // Must be taken same term
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301200", mock.Anything, 2023).Return(course2301200, nil)

	// Missing prerequisite but corequisite taken in same term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301173", Year: 2023, Semester: 2, Credits: 3}, // Corequisite (same term)
			{CourseCode: "2301200", Year: 2023, Semester: 2, Credits: 4}, // Missing prerequisite 2301170
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301200", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")
	assert.Empty(t, violations[0].CoreqsTakenInWrongTerm) // Corequisite is correct
}

func TestValidatePrerequisites_MixedPrereqCoreq_CorequisiteWrongTerm(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course with both prerequisites and corequisites
	course2301200 := &model.Course{
		Code: "2301200",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}}, // Must be taken before
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301173"}}, // Must be taken same term
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301200", mock.Anything, 2023).Return(course2301200, nil)

	// Prerequisite taken before but corequisite in wrong term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 3}, // Prerequisite (before)
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 3}, // Corequisite (wrong term)
			{CourseCode: "2301200", Year: 2023, Semester: 2, Credits: 4}, // Target course
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301200", violations[0].CourseCode)
	assert.Empty(t, violations[0].MissingPrereqs) // Prerequisite is satisfied
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301173")
}

func TestValidatePrerequisites_MixedPrereqCoreq_BothViolations(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course with both prerequisites and corequisites
	course2301200 := &model.Course{
		Code: "2301200",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}}, // Must be taken before
				},
			},
		},
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301173"}}, // Must be taken same term
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301200", mock.Anything, 2023).Return(course2301200, nil)

	// Missing prerequisite AND corequisite in wrong term
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301173", Year: 2023, Semester: 1, Credits: 3}, // Corequisite (wrong term)
			{CourseCode: "2301200", Year: 2023, Semester: 2, Credits: 4}, // Missing prerequisite 2301170
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301200", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")
	assert.True(t, violations[0].TakenInWrongTerm)
	assert.Contains(t, violations[0].CoreqsTakenInWrongTerm, "2301173")
}

func TestValidatePrerequisites_TransitiveChain_Complete(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Set up transitive chain: 2301365 -> 2301263 -> 2301260
	course2301260 := &model.Course{
		Code:               "2301260",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	course2301263 := &model.Course{
		Code: "2301263",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301260"}},
				},
			},
		},
	}

	course2301365 := &model.Course{
		Code: "2301365",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301263"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", mock.Anything, 2023).Return(course2301260, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301263", mock.Anything, 2023).Return(course2301263, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301365", mock.Anything, 2023).Return(course2301365, nil)

	// Taking all courses in correct order
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301260", Year: 2023, Semester: 1, Credits: 4},
			{CourseCode: "2301263", Year: 2023, Semester: 2, Credits: 4},
			{CourseCode: "2301365", Year: 2024, Semester: 1, Credits: 4},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}

// ==================== GRADUATION ELIGIBILITY TESTS (via CheckGraduation) ====================

func TestCheckGraduation_InsufficientTotalCredits(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 120,
	}

	// Simple courses with no prerequisites
	course2301170 := &model.Course{Code: "2301170", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}

	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			MinCredits: 30,
			Courses:    []model.Course{{Code: "2301170", Credits: 30}},
		},
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 25}, // Only 25 credits, need 30
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.False(t, result.CanGraduate)
	assert.Equal(t, 25, result.TotalCredits)
	assert.Equal(t, 120, result.RequiredCredits)
}

func TestCheckGraduation_CategoryNotSatisfied(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 36,
	}

	// Simple courses with no prerequisites
	course2301170 := &model.Course{Code: "2301170", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}
	course2301173 := &model.Course{Code: "2301173", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}

	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			MinCredits: 30,
			Courses:    []model.Course{{Code: "2301170", Credits: 30}},
		},
		{
			NameTH:     "วิชาเฉพาะ",
			MinCredits: 6,
			Courses:    []model.Course{{Code: "2301173", Credits: 6}},
		},
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(course2301173, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 30}, // General ed satisfied
			{CourseCode: "2301173", Year: 2023, Semester: 2, Credits: 4},  // Major not satisfied (4 < 6)
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.False(t, result.CanGraduate)
	// Check if there are missing courses
	if len(result.MissingCourses) > 0 {
		assert.Contains(t, result.MissingCourses[0], "วิชาเฉพาะ")
		assert.Contains(t, result.MissingCourses[0], "2") // Missing 2 credits
	}
}

func TestCheckGraduation_HasPrerequisiteViolations(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 30,
	}

	// 2301180 requires 2301170 OR 2301173
	course2301180 := &model.Course{
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
	}

	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			MinCredits: 30,
			Courses:    []model.Course{{Code: "2301180", Credits: 30}},
		},
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(&model.Course{Code: "2301173"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301180", mock.Anything, 2023).Return(course2301180, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301180", Year: 2023, Semester: 1, Credits: 30}, // Missing prerequisites
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.False(t, result.CanGraduate)
	assert.Len(t, result.PrerequisiteViolations, 1)
}

func TestCheckGraduation_HasCreditLimitViolations(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 30,
	}

	// Simple courses with no prerequisites
	course2301170 := &model.Course{Code: "2301170", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}

	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			MinCredits: 30,
			Courses:    []model.Course{{Code: "2301170", Credits: 30}},
		},
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 30}, // 30 credits in one semester (exceeds 22 limit)
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.False(t, result.CanGraduate)
	assert.Len(t, result.CreditLimitViolations, 1)
}

func TestCheckGraduation_MultipleCategoriesMissing(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 40,
	}

	// Simple courses with no prerequisites
	course2301170 := &model.Course{Code: "2301170", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}
	course2301173 := &model.Course{Code: "2301173", PrerequisiteGroups: []model.PrerequisiteGroup{}, CorequisiteGroups: []model.PrerequisiteGroup{}}

	categories := []model.Category{
		{
			NameTH:     "วิชาศึกษาทั่วไป",
			MinCredits: 30,
			Courses:    []model.Course{{Code: "2301170", Credits: 30}},
		},
		{
			NameTH:     "วิชาเฉพาะ",
			MinCredits: 10,
			Courses:    []model.Course{{Code: "2301173", Credits: 10}},
		},
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2023).Return(course2301170, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2023).Return(course2301173, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2023, Semester: 1, Credits: 25}, // Missing 5 credits
			{CourseCode: "2301173", Year: 2023, Semester: 2, Credits: 5},  // Missing 5 credits
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)

	assert.NoError(t, err)
	assert.False(t, result.CanGraduate)
	// At least the total credits should be insufficient
	assert.Equal(t, 30, result.TotalCredits)
	assert.Equal(t, 40, result.RequiredCredits)
}

// ==================== GUIDE SPECIFIC CASES ====================

func TestGuide_Case1_BasicPrerequisite_2301173(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Case 1: 2301173 has no prerequisites
	course2301173 := &model.Course{
		Code:               "2301173",
		PrerequisiteGroups: []model.PrerequisiteGroup{},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2566).Return(course2301173, nil)

	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301173", Year: 2024, Semester: 1, Credits: 4},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)
	assert.NoError(t, err)
	assert.Empty(t, violations)
}

func TestGuide_Case2_8_ComplexPrerequisite_2301260(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// 2301260 requires (2301170 AND 2301172) OR 2301173
	// Note: In real DB/service, (A AND B) OR C is stored as two groups: (A OR C) AND (B OR C)
	course2301260_Internal := &model.Course{
		Code: "2301260",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301170"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301172"}},
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", mock.Anything, 2566).Return(course2301260_Internal, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", mock.Anything, 2566).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", mock.Anything, 2566).Return(&model.Course{Code: "2301172"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", mock.Anything, 2566).Return(&model.Course{Code: "2301173"}, nil)

	// Sub-case 2.1: Missing all
	progressFail := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301260", Year: 2024, Semester: 1, Credits: 4},
		},
	}
	violations, _ := (*graduationService).ValidatePrerequisites(progressFail)
	assert.Len(t, violations, 1)
	assert.Contains(t, violations[0].MissingPrereqs, "2301170")
	assert.Contains(t, violations[0].MissingPrereqs, "2301173")
	assert.Contains(t, violations[0].MissingPrereqs, "2301172")

	// Sub-case 2.2: Satisfied by 2301173
	progressOK1 := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301173", Year: 2024, Semester: 1, Credits: 4},
			{CourseCode: "2301260", Year: 2024, Semester: 2, Credits: 4},
		},
	}
	violations, _ = (*graduationService).ValidatePrerequisites(progressOK1)
	assert.Empty(t, violations)

	// Sub-case 2.3: Satisfied by 2301170 AND 2301172
	progressOK2 := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301170", Year: 2024, Semester: 1, Credits: 3},
			{CourseCode: "2301172", Year: 2024, Semester: 1, Credits: 1},
			{CourseCode: "2301260", Year: 2024, Semester: 2, Credits: 4},
		},
	}
	violations, _ = (*graduationService).ValidatePrerequisites(progressOK2)
	assert.Empty(t, violations)
}

func TestGuide_Case5_Corequisite_2301362(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// 2301362 requires 2301279 OR 2301369 as corequisite
	course2301362 := &model.Course{
		Code: "2301362",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: true,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301279"}},
					{PrerequisiteCourse: model.Course{Code: "2301369"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301362", mock.Anything, 2566).Return(course2301362, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301279", mock.Anything, 2566).Return(&model.Course{Code: "2301279"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301369", mock.Anything, 2566).Return(&model.Course{Code: "2301369"}, nil)

	// Test: Different term violation
	progressWrongTerm := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301279", Year: 2024, Semester: 1, Credits: 3},
			{CourseCode: "2301362", Year: 2024, Semester: 2, Credits: 3},
		},
	}
	violations, _ := (*graduationService).ValidatePrerequisites(progressWrongTerm)
	assert.Len(t, violations, 1)
	assert.True(t, violations[0].TakenInWrongTerm)

	// Test: Valid same term
	progressOK := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301279", Year: 2024, Semester: 1, Credits: 3},
			{CourseCode: "2301362", Year: 2024, Semester: 1, Credits: 3},
		},
	}
	violations, _ = (*graduationService).ValidatePrerequisites(progressOK)
	assert.Empty(t, violations)
}

func TestGuide_Case15_CompleteGraduationCheck(t *testing.T) {
	graduationService, mockCurriculum, mockCourse, mockCategory := setupGraduationService()

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 136,
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)

	categories := []model.Category{
		{NameTH: "วิชาแกน", MinCredits: 14},
		{NameTH: "วิชาเฉพาะด้าน", MinCredits: 39},
		{NameTH: "วิชาพื้นฐานวิทยาศาสตร์", MinCredits: 12},
		{NameTH: "กลุ่มวิชาภาษา", MinCredits: 12},
	}
	mockCategory.On("GetByCurriculumID", curriculumID).Return(categories, nil)

	// Mock courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", mock.Anything, mock.Anything, 2566).Return(&model.Course{
		PrerequisiteGroups: []model.PrerequisiteGroup{},
		CorequisiteGroups:  []model.PrerequisiteGroup{},
	}, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2566,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301173", Year: 2024, Semester: 1, Credits: 15, Grade: "A"},  // Total 15
			{CourseCode: "2301260", Year: 2024, Semester: 2, Credits: 15, Grade: "B+"}, // Total 30
			{CourseCode: "2301263", Year: 2025, Semester: 1, Credits: 22, Grade: "A"},  // Total 52
			{CourseCode: "2301263", Year: 2025, Semester: 2, Credits: 22, Grade: "B"},  // Total 74
			{CourseCode: "2301263", Year: 2026, Semester: 1, Credits: 22, Grade: "A"},  // Total 96
			{CourseCode: "2301263", Year: 2026, Semester: 2, Credits: 22, Grade: "A"},  // Total 118
			{CourseCode: "2301263", Year: 2027, Semester: 1, Credits: 18, Grade: "B+"}, // Total 136
		},
		ManualCredits: map[string]int{
			"วิชาพื้นฐานวิทยาศาสตร์": 12,
			"วิชาแกน":       14,
			"กลุ่มวิชาภาษา": 12,
			"วิชาเฉพาะด้าน": 39, // Ensure categories are satisfied for this complex check
		},
	}

	result, err := (*graduationService).CheckGraduation(progress)
	assert.NoError(t, err)
	assert.True(t, result.CanGraduate)
	assert.Equal(t, 136, result.TotalCredits)
	assert.Greater(t, result.GPAX, 3.0)
}
