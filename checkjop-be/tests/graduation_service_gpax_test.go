package tests

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateGPAX_AllGraded(t *testing.T) {
	// graduationService, _, _, _ := setupGraduationService()

	// Create a service instance to access private method via reflection or just test public behavior
	// Since calculateGPAX is private, we'll test it through CheckGraduation or similar,
	// OR we can export it for testing, OR we can just rely on CheckGraduation result.
	// Let's rely on CheckGraduation result.

	mockCurriculum := &MockCurriculumRepository{}
	mockCourse := &MockCourseRepository{}
	mockCategory := &MockCategoryRepository{}

	// Re-setup with mocks
	svc := service.NewGraduationService(mockCurriculum, mockCourse, mockCategory)

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{
		ID:              curriculumID,
		MinTotalCredits: 120,
	}

	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return([]model.Category{}, nil)

	// Mock GetByCodeAndYear for courses to avoid panic in ValidatePrerequisites
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "A", mock.Anything, 2023).Return(&model.Course{Code: "A"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "B", mock.Anything, 2023).Return(&model.Course{Code: "B"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "C", mock.Anything, 2023).Return(&model.Course{Code: "C"}, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "A", Credits: 3, Grade: "A"},  // 4.0 * 3 = 12
			{CourseCode: "B", Credits: 3, Grade: "B"},  // 3.0 * 3 = 9
			{CourseCode: "C", Credits: 2, Grade: "C+"}, // 2.5 * 2 = 5
		},
	}
	// Total Points: 12 + 9 + 5 = 26
	// Total Credits: 3 + 3 + 2 = 8
	// GPAX: 26 / 8 = 3.25

	result, err := svc.CheckGraduation(progress)
	assert.NoError(t, err)
	assert.Equal(t, 3.25, result.GPAX)
}

func TestCalculateGPAX_WithNonGraded(t *testing.T) {
	mockCurriculum := &MockCurriculumRepository{}
	mockCourse := &MockCourseRepository{}
	mockCategory := &MockCategoryRepository{}
	svc := service.NewGraduationService(mockCurriculum, mockCourse, mockCategory)

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{ID: curriculumID}
	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return([]model.Category{}, nil)

	// Mock GetByCodeAndYear for courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "A", mock.Anything, 2023).Return(&model.Course{Code: "A"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "B", mock.Anything, 2023).Return(&model.Course{Code: "B"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "C", mock.Anything, 2023).Return(&model.Course{Code: "C"}, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "A", Credits: 3, Grade: "A"}, // 4.0 * 3 = 12
			{CourseCode: "B", Credits: 3, Grade: "W"}, // Ignored
			{CourseCode: "C", Credits: 3, Grade: "S"}, // Ignored
		},
	}
	// Total Points: 12
	// Total Credits: 3
	// GPAX: 4.0

	result, err := svc.CheckGraduation(progress)
	assert.NoError(t, err)
	assert.Equal(t, 4.0, result.GPAX)
}

func TestCalculateGPAX_WithF(t *testing.T) {
	mockCurriculum := &MockCurriculumRepository{}
	mockCourse := &MockCourseRepository{}
	mockCategory := &MockCategoryRepository{}
	svc := service.NewGraduationService(mockCurriculum, mockCourse, mockCategory)

	curriculumID := uuid.New()
	curriculum := &model.Curriculum{ID: curriculumID}
	mockCurriculum.On("GetByID", curriculumID).Return(curriculum, nil)
	mockCategory.On("GetByCurriculumID", curriculumID).Return([]model.Category{}, nil)

	// Mock GetByCodeAndYear for courses
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "A", mock.Anything, 2023).Return(&model.Course{Code: "A"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "B", mock.Anything, 2023).Return(&model.Course{Code: "B"}, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "A", Credits: 3, Grade: "A"}, // 4.0 * 3 = 12
			{CourseCode: "B", Credits: 3, Grade: "F"}, // 0.0 * 3 = 0
		},
	}
	// Total Points: 12
	// Total Credits: 6
	// GPAX: 2.0

	result, err := svc.CheckGraduation(progress)
	assert.NoError(t, err)
	assert.Equal(t, 2.0, result.GPAX)
}

func TestValidatePrerequisites_FGrade_Violation(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "PRE"
	courseMain := &model.Course{
		Code: "MAIN",
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "PRE"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "PRE", mock.Anything, 2023).Return(&model.Course{Code: "PRE"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "MAIN", mock.Anything, 2023).Return(courseMain, nil)

	// Taking MAIN with PRE having grade F
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "PRE", Year: 2023, Semester: 1, Credits: 3, Grade: "F"},
			{CourseCode: "MAIN", Year: 2023, Semester: 2, Credits: 3, Grade: "A"},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "MAIN", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "PRE")
}

func TestValidatePrerequisites_Corequisite_FGrade_Allowed(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Course requiring "CO" as corequisite
	courseMain := &model.Course{
		Code: "MAIN",
		CorequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "CO"}},
				},
			},
		},
	}

	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "CO", mock.Anything, 2023).Return(&model.Course{Code: "CO"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "MAIN", mock.Anything, 2023).Return(courseMain, nil)

	// Taking MAIN with CO having grade F (should be allowed for corequisite)
	progress := &model.StudentProgress{
		CurriculumID:  uuid.New(),
		AdmissionYear: 2023,
		Courses: []model.CompletedCourse{
			{CourseCode: "CO", Year: 2023, Semester: 1, Credits: 3, Grade: "F"},
			{CourseCode: "MAIN", Year: 2023, Semester: 1, Credits: 3, Grade: "A"},
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Empty(t, violations)
}
