package repository

import (
	"checkjop-be/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id uuid.UUID) (*model.Category, error)
	GetByName(name string) (*model.Category, error)
	GetByCurriculumID(curriculumID uuid.UUID) ([]model.Category, error)
	GetAll() ([]model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
	Upsert(category *model.Category) error
	BulkUpsert(categories []model.Category) error
	DeleteAll() error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	err := r.db.Preload("Courses").First(&category, id).Error
	return &category, err
}

func (r *categoryRepository) GetByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name_th = ? OR name_en = ?", name, name).Preload("Courses").First(&category).Error
	return &category, err
}

func (r *categoryRepository) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("curriculum_id = ?", curriculumID).Preload("Courses").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Preload("Courses").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *categoryRepository) Upsert(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) BulkUpsert(categories []model.Category) error {
	if len(categories) == 0 {
		return nil
	}

	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing categories with same name_th, name_en, and curriculum_id combinations
	for _, category := range categories {
		err := tx.Unscoped().Where("name_th = ? AND name_en = ? AND curriculum_id = ?",
			category.NameTH, category.NameEN, category.CurriculumID).
			Delete(&model.Category{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert all new categories
	err := tx.Create(&categories).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *categoryRepository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Category{}).Error
}
