package repository

import (
	"checkjop-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SetDefaultRepository interface {
	Create(setDefault *model.SetDefault) error
	GetByID(id uuid.UUID) (*model.SetDefault, error)
	GetByCurriculumName(curriculumName string) ([]model.SetDefault, error)
	GetByCurriculumID(curriculumID uuid.UUID) ([]model.SetDefault, error)
	GetAll() ([]model.SetDefault, error)
	Update(setDefault *model.SetDefault) error
	Delete(id uuid.UUID) error
	BulkUpsert(setDefaults []model.SetDefault) error
	DeleteByCurriculumID(curriculumID uuid.UUID) error
}

type setDefaultRepository struct {
	db *gorm.DB
}

func NewSetDefaultRepository(db *gorm.DB) SetDefaultRepository {
	return &setDefaultRepository{db}
}

func (r *setDefaultRepository) Create(setDefault *model.SetDefault) error {
	return r.db.Create(setDefault).Error
}

func (r *setDefaultRepository) GetByID(id uuid.UUID) (*model.SetDefault, error) {
	var setDefault model.SetDefault
	err := r.db.Preload("Curriculum").Preload("Course").
		First(&setDefault, id).Error
	if err != nil {
		return nil, err
	}
	return &setDefault, nil
}

func (r *setDefaultRepository) GetByCurriculumName(curriculumName string) ([]model.SetDefault, error) {
	var setDefaults []model.SetDefault
	err := r.db.Joins("JOIN curriculums ON set_defaults.curriculum_id = curriculums.id").
		Where("curriculums.name_th = ? OR curriculums.name_en = ?", curriculumName, curriculumName).
		Preload("Curriculum").Preload("Course").
		Find(&setDefaults).Error
	if err != nil {
		return nil, err
	}
	return setDefaults, nil
}

func (r *setDefaultRepository) GetByCurriculumID(curriculumID uuid.UUID) ([]model.SetDefault, error) {
	var setDefaults []model.SetDefault
	err := r.db.Where("curriculum_id = ?", curriculumID).
		Preload("Curriculum").Preload("Course").
		Find(&setDefaults).Error
	if err != nil {
		return nil, err
	}
	return setDefaults, nil
}

func (r *setDefaultRepository) GetAll() ([]model.SetDefault, error) {
	var setDefaults []model.SetDefault
	err := r.db.Preload("Curriculum").Preload("Course").
		Find(&setDefaults).Error
	if err != nil {
		return nil, err
	}
	return setDefaults, nil
}

func (r *setDefaultRepository) Update(setDefault *model.SetDefault) error {
	return r.db.Save(setDefault).Error
}

func (r *setDefaultRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.SetDefault{}, id).Error
}

func (r *setDefaultRepository) BulkUpsert(setDefaults []model.SetDefault) error {
	if len(setDefaults) == 0 {
		return nil
	}

	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing setDefaults with same curriculum_id and course_id combinations
	for _, setDefault := range setDefaults {
		err := tx.Where("curriculum_id = ? AND course_id = ?", setDefault.CurriculumID, setDefault.CourseID).
			Delete(&model.SetDefault{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert all new setDefaults
	err := tx.Create(&setDefaults).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *setDefaultRepository) DeleteByCurriculumID(curriculumID uuid.UUID) error {
	return r.db.Where("curriculum_id = ?", curriculumID).Delete(&model.SetDefault{}).Error
}
