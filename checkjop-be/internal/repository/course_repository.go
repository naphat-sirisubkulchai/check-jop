package repository

import (
	"checkjop-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CourseRepository interface {
	Create(course *model.Course) error
	GetByID(id uuid.UUID) (*model.Course, error)
	GetByCodeAndYear(code string, year int) (*model.Course, error)
	GetByName(name string) (*model.Course, error)
	GetByCurriculumID(curriculumID uuid.UUID) ([]model.Course, error)
	GetByCategoryID(categoryID uuid.UUID) ([]model.Course, error)
	GetAll() ([]model.Course, error)
	Update(course *model.Course) error
	Delete(id uuid.UUID) error
	CreateFromCSV(courses []model.Course) error
	Upsert(course *model.Course) error
	BulkUpsert(courses []model.Course) error
	SetPrerequisites(courseID uuid.UUID, prerequisiteIDs []uuid.UUID) error
	SetCorequisites(courseID uuid.UUID, corequisiteIDs []uuid.UUID) error
	GetByCodeAndCurriculumIDAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error)
	ExistsByCodeAndCurriculumID(code string, curriculumID uuid.UUID) bool
	CourseHasCFOptionInAnyCatalogYear(code string, curriculumID uuid.UUID) bool
	CatalogYearExists(curriculumID uuid.UUID, year int) bool
	GetLatestAvailableCatalogYear(curriculumID uuid.UUID, maxYear int) (int, bool)
	SetPrerequisiteGroups(courseID uuid.UUID, groups []model.PrerequisiteGroup) error
	SetCorequisiteGroups(courseID uuid.UUID, groups []model.PrerequisiteGroup) error
	DeleteAll() error
	DeleteByYear(year int) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db}
}

func (r *courseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) GetByID(id uuid.UUID) (*model.Course, error) {
	var course model.Course
	err := r.db.Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		First(&course, id).Error
	r.initializeSlices(&course)
	return &course, err
}

func (r *courseRepository) GetByCodeAndYear(code string, year int) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("code = ? AND year = ?", code, year).
		Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		First(&course).Error
	r.initializeSlices(&course)
	return &course, err
}

func (r *courseRepository) GetByName(name string) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("name_en = ?", name).
		Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		First(&course).Error
	r.initializeSlices(&course)
	return &course, err
}

func (r *courseRepository) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Where("curriculum_id = ?", curriculumID).
		Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Find(&courses).Error
	r.initializeSlicesForCourses(courses)
	return courses, err
}

func (r *courseRepository) GetByCategoryID(categoryID uuid.UUID) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Where("category_id = ?", categoryID).
		Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Find(&courses).Error
	r.initializeSlicesForCourses(courses)
	return courses, err
}

func (r *courseRepository) GetAll() ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Find(&courses).Error
	r.initializeSlicesForCourses(courses)
	return courses, err
}

func (r *courseRepository) CreateFromCSV(courses []model.Course) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, course := range courses {
			if err := tx.Create(&course).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *courseRepository) Update(course *model.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Course{}, id).Error
}

func (r *courseRepository) Upsert(course *model.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) BulkUpsert(courses []model.Course) error {
	if len(courses) == 0 {
		return nil
	}

	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Use OnConflict to handle upserts efficiently and avoid unique constraint violations
	// This handles cases where DeleteAll might not have fully cleared data or if we want to update existing records
	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "year"}, {Name: "curriculum_id"}, {Name: "category_id"}},
		UpdateAll: true,
	}).Create(&courses).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// initializeSlices ensures that PrerequisiteGroups and CorequisiteGroups are empty arrays instead of nil
func (r *courseRepository) initializeSlices(course *model.Course) {
	if course.PrerequisiteGroups == nil {
		course.PrerequisiteGroups = []model.PrerequisiteGroup{}
	}
	if course.CorequisiteGroups == nil {
		course.CorequisiteGroups = []model.PrerequisiteGroup{}
	}
}

