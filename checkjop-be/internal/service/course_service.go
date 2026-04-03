package service

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/repository"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CourseService interface {
	Create(course *model.Course) error
	GetByID(id uuid.UUID) (*model.Course, error)
	GetByCodeAndYear(code string, year int) (*model.Course, error)
	GetByCodeAndCurriculumIDAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error)
	GetByName(name string) (*model.Course, error)
	GetByCurriculumID(curriculumID uuid.UUID) ([]model.Course, error)
	GetAll() ([]model.Course, error)
	Update(course *model.Course) error
	Delete(id uuid.UUID) error
	ImportFromCSV(reader io.Reader) error
	ImportFromCSVWithYear(reader io.Reader, year int) error
	SetCourseRelationships(courseID uuid.UUID, prerequisiteIDs []uuid.UUID, corequisiteIDs []uuid.UUID) error
	SetCourseRelationshipsWithGroups(courseID uuid.UUID, prerequisiteGroups []model.PrerequisiteGroup, corequisiteGroups []model.PrerequisiteGroup) error
	ResetDatabase() error
}

type courseService struct {
	courseRepo     repository.CourseRepository
	categoryRepo   repository.CategoryRepository
	curriculumRepo repository.CurriculumRepository
}

func NewCourseService(
	courseRepo repository.CourseRepository,
	categoryRepo repository.CategoryRepository,
	curriculumRepo repository.CurriculumRepository,
) CourseService {
	return &courseService{
		courseRepo:     courseRepo,
		categoryRepo:   categoryRepo,
		curriculumRepo: curriculumRepo,
	}
}

func (s *courseService) Create(course *model.Course) error {
	return s.courseRepo.Create(course)
}

func (s *courseService) GetByID(id uuid.UUID) (*model.Course, error) {
	return s.courseRepo.GetByID(id)
}

func (s *courseService) GetByCodeAndYear(code string, year int) (*model.Course, error) {
	return s.courseRepo.GetByCodeAndYear(code, year)
}

func (s *courseService) GetByCodeAndCurriculumIDAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error) {
	return s.courseRepo.GetByCodeAndCurriculumIDAndYear(code, curriculumID, year)
}

func (s *courseService) GetByName(name string) (*model.Course, error) {
	return s.courseRepo.GetByName(name)
}

func (s *courseService) GetByCurriculumID(curriculumID uuid.UUID) ([]model.Course, error) {
	return s.courseRepo.GetByCurriculumID(curriculumID)
}

func (s *courseService) GetAll() ([]model.Course, error) {
	return s.courseRepo.GetAll()
}

func (s *courseService) Update(course *model.Course) error {
	return s.courseRepo.Update(course)
}

func (s *courseService) Delete(id uuid.UUID) error {
	return s.courseRepo.Delete(id)
}

