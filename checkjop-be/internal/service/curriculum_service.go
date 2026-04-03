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

type CurriculumService interface {
	Create(curriculum *model.Curriculum) error
	GetByID(id uuid.UUID) (*model.Curriculum, error)
	GetByName(name string) (*model.Curriculum, error)
	GetAll() ([]model.Curriculum, error)
	GetActiveByYear(year int) ([]model.Curriculum, error)
	Update(curriculum *model.Curriculum) error
	Delete(id uuid.UUID) error
	ImportFromCSV(reader io.Reader) error
	GetAllWithOutCatAndCourse() ([]model.Curriculum, error)
}

type curriculumService struct {
	curriculumRepo repository.CurriculumRepository
}

func NewCurriculumService(curriculumRepo repository.CurriculumRepository) CurriculumService {
	return &curriculumService{
		curriculumRepo: curriculumRepo,
	}
}

func (s *curriculumService) Create(curriculum *model.Curriculum) error {
	return s.curriculumRepo.Create(curriculum)
}

func (s *curriculumService) GetByID(id uuid.UUID) (*model.Curriculum, error) {
	return s.curriculumRepo.GetByID(id)
}

func (s *curriculumService) GetByName(name string) (*model.Curriculum, error) {
	return s.curriculumRepo.GetByName(name)
}

func (s *curriculumService) GetAll() ([]model.Curriculum, error) {
	return s.curriculumRepo.GetAll()
}
func (s *curriculumService) GetAllWithOutCatAndCourse() ([]model.Curriculum, error) {
	return s.curriculumRepo.GetAllWithOutCatAndCourse()
}

func (s *curriculumService) GetActiveByYear(year int) ([]model.Curriculum, error) {
	return s.curriculumRepo.GetActiveByYear(year)
}

func (s *curriculumService) Update(curriculum *model.Curriculum) error {
	return s.curriculumRepo.Update(curriculum)
}

func (s *curriculumService) Delete(id uuid.UUID) error {
	return s.curriculumRepo.Delete(id)
}

func (s *curriculumService) ImportFromCSV(reader io.Reader) error {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV must have at least header and one data row")
	}

	// CSV format: curriculumNameEN,curriculumNameTH,Year,minTotalCredit,isActive
	var curriculums []model.Curriculum
	for i, record := range records[1:] { // Skip header
		if len(record) < 5 {
			return fmt.Errorf("row %d: insufficient columns, expected 5 got %d", i+2, len(record))
		}

		year, err := strconv.Atoi(strings.TrimSpace(record[2]))
		if err != nil {
			return fmt.Errorf("row %d: invalid year '%s': %w", i+2, record[2], err)
		}

		minCredits, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			return fmt.Errorf("row %d: invalid min_total_credits '%s': %w", i+2, record[3], err)
		}

		isActive := strings.ToUpper(strings.TrimSpace(record[4])) == "TRUE"

		curriculum := model.Curriculum{
			NameEN:          strings.TrimSpace(record[0]),
			NameTH:          strings.TrimSpace(record[1]),
			Year:            year,
			MinTotalCredits: minCredits,
			IsActive:        isActive,
		}

		curriculums = append(curriculums, curriculum)
	}

	return s.curriculumRepo.BulkUpsert(curriculums)
}
