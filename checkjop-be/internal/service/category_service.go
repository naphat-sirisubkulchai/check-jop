package service

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/repository"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CategoryService interface {
	Create(category *model.Category) error
	GetByID(id uuid.UUID) (*model.Category, error)
	GetByName(name string) (*model.Category, error)
	GetByCurriculumID(curriculumID uuid.UUID) ([]model.Category, error)
	GetAll() ([]model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
	ImportFromCSV(reader io.Reader) error
}

type categoryService struct {
	categoryRepo   repository.CategoryRepository
	curriculumRepo repository.CurriculumRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, curriculumRepo repository.CurriculumRepository) CategoryService {
	return &categoryService{
		categoryRepo:   categoryRepo,
		curriculumRepo: curriculumRepo,
	}
}

func (s *categoryService) Create(category *model.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *categoryService) GetByID(id uuid.UUID) (*model.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetByName(name string) (*model.Category, error) {
	return s.categoryRepo.GetByName(name)
}

func (s *categoryService) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Category, error) {
	return s.categoryRepo.GetByCurriculumID(curriculumID)
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) Update(category *model.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *categoryService) Delete(id uuid.UUID) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) ImportFromCSV(reader io.Reader) error {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV must have header and at least one data row")
	}

	// CSV format: categoryNameEn,categoryNameTH,curriculumName,minCredit,Year
	if len(records[0]) < 5 {
		return fmt.Errorf("CSV must have at least 5 columns")
	}

	// Build curriculum cache
	curriculumCache := make(map[string]*model.Curriculum)

	var categories []model.Category
	// Use a map to deduplicate categories based on curriculum_id + name_th + name_en
	categoryMap := make(map[string]model.Category)

	for i, record := range records[1:] {
		if len(record) < 5 {
			return fmt.Errorf("row %d: insufficient columns", i+2)
		}

		// Parse multiple curricula (comma-separated)
		curriculumNames := strings.Split(strings.TrimSpace(record[2]), ",")

		minCredits, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			return fmt.Errorf("row %d: invalid min_credits: %w", i+2, err)
		}

		// Create category for each curriculum
		for _, currName := range curriculumNames {
			currName = strings.TrimSpace(currName)
			if currName == "" {
				continue
			}

			// Get or cache curriculum
			curriculum, exists := curriculumCache[currName]
			if !exists {
				curr, err := s.curriculumRepo.GetByName(currName)
				if err != nil {
					return fmt.Errorf("row %d: curriculum '%s' not found: %w", i+2, currName, err)
				}
				curriculum = curr
				curriculumCache[currName] = curriculum
			}

			category := model.Category{
				CurriculumID: curriculum.ID,
				NameEN:       strings.TrimSpace(record[0]),
				NameTH:       strings.TrimSpace(record[1]),
				MinCredits:   minCredits,
			}

			// Create unique key for deduplication
			key := fmt.Sprintf("%s_%s_%s", curriculum.ID.String(), category.NameTH, category.NameEN)
			categoryMap[key] = category
		}
	}

	// Convert map to slice
	for _, category := range categoryMap {
		categories = append(categories, category)
	}

	return s.categoryRepo.BulkUpsert(categories)
}