func (s *courseService) ImportFromCSV(reader io.Reader) error {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV must have header and at least one data row")
	}

	// Clear existing courses before import
	if err := s.courseRepo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to clear existing courses: %w", err)
	}

	// CSV format: code,courseNameEN,courseNameTH,credit,pre,co,category,curriculum,Year
	if len(records[0]) < 9 {
		return fmt.Errorf("CSV must have at least 9 columns")
	}

	// Build curriculum cache
	curriculumCache := make(map[string]*model.Curriculum)
	// Build category cache per curriculum
	categoryCache := make(map[string]map[string]uuid.UUID)

	var courses []model.Course
	// Use a map to deduplicate courses based on code + curriculum_id
	courseMap := make(map[string]model.Course)

	for i, record := range records[1:] {
		if len(record) < 9 {
			return fmt.Errorf("row %d: insufficient columns", i+2)
		}

		// Parse multiple curricula (comma-separated)
		curriculumNames := strings.Split(strings.TrimSpace(record[7]), ",")
		// Parse multiple categories (comma-separated)
		categoryNames := strings.Split(strings.TrimSpace(record[6]), ",")

		credits, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			return fmt.Errorf("row %d: invalid credits: %w", i+2, err)
		}

		year, err := strconv.Atoi(strings.TrimSpace(record[8]))
		if err != nil {
			return fmt.Errorf("row %d: invalid year: %w", i+2, err)
		}

		// Create course for each curriculum and category combination
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

			// Get or cache categories for this curriculum
			categoryMappings, exists := categoryCache[currName]
			if !exists {
				categories, err := s.categoryRepo.GetByCurriculumID(curriculum.ID)
				if err != nil {
					return fmt.Errorf("row %d: failed to get categories for curriculum '%s': %w", i+2, currName, err)
				}

				categoryMappings = make(map[string]uuid.UUID)
				for _, cat := range categories {
					categoryMappings[cat.NameTH] = cat.ID
					categoryMappings[cat.NameEN] = cat.ID
				}
				categoryCache[currName] = categoryMappings
			}

			// Create course for each category
			for _, categoryName := range categoryNames {
				categoryName = strings.TrimSpace(categoryName)
				if categoryName == "" {
					continue
				}

				categoryID, exists := categoryMappings[categoryName]
				if !exists {
					return fmt.Errorf("row %d: category '%s' not found in curriculum '%s'", i+2, categoryName, currName)
				}

				course := model.Course{
					Code:         strings.TrimSpace(record[0]),
					NameEN:       strings.TrimSpace(record[1]),
					NameTH:       strings.TrimSpace(record[2]),
					Credits:      credits,
					CurriculumID: curriculum.ID,
					CategoryID:   categoryID,
					Year:         year,
				}

				// Create unique key for deduplication (include category and year to allow same course in different categories/years)
				key := fmt.Sprintf("%s_%s_%s_%d", course.Code, curriculum.ID.String(), categoryID.String(), course.Year)
				courseMap[key] = course
			}
		}
	}

	// Convert map to slice
	for _, course := range courseMap {
		courses = append(courses, course)
	}

	// First pass: Import all courses
	if err := s.courseRepo.BulkUpsert(courses); err != nil {
		return fmt.Errorf("failed to upsert courses: %w", err)
	}

	// Second pass: Set up prerequisites and corequisites relationships
	return s.setupCourseRelationships(records[1:])
}

// SetCourseRelationships sets prerequisite and corequisite relationships for a course
func (s *courseService) SetCourseRelationships(courseID uuid.UUID, prerequisiteIDs []uuid.UUID, corequisiteIDs []uuid.UUID) error {
	if err := s.courseRepo.SetPrerequisites(courseID, prerequisiteIDs); err != nil {
		return fmt.Errorf("failed to set prerequisites: %w", err)
	}
	if err := s.courseRepo.SetCorequisites(courseID, corequisiteIDs); err != nil {
		return fmt.Errorf("failed to set corequisites: %w", err)
	}
	return nil
}

// SetCourseRelationshipsWithGroups sets prerequisite and corequisite groups with OR/AND logic for a course
func (s *courseService) SetCourseRelationshipsWithGroups(courseID uuid.UUID, prerequisiteGroups []model.PrerequisiteGroup, corequisiteGroups []model.PrerequisiteGroup) error {
	if err := s.courseRepo.SetPrerequisiteGroups(courseID, prerequisiteGroups); err != nil {
		return fmt.Errorf("failed to set prerequisite groups: %w", err)
	}
	if err := s.courseRepo.SetCorequisiteGroups(courseID, corequisiteGroups); err != nil {
		return fmt.Errorf("failed to set corequisite groups: %w", err)
	}
	return nil
}

