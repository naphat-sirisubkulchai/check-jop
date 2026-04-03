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

type SetDefaultService interface {
	Create(setDefault *model.SetDefault) error
	GetByID(id uuid.UUID) (*model.SetDefault, error)
	GetByCurriculumName(curriculumName string) ([]model.SetDefault, error)
	GetByCurriculumID(curriculumID uuid.UUID) ([]model.SetDefault, error)
	GetAll() ([]model.SetDefault, error)
	Update(setDefault *model.SetDefault) error
	Delete(id uuid.UUID) error
	ImportFromCSV(reader io.Reader) error
}

type setDefaultService struct {
	setDefaultRepo repository.SetDefaultRepository
	curriculumRepo repository.CurriculumRepository
	courseRepo     repository.CourseRepository
}

func NewSetDefaultService(
	setDefaultRepo repository.SetDefaultRepository,
	curriculumRepo repository.CurriculumRepository,
	courseRepo repository.CourseRepository,
) SetDefaultService {
	return &setDefaultService{
		setDefaultRepo: setDefaultRepo,
		curriculumRepo: curriculumRepo,
		courseRepo:     courseRepo,
	}
}

func (s *setDefaultService) Create(setDefault *model.SetDefault) error {
	return s.setDefaultRepo.Create(setDefault)
}

func (s *setDefaultService) GetByID(id uuid.UUID) (*model.SetDefault, error) {
	return s.setDefaultRepo.GetByID(id)
}

func (s *setDefaultService) GetByCurriculumName(curriculumName string) ([]model.SetDefault, error) {
	return s.setDefaultRepo.GetByCurriculumName(curriculumName)
}

func (s *setDefaultService) GetByCurriculumID(curriculumID uuid.UUID) ([]model.SetDefault, error) {
	return s.setDefaultRepo.GetByCurriculumID(curriculumID)
}

func (s *setDefaultService) GetAll() ([]model.SetDefault, error) {
	return s.setDefaultRepo.GetAll()
}

func (s *setDefaultService) Update(setDefault *model.SetDefault) error {
	return s.setDefaultRepo.Update(setDefault)
}

func (s *setDefaultService) Delete(id uuid.UUID) error {
	return s.setDefaultRepo.Delete(id)
}

func (s *setDefaultService) ImportFromCSV(reader io.Reader) error {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	// Skip header row
	records = records[1:]

	var setDefaults []model.SetDefault

	for i, record := range records {
		if len(record) != 4 {
			return fmt.Errorf("invalid CSV format at row %d: expected 4 columns, got %d", i+2, len(record))
		}

		curriculumNames := strings.Split(record[0], ",")
		courseCode := strings.TrimSpace(record[1])

		year, err := strconv.Atoi(strings.TrimSpace(record[2]))
		if err != nil {
			return fmt.Errorf("invalid year at row %d: %v", i+2, err)
		}

		semester, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			return fmt.Errorf("invalid semester at row %d: %v", i+2, err)
		}

		// For each curriculum name in the comma-separated list
		for _, curriculumName := range curriculumNames {
			curriculumName = strings.TrimSpace(curriculumName)
			if curriculumName == "" {
				continue
			}

			// Find curriculum by name
			curriculum, err := s.curriculumRepo.GetByName(curriculumName)
			if err != nil {
				return fmt.Errorf("curriculum not found for name %s at row %d", curriculumName, i+2)
			}

			// Find the course by code and curriculum ID and year
			course, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(courseCode, curriculum.ID, year)
			if err != nil {
				return fmt.Errorf("course %s (year %d) not found in curriculum %s at row %d", courseCode, year, curriculumName, i+2)
			}

			setDefault := model.SetDefault{
				CurriculumID: curriculum.ID,
				CourseID:     course.ID,
				Year:         year,
				Semester:     semester,
			}

			setDefaults = append(setDefaults, setDefault)
		}
	}

	return s.setDefaultRepo.BulkUpsert(setDefaults)
}