// initializeSlicesForCourses initializes empty slices for all courses in a slice
func (r *courseRepository) initializeSlicesForCourses(courses []model.Course) {
	for i := range courses {
		r.initializeSlices(&courses[i])
	}
}

// SetPrerequisites sets the prerequisite relationships for a course
func (r *courseRepository) SetPrerequisites(courseID uuid.UUID, prerequisiteIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Clear existing prerequisites
		if err := tx.Exec("DELETE FROM course_prerequisites WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}

		// Add new prerequisites
		for _, prereqID := range prerequisiteIDs {
			if err := tx.Exec("INSERT INTO course_prerequisites (course_id, prerequisite_id) VALUES (?, ?)", courseID, prereqID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SetCorequisites sets the corequisite relationships for a course
func (r *courseRepository) SetCorequisites(courseID uuid.UUID, corequisiteIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Clear existing corequisites
		if err := tx.Exec("DELETE FROM course_corequisites WHERE course_id = ?", courseID).Error; err != nil {
			return err
		}

		// Add new corequisites
		for _, coreqID := range corequisiteIDs {
			if err := tx.Exec("INSERT INTO course_corequisites (course_id, corequisite_id) VALUES (?, ?)", courseID, coreqID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CatalogYearExists returns true if there is at least one course in the DB for the given
// curriculum and year — i.e. the catalog for that academic year has been imported.
func (r *courseRepository) CatalogYearExists(curriculumID uuid.UUID, year int) bool {
	var count int64
	r.db.Model(&model.Course{}).Where("curriculum_id = ? AND year = ?", curriculumID, year).Count(&count)
	return count > 0
}

// GetLatestAvailableCatalogYear returns the most recent catalog year <= maxYear that has
// been imported for the given curriculum. Returns (year, true) if found, (0, false) otherwise.
func (r *courseRepository) GetLatestAvailableCatalogYear(curriculumID uuid.UUID, maxYear int) (int, bool) {
	var year int
	err := r.db.Model(&model.Course{}).
		Where("curriculum_id = ? AND year <= ?", curriculumID, maxYear).
		Select("MAX(year)").
		Scan(&year).Error
	if err != nil || year == 0 {
		return 0, false
	}
	return year, true
}

// ExistsByCodeAndCurriculumID returns true if a course with the given code exists in the
// curriculum under any admission year. Used to distinguish real course codes from
// custom manual strings (e.g. "FREE II (3)").
func (r *courseRepository) ExistsByCodeAndCurriculumID(code string, curriculumID uuid.UUID) bool {
	var count int64
	r.db.Model(&model.Course{}).Where("code = ? AND curriculum_id = ?", code, curriculumID).Count(&count)
	return count > 0
}

// CourseHasCFOptionInAnyCatalogYear returns true if the course has has_cf_option=true in any catalog year.
func (r *courseRepository) CourseHasCFOptionInAnyCatalogYear(code string, curriculumID uuid.UUID) bool {
	var count int64
	r.db.Model(&model.Course{}).Where("code = ? AND curriculum_id = ? AND has_cf_option = true", code, curriculumID).Count(&count)
	return count > 0
}

// GetByCodeAndCurriculumIDAndYear finds a course by code within a specific curriculum and year
func (r *courseRepository) GetByCodeAndCurriculumIDAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("code = ? AND curriculum_id = ? AND year = ?", code, curriculumID, year).
		Preload("Curriculum").Preload("Category").
		Preload("PrerequisiteGroups", "group_type = ?", "prerequisite").
		Preload("PrerequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		Preload("CorequisiteGroups", "group_type = ?", "corequisite").
		Preload("CorequisiteGroups.PrerequisiteCourses.PrerequisiteCourse").
		First(&course).Error
	if err != nil {
		return nil, err
	}
	r.initializeSlices(&course)
	return &course, nil
}

// SetPrerequisiteGroups sets prerequisite groups with OR/AND logic for a course
func (r *courseRepository) SetPrerequisiteGroups(courseID uuid.UUID, groups []model.PrerequisiteGroup) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, get existing group IDs
		var existingGroups []model.PrerequisiteGroup
		if err := tx.Select("id").Where("course_id = ? AND group_type = ?", courseID, "prerequisite").Find(&existingGroups).Error; err != nil {
			return err
		}

		// Delete prerequisite course links first (child records)
		for _, group := range existingGroups {
			if err := tx.Where("group_id = ?", group.ID).Delete(&model.PrerequisiteCourseLink{}).Error; err != nil {
				return err
			}
		}

		// Then delete prerequisite groups (parent records)
		if err := tx.Where("course_id = ? AND group_type = ?", courseID, "prerequisite").Delete(&model.PrerequisiteGroup{}).Error; err != nil {
			return err
		}

		// Insert new prerequisite groups
		for _, group := range groups {
			group.CourseID = courseID
			group.GroupType = "prerequisite"
			if err := tx.Create(&group).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SetCorequisiteGroups sets corequisite groups with OR/AND logic for a course
func (r *courseRepository) SetCorequisiteGroups(courseID uuid.UUID, groups []model.PrerequisiteGroup) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, get existing group IDs
		var existingGroups []model.PrerequisiteGroup
		if err := tx.Select("id").Where("course_id = ? AND group_type = ?", courseID, "corequisite").Find(&existingGroups).Error; err != nil {
			return err
		}

		// Delete prerequisite course links first (child records)
		for _, group := range existingGroups {
			if err := tx.Where("group_id = ?", group.ID).Delete(&model.PrerequisiteCourseLink{}).Error; err != nil {
				return err
			}
		}

		// Then delete corequisite groups (parent records)
		if err := tx.Where("course_id = ? AND group_type = ?", courseID, "corequisite").Delete(&model.PrerequisiteGroup{}).Error; err != nil {
			return err
		}

		// Insert new corequisite groups
		for _, group := range groups {
			group.CourseID = courseID
			group.GroupType = "corequisite"
			if err := tx.Create(&group).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteAll deletes all courses from the database (hard delete)
func (r *courseRepository) DeleteAll() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete all prerequisite course links first (child of groups, references courses)
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.PrerequisiteCourseLink{}).Error; err != nil {
			return err
		}

		// Delete all prerequisite/corequisite groups (child of courses)
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.PrerequisiteGroup{}).Error; err != nil {
			return err
		}

		// Delete all courses
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Course{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteByYear deletes all courses for a specific year from the database (hard delete)
func (r *courseRepository) DeleteByYear(year int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, get all course IDs for this year
		var courseIDs []uuid.UUID
		if err := tx.Model(&model.Course{}).Where("year = ?", year).Pluck("id", &courseIDs).Error; err != nil {
			return err
		}

		if len(courseIDs) == 0 {
			return nil // No courses to delete
		}

		// Get all prerequisite/corequisite group IDs for these courses
		var groupIDs []uuid.UUID
		if err := tx.Model(&model.PrerequisiteGroup{}).Where("course_id IN ?", courseIDs).Pluck("id", &groupIDs).Error; err != nil {
			return err
		}

		// Delete prerequisite course links first (child of groups)
		if len(groupIDs) > 0 {
			if err := tx.Unscoped().Where("group_id IN ?", groupIDs).Delete(&model.PrerequisiteCourseLink{}).Error; err != nil {
				return err
			}
		}

		// Delete prerequisite/corequisite groups (child of courses)
		if err := tx.Unscoped().Where("course_id IN ?", courseIDs).Delete(&model.PrerequisiteGroup{}).Error; err != nil {
			return err
		}

		// Delete courses for this year
		if err := tx.Unscoped().Where("year = ?", year).Delete(&model.Course{}).Error; err != nil {
			return err
		}

		return nil
	})
}