// setupCourseRelationships processes CSV records to establish course relationships
func (s *courseService) setupCourseRelationships(records [][]string) error {
	// Build curriculum cache
	curriculumCache := make(map[string]*model.Curriculum)

	// Collect all errors instead of failing on first error
	var allErrors []string

	for i, record := range records {
		if len(record) < 9 {
			continue // Skip invalid records
		}

		courseCode := strings.TrimSpace(record[0])
		prerequisitesStr := strings.TrimSpace(record[4])
		corequisitesStr := strings.TrimSpace(record[5])
		curriculumNames := strings.Split(strings.TrimSpace(record[7]), ",")
		year, err := strconv.Atoi(strings.TrimSpace(record[8]))
		if err != nil {
			continue // Skip invalid year
		}

		// Skip if no relationships to set up
		if prerequisitesStr == "" && corequisitesStr == "" {
			continue
		}

		// Process for each curriculum this course belongs to
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
					allErrors = append(allErrors, fmt.Sprintf("row %d: curriculum '%s' not found", i+2, currName))
					continue
				}
				curriculum = curr
				curriculumCache[currName] = curriculum
			}

			// Find the course in this curriculum and year
			course, err := s.findCourseByCodeAndCurriculumAndYear(courseCode, curriculum.ID, year)
			if err != nil {
				allErrors = append(allErrors, fmt.Sprintf("row %d: course '%s' (year %d) not found in curriculum '%s'", i+2, courseCode, year, currName))
				continue
			}

			// Parse prerequisites using new OR logic
			var prerequisiteGroups []model.PrerequisiteGroup
			hasCFOption := false
			if prerequisitesStr != "" {
				parsedGroups, err := s.parsePrerequisiteString(prerequisitesStr)
				if err != nil {
					allErrors = append(allErrors, fmt.Sprintf("row %d: failed to parse prerequisites '%s': %v", i+2, prerequisitesStr, err))
					continue
				}

				for _, parsedGroup := range parsedGroups {
					if parsedGroup.HasCFCondition {
						hasCFOption = true
					}
					var prerequisiteCourseLinks []model.PrerequisiteCourseLink
					for _, courseCode := range parsedGroup.CourseCodes {
						prereqCourse, err := s.findCourseByCodeAndCurriculumAndYear(courseCode, curriculum.ID, year)
						if err != nil {
							allErrors = append(allErrors, fmt.Sprintf("row %d: prerequisite course '%s' (year %d) not found in curriculum '%s' for course '%s'", i+2, courseCode, year, currName, course.Code))
							continue
						}
						prerequisiteCourseLinks = append(prerequisiteCourseLinks, model.PrerequisiteCourseLink{
							PrerequisiteCourseID: prereqCourse.ID,
						})
					}

					prerequisiteGroups = append(prerequisiteGroups, model.PrerequisiteGroup{
						IsOrGroup:           parsedGroup.IsOrGroup,
						HasCFCondition:      parsedGroup.HasCFCondition,
						PrerequisiteCourses: prerequisiteCourseLinks,
					})
				}
			}

			// Parse corequisites using new OR logic
			var corequisiteGroups []model.PrerequisiteGroup
			if corequisitesStr != "" {
				parsedGroups, err := s.parsePrerequisiteString(corequisitesStr)
				if err != nil {
					allErrors = append(allErrors, fmt.Sprintf("row %d: failed to parse corequisites '%s': %v", i+2, corequisitesStr, err))
					continue
				}

				for _, parsedGroup := range parsedGroups {
					if parsedGroup.HasCFCondition {
						hasCFOption = true
					}
					var corequisiteCourseLinks []model.PrerequisiteCourseLink
					for _, courseCode := range parsedGroup.CourseCodes {
						coreqCourse, err := s.findCourseByCodeAndCurriculumAndYear(courseCode, curriculum.ID, year)
						if err != nil {
							allErrors = append(allErrors, fmt.Sprintf("row %d: corequisite course '%s' (year %d) not found in curriculum '%s' for course '%s'", i+2, courseCode, year, currName, course.Code))
							continue
						}
						corequisiteCourseLinks = append(corequisiteCourseLinks, model.PrerequisiteCourseLink{
							PrerequisiteCourseID: coreqCourse.ID,
						})
					}

					corequisiteGroups = append(corequisiteGroups, model.PrerequisiteGroup{
						IsOrGroup:           parsedGroup.IsOrGroup,
						HasCFCondition:      parsedGroup.HasCFCondition,
						PrerequisiteCourses: corequisiteCourseLinks,
					})
				}
			}

			// Update course HasCFOption if C.F. was found in prerequisites or corequisites
			if hasCFOption {
				course.HasCFOption = true
				if err := s.courseRepo.Update(course); err != nil {
					allErrors = append(allErrors, fmt.Sprintf("row %d: failed to update HasCFOption for course '%s': %v", i+2, courseCode, err))
					continue
				}
			}

			// Set the relationships using new groups
			if err := s.SetCourseRelationshipsWithGroups(course.ID, prerequisiteGroups, corequisiteGroups); err != nil {
				allErrors = append(allErrors, fmt.Sprintf("row %d: failed to set relationships for course '%s': %v", i+2, courseCode, err))
				continue
			}
		}
	}

	// If there were any errors, return them all at once
	if len(allErrors) > 0 {
		return fmt.Errorf("found %d error(s) during course relationship setup:\n- %s", len(allErrors), strings.Join(allErrors, "\n- "))
	}

	return nil
}

