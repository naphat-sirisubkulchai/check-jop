package tests

import (
	"checkjop-be/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCurriculumRepository provides a mock implementation of CurriculumRepository
type MockCurriculumRepository struct {
	mock.Mock
}

func (m *MockCurriculumRepository) Create(curriculum *model.Curriculum) error {
	args := m.Called(curriculum)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetByID(id uuid.UUID) (*model.Curriculum, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Curriculum), args.Error(1)
}

func (m *MockCurriculumRepository) GetByName(name string) (*model.Curriculum, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Curriculum), args.Error(1)
}

func (m *MockCurriculumRepository) GetAll() ([]model.Curriculum, error) {
	args := m.Called()
	return args.Get(0).([]model.Curriculum), args.Error(1)
}

func (m *MockCurriculumRepository) GetActiveByYear(year int) ([]model.Curriculum, error) {
	args := m.Called(year)
	return args.Get(0).([]model.Curriculum), args.Error(1)
}

func (m *MockCurriculumRepository) Update(curriculum *model.Curriculum) error {
	args := m.Called(curriculum)
	return args.Error(0)
}

func (m *MockCurriculumRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCurriculumRepository) Upsert(curriculum *model.Curriculum) error {
	args := m.Called(curriculum)
	return args.Error(0)
}

func (m *MockCurriculumRepository) BulkUpsert(curriculums []model.Curriculum) error {
	args := m.Called(curriculums)
	return args.Error(0)
}

func (m *MockCurriculumRepository) GetAllWithOutCatAndCourse() ([]model.Curriculum, error) {
	args := m.Called()
	return args.Get(0).([]model.Curriculum), args.Error(1)
}

func (m *MockCurriculumRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

// MockCourseRepository provides a mock implementation of CourseRepository
type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) Create(course *model.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) GetByID(id uuid.UUID) (*model.Course, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByCodeAndYear(code string, year int) (*model.Course, error) {
	args := m.Called(code, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByName(name string) (*model.Course, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Course, error) {
	args := m.Called(curriculumID)
	return args.Get(0).([]model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByCategoryID(categoryID uuid.UUID) ([]model.Course, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]model.Course), args.Error(1)
}

func (m *MockCourseRepository) GetAll() ([]model.Course, error) {
	args := m.Called()
	return args.Get(0).([]model.Course), args.Error(1)
}

func (m *MockCourseRepository) Update(course *model.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCourseRepository) CreateFromCSV(courses []model.Course) error {
	args := m.Called(courses)
	return args.Error(0)
}

func (m *MockCourseRepository) Upsert(course *model.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) BulkUpsert(courses []model.Course) error {
	args := m.Called(courses)
	return args.Error(0)
}

func (m *MockCourseRepository) SetPrerequisites(courseID uuid.UUID, prerequisiteIDs []uuid.UUID) error {
	args := m.Called(courseID, prerequisiteIDs)
	return args.Error(0)
}

func (m *MockCourseRepository) SetCorequisites(courseID uuid.UUID, corequisiteIDs []uuid.UUID) error {
	args := m.Called(courseID, corequisiteIDs)
	return args.Error(0)
}

func (m *MockCourseRepository) GetByCodeAndCurriculumIDAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error) {
	args := m.Called(code, curriculumID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Course), args.Error(1)
}

func (m *MockCourseRepository) ExistsByCodeAndCurriculumID(code string, curriculumID uuid.UUID) bool {
	args := m.Called(code, curriculumID)
	return args.Bool(0)
}

func (m *MockCourseRepository) CatalogYearExists(curriculumID uuid.UUID, year int) bool {
	args := m.Called(curriculumID, year)
	return args.Bool(0)
}

func (m *MockCourseRepository) GetLatestAvailableCatalogYear(curriculumID uuid.UUID, maxYear int) (int, bool) {
	args := m.Called(curriculumID, maxYear)
	return args.Int(0), args.Bool(1)
}

func (m *MockCourseRepository) CourseHasCFOptionInAnyCatalogYear(code string, curriculumID uuid.UUID) bool {
	args := m.Called(code, curriculumID)
	return args.Bool(0)
}

func (m *MockCourseRepository) SetPrerequisiteGroups(courseID uuid.UUID, groups []model.PrerequisiteGroup) error {
	args := m.Called(courseID, groups)
	return args.Error(0)
}

func (m *MockCourseRepository) SetCorequisiteGroups(courseID uuid.UUID, groups []model.PrerequisiteGroup) error {
	args := m.Called(courseID, groups)
	return args.Error(0)
}

func (m *MockCourseRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCourseRepository) DeleteByYear(year int) error {
	args := m.Called(year)
	return args.Error(0)
}

// MockCategoryRepository provides a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(id uuid.UUID) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByName(name string) (*model.Category, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Category, error) {
	args := m.Called(curriculumID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetAll() ([]model.Category, error) {
	args := m.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) Upsert(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) DeleteAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCategoryRepository) BulkUpsert(categories []model.Category) error {
	args := m.Called(categories)
	return args.Error(0)
}
