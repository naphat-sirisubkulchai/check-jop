package tests

import (
	"checkjop-be/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCrossCurriculum_PrerequisiteCheck_NoEnforcement(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	// Curriculum A and B
	// curriculumID_A := uuid.New()
	curriculumID_B := uuid.New()

	// Course 2301172 (from guide Case 9) is in Curriculum A but NOT in B
	// 2301172 has a corequisite rule in its original curriculum (Curr A)

	// Mock GetByCodeAndCurriculumIDAndYear:
	// When checking for student in Curr B, searching for 2301170 and 2301172
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301170", curriculumID_B, 2024).Return(&model.Course{Code: "2301170"}, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301172", curriculumID_B, 2024).Return((*model.Course)(nil), assert.AnError)

	// Student is in Curriculum B
	progress := &model.StudentProgress{
		CurriculumID:  curriculumID_B,
		AdmissionYear: 2024,
		Courses: []model.CompletedCourse{
			{
				CourseCode: "2301170",
				Year:       2024,
				Semester:   1,
				Credits:    3,
			},
			{
				CourseCode: "2301172",
				Year:       2024,
				Semester:   2,
				Credits:    1,
			},
		},
	}

	// Validate Prerequisites
	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	// As per Case 9 in the guide: "วิชาที่ไม่อยู่ในหลักสูตรของนักศึกษา จะไม่มีการตรวจ prerequisites/corequisites (ถือว่าไม่มีกฎ)"
	// Since 2301172 is not found in Curr B, no corequisite violation should be reported for 2301172.
	assert.Empty(t, violations)
}

func TestCrossCurriculum_PrerequisiteEnforced_WhenInCurriculum(t *testing.T) {
	graduationService, _, mockCourse, _ := setupGraduationService()

	curriculumID := uuid.New()

	// Course 2301260 is IN the curriculum and has prerequisites
	course2301260 := &model.Course{
		Code: "2301260",
		Year: 2024,
		PrerequisiteGroups: []model.PrerequisiteGroup{
			{
				IsOrGroup: false,
				PrerequisiteCourses: []model.PrerequisiteCourseLink{
					{PrerequisiteCourse: model.Course{Code: "2301173"}},
				},
			},
		},
	}

	// Mock repository finding the course in the student's curriculum
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301260", curriculumID, 2024).Return(course2301260, nil)
	mockCourse.On("GetByCodeAndCurriculumIDAndYear", "2301173", curriculumID, 2024).Return(&model.Course{Code: "2301173"}, nil)

	progress := &model.StudentProgress{
		CurriculumID:  curriculumID,
		AdmissionYear: 2024,
		Courses: []model.CompletedCourse{
			{CourseCode: "2301260", Year: 2024, Semester: 1, Credits: 4}, // Missing internal prerequisite 2301173
		},
	}

	violations, err := (*graduationService).ValidatePrerequisites(progress)

	assert.NoError(t, err)
	assert.Len(t, violations, 1)
	assert.Equal(t, "2301260", violations[0].CourseCode)
	assert.Contains(t, violations[0].MissingPrereqs, "2301173")
}