// findCourseByCodeAndCurriculumAndYear finds a course by code within a specific curriculum and year
// If not found in the specific curriculum, it falls back to finding ANY course with the same code and year
func (s *courseService) findCourseByCodeAndCurriculumAndYear(code string, curriculumID uuid.UUID, year int) (*model.Course, error) {
	// Try to find in the specific curriculum first
	course, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(code, curriculumID, year)
	if err == nil {
		return course, nil
	}

	// Fallback: Try to find any course with the same code and year
	// This handles cross-curriculum prerequisites
	course, err = s.courseRepo.GetByCodeAndYear(code, year)
	if err == nil {
		return course, nil
	}

	return nil, err
}

// parsePrerequisiteString parses a prerequisite string that may contain OR/AND logic
// It converts the expression into Conjunctive Normal Form (CNF): (A OR B) AND (C OR D) ...
func (s *courseService) parsePrerequisiteString(prerequisiteStr string) ([]ParsedPrerequisiteGroup, error) {
	if prerequisiteStr == "" {
		return []ParsedPrerequisiteGroup{}, nil
	}

	// Replace "AND" with "," to normalize AND delimiters
	// This handles "A AND B" -> "A, B"
	re := regexp.MustCompile(`(?i)\s+AND\s+`)
	prerequisiteStr = re.ReplaceAllString(prerequisiteStr, ",")

	// Parse into CNF structure: list of groups (AND), where each group is a list of codes (OR)
	cnfGroups := s.parseToCNF(prerequisiteStr)

	var groups []ParsedPrerequisiteGroup
	for _, codes := range cnfGroups {
		group := ParsedPrerequisiteGroup{
			CourseCodes:    []string{},
			IsOrGroup:      len(codes) > 1, // If more than 1 code in CNF group, it's an OR group
			HasCFCondition: false,
		}

		// If single item, IsOrGroup is false (AND logic for single item is same as OR logic)
		// But in our model, IsOrGroup=true means "Satisfy ANY". IsOrGroup=false means "Satisfy ALL".
		// For a group of 1, "ANY" == "ALL".
		// However, to be consistent with previous behavior, single items are usually IsOrGroup=false (unless explicitly OR).
		// But here, CNF groups are inherently OR groups.
		// If [A, B], it means A OR B.
		// If [A], it means A.
		// So IsOrGroup=true is correct for >1.
		// For =1, it doesn't matter, but false is cleaner.

		for _, code := range codes {
			code = strings.TrimSpace(code)
			if strings.ToUpper(code) == "C.F." {
				group.HasCFCondition = true
			} else if code != "" {
				group.CourseCodes = append(group.CourseCodes, code)
			}
		}

		if len(group.CourseCodes) > 0 || group.HasCFCondition {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// parseToCNF recursively parses the string into CNF: list of (list of codes)
// Outer list is AND, inner list is OR.
func (s *courseService) parseToCNF(str string) [][]string {
	str = strings.TrimSpace(str)
	if str == "" {
		return [][]string{}
	}

	// 1. Split by comma (AND) outside parentheses
	parts := s.splitByDelimiterOutsideParentheses(str, ',')
	if len(parts) > 1 {
		var result [][]string
		for _, part := range parts {
			result = append(result, s.parseToCNF(part)...)
		}
		return result
	}

	// 2. Split by OR outside parentheses
	// Note: We need to handle case-insensitive " OR "
	orParts := s.splitByOROutsideParentheses(str)
	if len(orParts) > 1 {
		var cnfs [][][]string
		for _, part := range orParts {
			cnfs = append(cnfs, s.parseToCNF(part))
		}

		// Cartesian Product to merge OR branches
		// (A AND B) OR C -> CNF1=[[A],[B]], CNF2=[[C]]
		// Result = [ [A,C], [B,C] ]
		return s.mergeCNFsOr(cnfs)
	}

	// 3. Remove outer parentheses and recurse
	if strings.HasPrefix(str, "(") && strings.HasSuffix(str, ")") {
		// Check if parens enclose the whole string matching-ly
		// e.g. "(A) OR (B)" -> starts with ( ends with ) but shouldn't remove.
		// But we already tried splitting by OR and AND. If we are here, it means
		// the top level structure is likely enclosed in parens or is a single atom.
		// We can try removing parens.
		inner := str[1 : len(str)-1]
		// Verify balancing? splitBy... handles balancing.
		// If we are here, it's either "(...)" or "Atom".
		// If it was "(A) OR (B)", splitByOR would have caught it.
		// So it must be "(...)" or "Atom".
		// Let's check if it's truly enclosed.
		if s.isEnclosedInParentheses(str) {
			return s.parseToCNF(inner)
		}
	}

	// 4. Base case: Single code
	return [][]string{{str}}
}

// mergeCNFsOr merges multiple CNF structures using Cartesian Product (OR logic)
func (s *courseService) mergeCNFsOr(cnfs [][][]string) [][]string {
	if len(cnfs) == 0 {
		return [][]string{}
	}

	result := cnfs[0]
	for i := 1; i < len(cnfs); i++ {
		nextCNF := cnfs[i]
		var newResult [][]string
		for _, g1 := range result {
			for _, g2 := range nextCNF {
				// Union of g1 and g2
				// g1 and g2 are lists of codes (OR groups)
				// Merging them creates a larger OR group
				merged := make([]string, len(g1)+len(g2))
				copy(merged, g1)
				copy(merged[len(g1):], g2)
				newResult = append(newResult, merged)
			}
		}
		result = newResult
	}
	return result
}

// splitByDelimiterOutsideParentheses splits a string by a delimiter char outside parentheses
func (s *courseService) splitByDelimiterOutsideParentheses(str string, delimiter rune) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for _, char := range str {
		switch char {
		case '(':
			depth++
			current.WriteRune(char)
		case ')':
			depth--
			current.WriteRune(char)
		case delimiter:
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// splitByOROutsideParentheses splits a string by " OR " (case-insensitive) outside parentheses
func (s *courseService) splitByOROutsideParentheses(str string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	// We need to look ahead/behind for " OR ".
	// Simpler approach: Iterate and buffer. Check buffer suffix.

	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		switch char {
		case '(':
			depth++
			current.WriteRune(char)
		case ')':
			depth--
			current.WriteRune(char)
		default:
			current.WriteRune(char)
			// Check for " OR " split
			if depth == 0 {
				buf := current.String()
				if len(buf) >= 4 {
					suffix := buf[len(buf)-4:]
					if strings.EqualFold(suffix, " OR ") {
						// Found split
						part := buf[:len(buf)-4]
						parts = append(parts, part)
						current.Reset()
					}
				}
			}
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// isEnclosedInParentheses checks if the string is fully enclosed in parentheses
// e.g. "(A)" -> true, "(A) OR (B)" -> false
func (s *courseService) isEnclosedInParentheses(str string) bool {
	if !strings.HasPrefix(str, "(") || !strings.HasSuffix(str, ")") {
		return false
	}
	depth := 0
	for i, char := range str {
		if char == '(' {
			depth++
		} else if char == ')' {
			depth--
		}
		if depth == 0 && i < len(str)-1 {
			return false // Closed before end
		}
	}
	return depth == 0
}

// ParsedPrerequisiteGroup represents a parsed prerequisite group with raw course codes
type ParsedPrerequisiteGroup struct {
	CourseCodes    []string
	IsOrGroup      bool
	HasCFCondition bool
}

// ImportFromCSVWithYear imports courses from CSV with a specified year (for version 4 CSV without Year column)
func (s *courseService) ImportFromCSVWithYear(reader io.Reader, year int) error {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV must have header and at least one data row")
	}

	// Hard delete existing courses for this year before import
	if err := s.courseRepo.DeleteByYear(year); err != nil {
		return fmt.Errorf("failed to clear existing courses for year %d: %w", year, err)
	}

	// CSV format for version 4: code,courseNameEN,courseNameTH,credit,prerequisites,corequisites,category,curriculum
	// (no Year column - we'll add it from the parameter)
	if len(records[0]) < 8 {
		return fmt.Errorf("CSV must have at least 8 columns")
	}

	// Build curriculum cache
	curriculumCache := make(map[string]*model.Curriculum)
	// Build category cache per curriculum
	categoryCache := make(map[string]map[string]uuid.UUID)

	var courses []model.Course
	// Use a map to deduplicate courses based on code + curriculum_id
	courseMap := make(map[string]model.Course)

	for i, record := range records[1:] {
		if len(record) < 8 {
			return fmt.Errorf("row %d: insufficient columns", i+2)
		}

		// Parse multiple curricula (comma-separated)
		curriculumNames := strings.Split(strings.TrimSpace(record[7]), ",")
		// Parse multiple categories (comma-separated)
		categoryNames := strings.Split(strings.TrimSpace(record[6]), ",")

		credits, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			return fmt.Errorf("row %d: invalid credits: %w", i+2, err)
		}

		// Create course for each curriculum and category combination
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

			// Get or cache categories for this curriculum
			categoryMappings, exists := categoryCache[currName]
			if !exists {
				categories, err := s.categoryRepo.GetByCurriculumID(curriculum.ID)
				if err != nil {
					return fmt.Errorf("row %d: failed to get categories for curriculum '%s': %w", i+2, currName, err)
				}

				categoryMappings = make(map[string]uuid.UUID)
				for _, cat := range categories {
					categoryMappings[cat.NameTH] = cat.ID
					categoryMappings[cat.NameEN] = cat.ID
				}
				categoryCache[currName] = categoryMappings
			}

			// Create course for each category
			for _, categoryName := range categoryNames {
				categoryName = strings.TrimSpace(categoryName)
				if categoryName == "" {
					continue
				}

				categoryID, exists := categoryMappings[categoryName]
				if !exists {
					return fmt.Errorf("row %d: category '%s' not found in curriculum '%s'", i+2, categoryName, currName)
				}

				course := model.Course{
					Code:         strings.TrimSpace(record[0]),
					NameEN:       strings.TrimSpace(record[1]),
					NameTH:       strings.TrimSpace(record[2]),
					Credits:      credits,
					CurriculumID: curriculum.ID,
					CategoryID:   categoryID,
					Year:         year, // Use the year parameter instead of reading from CSV
				}

				// Create unique key for deduplication (include category and year to allow same course in different categories/years)
				key := fmt.Sprintf("%s_%s_%s_%d", course.Code, curriculum.ID.String(), categoryID.String(), course.Year)
				courseMap[key] = course
			}
		}
	}

	// Convert map to slice
	for _, course := range courseMap {
		courses = append(courses, course)
	}

	// First pass: Import all courses
	if err := s.courseRepo.BulkUpsert(courses); err != nil {
		return fmt.Errorf("failed to upsert courses: %w", err)
	}

	// Second pass: Set up prerequisites and corequisites relationships
	return s.setupCourseRelationshipsWithYear(records[1:], year)
}

// setupCourseRelationshipsWithYear processes CSV records to establish course relationships for a specific year
func (s *courseService) setupCourseRelationshipsWithYear(records [][]string, year int) error {
	// Build curriculum cache
	curriculumCache := make(map[string]*model.Curriculum)

	// Collect all errors instead of failing on first error
	var allErrors []string

	for i, record := range records {
		if len(record) < 8 {
			continue // Skip invalid records
		}

		courseCode := strings.TrimSpace(record[0])
		prerequisitesStr := strings.TrimSpace(record[4])
		corequisitesStr := strings.TrimSpace(record[5])
		curriculumNames := strings.Split(strings.TrimSpace(record[7]), ",")

		// Skip if no relationships to set up
		if prerequisitesStr == "" && corequisitesStr == "" {
			continue
		}

		// Process for each curriculum this course belongs to
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
					allErrors = append(allErrors, fmt.Sprintf("row %d: curriculum '%s' not found", i+2, currName))
					continue
				}
				curriculum = curr
				curriculumCache[currName] = curriculum
			}

			// Find the course in this curriculum and year
			course, err := s.findCourseByCodeAndCurriculumAndYear(courseCode, curriculum.ID, year)
			if err != nil {
				allErrors = append(allErrors, fmt.Sprintf("row %d: course '%s' (year %d) not found in curriculum '%s'", i+2, courseCode, year, currName))
				continue
			}

			// Parse prerequisites using new OR logic
			var prerequisiteGroups []model.PrerequisiteGroup
			hasCFOption := false
			if prerequisitesStr != "" {
				parsedGroups, err := s.parsePrerequisiteString(prerequisitesStr)
				if err != nil {
					allErrors = append(allErrors, fmt.Sprintf("row %d: failed to parse prerequisites '%s': %v", i+2, prerequisitesStr, err))
					continue
				}

				for _, parsedGroup := range parsedGroups {
					if parsedGroup.HasCFCondition {
						hasCFOption = true
					}
					var prerequisiteCourseLinks []model.PrerequisiteCourseLink
					for _, courseCode := range parsedGroup.CourseCodes {
						prereqCourse, err := s.findCourseByCodeAndCurriculumAndYear(courseCode, curriculum.ID, year)
						if err != nil {
							allErrors = append(allErrors, fmt.Sprintf("row %d: prerequisite course '%s' (year %d) not found in curriculum '%s' for course '%s'", i+2, courseCode, year, currName, course.Code))
							continue
						}
						prerequisiteCourseLinks = append(prerequisiteCourseLinks, model.PrerequisiteCourseLink{
							PrerequisiteCourseID: prereqCourse.ID,
						})
					}

					prerequisiteGroups = append(prerequisiteGroups, model.PrerequisiteGroup{
						IsOrGroup:           parsedGroup.IsOrGroup,
						HasCFCondition:      parsedGroup.HasCFCondition,
						PrerequisiteCourses: prerequisiteCourseLinks,
					})
				}
			}

			// Parse corequisites using new OR logic
			var corequisiteGroups []model.PrerequisiteGroup
			if corequisitesStr != "" {
				parsedGroups, err := s.parsePrerequisiteString(corequisitesStr)
				if err != nil {
					allErrors = append(allErrors, fmt.Sprintf("row %d: failed to parse corequisites '%s': %v", i+2, corequisitesStr, err))
					continue
				}

				for _, parsedGroup := range parsedGroups {
					if parsedGroup.HasCFCondition {
						hasCFOption = true
					}
					var corequisiteCourseLinks []model.PrerequisiteCourseLink
					for _, courseCode := range parsedGroup.CourseCodes {
						coreqCourse, err := s.findCourseByCodeAndCurriculumAndYear(courseCode, curriculum.ID, year)
						if err != nil {
							allErrors = append(allErrors, fmt.Sprintf("row %d: corequisite course '%s' (year %d) not found in curriculum '%s' for course '%s'", i+2, courseCode, year, currName, course.Code))
							continue
						}
						corequisiteCourseLinks = append(corequisiteCourseLinks, model.PrerequisiteCourseLink{
							PrerequisiteCourseID: coreqCourse.ID,
						})
					}

					corequisiteGroups = append(corequisiteGroups, model.PrerequisiteGroup{
						IsOrGroup:           parsedGroup.IsOrGroup,
						HasCFCondition:      parsedGroup.HasCFCondition,
						PrerequisiteCourses: corequisiteCourseLinks,
					})
				}
			}

			// Update course HasCFOption if C.F. was found in prerequisites or corequisites
			if hasCFOption {
				course.HasCFOption = true
				if err := s.courseRepo.Update(course); err != nil {
					allErrors = append(allErrors, fmt.Sprintf("row %d: failed to update HasCFOption for course '%s': %v", i+2, courseCode, err))
					continue
				}
			}

			// Set the relationships using new groups
			if err := s.SetCourseRelationshipsWithGroups(course.ID, prerequisiteGroups, corequisiteGroups); err != nil {
				allErrors = append(allErrors, fmt.Sprintf("row %d: failed to set relationships for course '%s': %v", i+2, courseCode, err))
				continue
			}
		}
	}

	// If there were any errors, return them all at once
	if len(allErrors) > 0 {
		return fmt.Errorf("found %d error(s) during course relationship setup:\n- %s", len(allErrors), strings.Join(allErrors, "\n- "))
	}

	return nil
}

// ResetDatabase clears all data from the database (Courses, Categories, Curricula)
func (s *courseService) ResetDatabase() error {
	// Delete in reverse order of dependencies

	// 1. Delete Courses (depend on Categories and Curricula)
	if err := s.courseRepo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to delete all courses: %w", err)
	}

	// 2. Delete Categories (depend on Curricula)
	if err := s.categoryRepo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to delete all categories: %w", err)
	}

	// 3. Delete Curricula (root)
	if err := s.curriculumRepo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to delete all curricula: %w", err)
	}

	return nil
}
