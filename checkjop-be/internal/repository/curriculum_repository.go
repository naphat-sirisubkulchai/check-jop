package repository

import (
	"checkjop-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CurriculumRepository interface {
	Create(curriculum *model.Curriculum) error
	GetByID(id uuid.UUID) (*model.Curriculum, error)
	GetByName(name string) (*model.Curriculum, error)
	GetAll() ([]model.Curriculum, error)
	GetActiveByYear(year int) ([]model.Curriculum, error)
	Update(curriculum *model.Curriculum) error
	Delete(id uuid.UUID) error
	Upsert(curriculum *model.Curriculum) error
	BulkUpsert(curriculums []model.Curriculum) error
	GetAllWithOutCatAndCourse() ([]model.Curriculum, error)
	DeleteAll() error
}

type curriculumRepository struct {
	db *gorm.DB
}

func NewCurriculumRepository(db *gorm.DB) CurriculumRepository {
	return &curriculumRepository{db}
}

func (r *curriculumRepository) Create(curriculum *model.Curriculum) error {
	return r.db.Create(curriculum).Error
}

func (r *curriculumRepository) GetByID(id uuid.UUID) (*model.Curriculum, error) {
	var curriculum model.Curriculum
	err := r.db.Preload("Categories").Preload("Courses").
		First(&curriculum, id).Error
	if err != nil {
		return nil, err
	}
	return &curriculum, err
}

func (r *curriculumRepository) GetByName(name string) (*model.Curriculum, error) {
	var curriculum model.Curriculum
	err := r.db.Where("name_th = ? OR name_en = ?", name, name).
		Preload("Categories").Preload("Courses").
		First(&curriculum).Error
	if err != nil {
		return nil, err
	}
	return &curriculum, err
}

func (r *curriculumRepository) GetAll() ([]model.Curriculum, error) {
	var curriculums []model.Curriculum
	err := r.db.Preload("Categories").
		Preload("Courses").Preload("Courses.PrerequisiteGroups").Preload("Courses.CorequisiteGroups").
		Find(&curriculums).Error
	if err != nil {
		return nil, err
	}

	return curriculums, err
}
func (r *curriculumRepository) GetAllWithOutCatAndCourse() ([]model.Curriculum, error) {
	var curriculums []model.Curriculum
	err := r.db.Find(&curriculums).Error
	if err != nil {
		return nil, err
	}
	return curriculums, err
}

func (r *curriculumRepository) GetActiveByYear(year int) ([]model.Curriculum, error) {
	var curricula []model.Curriculum
	err := r.db.Where("year = ? AND is_active = ?", year, true).Find(&curricula).Error
	return curricula, err
}

func (r *curriculumRepository) Update(curriculum *model.Curriculum) error {
	return r.db.Save(curriculum).Error
}

func (r *curriculumRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Curriculum{}, id).Error
}

func (r *curriculumRepository) Upsert(curriculum *model.Curriculum) error {
	return r.db.Save(curriculum).Error
}

func (r *curriculumRepository) BulkUpsert(curriculums []model.Curriculum) error {
	if len(curriculums) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name_th"}, {Name: "name_en"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"year", "min_total_credits", "is_active", "updated_at", "deleted_at",
		}),
	}).Create(&curriculums).Error
}

func (r *curriculumRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Curriculum{}).Error
}
