package service

import (
	"checkjop-be/internal/model"
	"checkjop-be/internal/repository"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type GraduationService interface {
	CheckGraduation(progress *model.StudentProgress) (*model.GraduationCheckResult, error)
	CheckGraduationByName(request *model.StudentProgressByName) (*model.GraduationCheckResult, error)
	ValidatePrerequisites(progress *model.StudentProgress) ([]model.PrerequisiteViolation, error)
	CheckCategoryRequirements(progress *model.StudentProgress) ([]model.CategoryCheckResult, error)
	ValidateCreditLimits(progress *model.StudentProgress) ([]model.CreditLimitViolation, error)
}

type graduationService struct {
	curriculumRepo repository.CurriculumRepository
	courseRepo     repository.CourseRepository
	categoryRepo   repository.CategoryRepository
}

func NewGraduationService(
	curriculumRepo repository.CurriculumRepository,
	courseRepo repository.CourseRepository,
	categoryRepo repository.CategoryRepository,
) GraduationService {
	return &graduationService{
		curriculumRepo: curriculumRepo,
		courseRepo:     courseRepo,
		categoryRepo:   categoryRepo,
	}
}

func (s *graduationService) CheckGraduation(progress *model.StudentProgress) (*model.GraduationCheckResult, error) {
	curriculum, err := s.curriculumRepo.GetByID(progress.CurriculumID)
	if err != nil {
		return nil, fmt.Errorf("curriculum not found: %w", err)
	}

	categoryResults, err := s.CheckCategoryRequirements(progress)
	if err != nil {
		return nil, err
	}

	prereqViolations, err := s.ValidatePrerequisites(progress)
	if err != nil {
		return nil, err
	}

	creditLimitViolations, err := s.ValidateCreditLimits(progress)
	if err != nil {
		return nil, err
	}

	unrecognizedCourses := s.findUnrecognizedCourses(progress)
	missingCatalogYears := s.findMissingCatalogYears(progress)
	catalogYearFallbacks := s.buildCatalogYearFallbacks(progress.CurriculumID, missingCatalogYears)

	totalCredits := s.calculateTotalCredits(progress.Courses)
	gpax := s.calculateGPAX(progress.Courses)
	canGraduate := s.determineGraduationEligibility(totalCredits, curriculum.MinTotalCredits, categoryResults, prereqViolations, creditLimitViolations, unrecognizedCourses)

	return &model.GraduationCheckResult{
		CanGraduate:            canGraduate,
		GPAX:                   gpax,
		TotalCredits:           totalCredits,
		RequiredCredits:        curriculum.MinTotalCredits,
		CategoryResults:        categoryResults,
		MissingCourses:         s.findMissingCourses(progress, categoryResults),
		UnrecognizedCourses:    unrecognizedCourses,
		MissingCatalogYears:    missingCatalogYears,
		CatalogYearFallbacks:   catalogYearFallbacks,
		PrerequisiteViolations: prereqViolations,
		CreditLimitViolations:  creditLimitViolations,
	}, nil
}

func (s *graduationService) CheckGraduationByName(request *model.StudentProgressByName) (*model.GraduationCheckResult, error) {
	curriculum, err := s.curriculumRepo.GetByName(request.NameTH)
	if err != nil {
		return nil, fmt.Errorf("curriculum not found: %w", err)
	}

	progress := &model.StudentProgress{
		CurriculumID:  curriculum.ID,
		Courses:       request.Courses,
		ManualCredits: request.ManualCredits,
		AdmissionYear: request.AdmissionYear,
		Exemptions:    request.Exemptions,
	}

	return s.CheckGraduation(progress)
}

func (s *graduationService) CheckCategoryRequirements(progress *model.StudentProgress) ([]model.CategoryCheckResult, error) {
	categories, err := s.categoryRepo.GetByCurriculumID(progress.CurriculumID)
	if err != nil {
		return nil, err
	}

	results := []model.CategoryCheckResult{}
	completedCourseMap := s.buildCompletedCourseMap(progress.Courses)

	for _, category := range categories {
		earnedCredits := 0

		// First check if courses specify this category directly
		for _, course := range progress.Courses {
			if course.CategoryName == category.NameTH || course.CategoryName == category.NameEN {
				earnedCredits += course.Credits
			}
		}

		// Check database courses for this category
		missingCourses := []string{}
		isElective := strings.Contains(category.NameTH, "วิชาเลือก") || strings.Contains(category.NameEN, "Elective")
		if earnedCredits == 0 {
			countedCourses := make(map[string]bool)
			for _, course := range category.Courses {
				// Only count courses matching the student's admission year
				if course.Year != progress.AdmissionYear {
					continue
				}
				if completedCourse, completed := completedCourseMap[course.Code]; completed {
					if !countedCourses[course.Code] {
						earnedCredits += completedCourse.Credits
						countedCourses[course.Code] = true
					}
				} else if !isElective {
					// Only report missing courses for non-elective categories
					missingCourses = append(missingCourses, course.Code)
				}
			}
		}

		// Fallback to manual credits if still zero
		if earnedCredits == 0 && progress.ManualCredits != nil {
			if manualCredit, exists := progress.ManualCredits[category.NameTH]; exists {
				earnedCredits = manualCredit
			} else if manualCredit, exists := progress.ManualCredits[category.NameEN]; exists {
				earnedCredits = manualCredit
			}
		}

		results = append(results, model.CategoryCheckResult{
			CategoryName:    category.NameTH,
			EarnedCredits:   earnedCredits,
			RequiredCredits: category.MinCredits,
			IsSatisfied:     earnedCredits >= category.MinCredits,
			MissingCourses:  missingCourses,
		})
	}

	return results, nil
}

func (s *graduationService) ValidatePrerequisites(progress *model.StudentProgress) ([]model.PrerequisiteViolation, error) {
	// Resolve a fallback year for exemption validation (use latest course year)
	latestYear := progress.AdmissionYear
	for _, c := range progress.Courses {
		if c.Year > latestYear {
			latestYear = c.Year
		}
	}
	exemptionYear := s.resolvePrereqYearForCourse(progress.CurriculumID, latestYear)

	// First validate exemptions
	if err := s.validateExemptions(progress.Exemptions, progress.CurriculumID, exemptionYear); err != nil {
		return nil, err
	}

	violations := []model.PrerequisiteViolation{}
	completedCourseMap := s.buildCompletedCourseMap(progress.Courses)

	for _, completedCourse := range progress.Courses {
		// Resolve catalog year per course: use course.Year, fallback to nearest earlier available
		prereqYear := s.resolvePrereqYearForCourse(progress.CurriculumID, completedCourse.Year)

		// Skip purely manual courses (have category_name but course not in DB)
		if completedCourse.CategoryName != "" {
			_, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(completedCourse.CourseCode, progress.CurriculumID, prereqYear)
			if err != nil {
				// Not in DB — purely manual, skip pre/co check
				continue
			}
			// Found in DB (e.g. free elective with real course code) — fall through to check pre/co
		}

		// Validate transitive prerequisites with proper OR/AND logic
		visited := make(map[string]bool)
		missingPrereqs, prereqsWrongTerm := s.validateTransitivePrerequisites(completedCourse.CourseCode, completedCourseMap, completedCourse, visited, progress.Exemptions, progress.CurriculumID, prereqYear)

		// Validate transitive corequisites with proper OR/AND logic
		coreqVisited := make(map[string]bool)
		missingCoreqs, coreqsWrongTerm := s.validateTransitiveCorequisites(completedCourse.CourseCode, completedCourseMap, completedCourse, coreqVisited, progress.CurriculumID, prereqYear)
		wrongTermCourse := len(coreqsWrongTerm) > 0

		if len(missingPrereqs) > 0 || len(prereqsWrongTerm) > 0 || len(missingCoreqs) > 0 || len(coreqsWrongTerm) > 0 || wrongTermCourse {
			// Deduplicate arrays to avoid showing duplicate course codes
			missingPrereqs = s.deduplicateStrings(missingPrereqs)
			prereqsWrongTerm = s.deduplicateStrings(prereqsWrongTerm)
			missingCoreqs = s.deduplicateStrings(missingCoreqs)
			coreqsWrongTerm = s.deduplicateStrings(coreqsWrongTerm)

			// Ensure empty arrays instead of nil
			if missingPrereqs == nil {
				missingPrereqs = []string{}
			}
			if prereqsWrongTerm == nil {
				prereqsWrongTerm = []string{}
			}
			if missingCoreqs == nil {
				missingCoreqs = []string{}
			}
			if coreqsWrongTerm == nil {
				coreqsWrongTerm = []string{}
			}

			violations = append(violations, model.PrerequisiteViolation{
				CourseCode:              completedCourse.CourseCode,
				MissingPrereqs:          missingPrereqs,
				PrereqsTakenInWrongTerm: prereqsWrongTerm,
				TakenInWrongTerm:        wrongTermCourse,
				MissingCoreqs:           missingCoreqs,
				CoreqsTakenInWrongTerm:  coreqsWrongTerm,
			})
		}
	}

	return violations, nil
}

func (s *graduationService) ValidateCreditLimits(progress *model.StudentProgress) ([]model.CreditLimitViolation, error) {
	const maxCreditsPerTerm = 22
	const maxCreditsPerSummer = 10 // Semester 3 (summer) limit
	violations := []model.CreditLimitViolation{}

	// Group courses by year and semester
	termCredits := make(map[string]int)

	for _, course := range progress.Courses {
		termKey := fmt.Sprintf("%d-%d", course.Year, course.Semester)
		termCredits[termKey] += course.Credits
	}

	// Check each term for credit limit violations
	for termKey, credits := range termCredits {
		// Parse the term key to get year and semester
		var year, semester int
		fmt.Sscanf(termKey, "%d-%d", &year, &semester)

		// Determine max credits based on semester
		var maxCredits int
		if semester == 3 {
			maxCredits = maxCreditsPerSummer // Summer semester
		} else {
			maxCredits = maxCreditsPerTerm // Regular semesters (1, 2)
		}

		if credits > maxCredits {
			violations = append(violations, model.CreditLimitViolation{
				Year:       year,
				Semester:   semester,
				Credits:    credits,
				MaxCredits: maxCredits,
			})
		}
	}
	return violations, nil
}

func (s *graduationService) isSameTerm(course1, course2 model.CompletedCourse) bool {
	return course1.Year == course2.Year && course1.Semester == course2.Semester
}

func (s *graduationService) isCourseTakenBefore(course1, course2 model.CompletedCourse) bool {
	if course1.Year < course2.Year {
		return true
	}
	return course1.Year == course2.Year && course1.Semester < course2.Semester
}

// validateTransitivePrerequisites checks all transitive prerequisites for a course recursively
func (s *graduationService) validateTransitivePrerequisites(courseCode string, completedMap map[string]model.CompletedCourse, currentCourse model.CompletedCourse, visited map[string]bool, exemptions []string, curriculumID uuid.UUID, admissionYear int) ([]string, []string) {
	if visited[courseCode] {
		return []string{}, []string{} // Avoid infinite loops
	}
	visited[courseCode] = true

	course, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(courseCode, curriculumID, admissionYear)
	if err != nil {
		return []string{}, []string{}
	}

	var allMissing []string
	var allWrongTerm []string

	// Check each prerequisite group
	for _, group := range course.PrerequisiteGroups {
		// Check for C.F. condition
		if group.HasCFCondition {
			hasPermission := false
			for _, exemption := range exemptions {
				if exemption == currentCourse.CourseCode {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				continue // Skip this group as it is satisfied by C.F.
			}
			// If no permission and this group only has C.F. (no course codes), it's a violation
			if len(group.PrerequisiteCourses) == 0 {
				allMissing = append(allMissing, "C.F.")
				continue
			}
		}

		if group.IsOrGroup {
			// OR logic: at least one course in the group must be satisfied
			groupSatisfied := false
			var groupMissing []string
			var groupWrongTerm []string

			for _, link := range group.PrerequisiteCourses {
				prereqCode := link.PrerequisiteCourse.Code
				if completed, exists := completedMap[prereqCode]; exists {
					// Prerequisite must be passed (not F) and taken before
					if completed.Grade != "F" && s.isCourseTakenBefore(completed, currentCourse) {
						groupSatisfied = true
						// If this prerequisite is satisfied, also check its transitive prerequisites
						transitMissing, transitWrongTerm := s.validateTransitivePrerequisites(prereqCode, completedMap, currentCourse, visited, exemptions, curriculumID, admissionYear)
						allMissing = append(allMissing, transitMissing...)
						allWrongTerm = append(allWrongTerm, transitWrongTerm...)
						break // OR condition satisfied, no need to check other options
					} else if completed.Grade == "F" {
						groupMissing = append(groupMissing, prereqCode)
					} else {
						groupWrongTerm = append(groupWrongTerm, prereqCode)
					}
				} else {
					groupMissing = append(groupMissing, prereqCode)
				}
			}

			if !groupSatisfied {
				// Add all individual courses from the OR group as missing
				allMissing = append(allMissing, groupMissing...)
				allWrongTerm = append(allWrongTerm, groupWrongTerm...)

				// Check transitive prerequisites for all courses in the OR group
				for _, link := range group.PrerequisiteCourses {
					prereqCode := link.PrerequisiteCourse.Code
					transitMissing, transitWrongTerm := s.validateTransitivePrerequisites(prereqCode, completedMap, currentCourse, visited, exemptions, curriculumID, admissionYear)
					allMissing = append(allMissing, transitMissing...)
					allWrongTerm = append(allWrongTerm, transitWrongTerm...)
				}
			}
		} else {
			// AND logic: all courses in the group must be satisfied
			for _, link := range group.PrerequisiteCourses {
				prereqCode := link.PrerequisiteCourse.Code
				if completed, exists := completedMap[prereqCode]; !exists {
					allMissing = append(allMissing, prereqCode)
					// Even if this prerequisite is missing, check what its prerequisites would be
					transitMissing, transitWrongTerm := s.validateTransitivePrerequisites(prereqCode, completedMap, currentCourse, visited, exemptions, curriculumID, admissionYear)
					allMissing = append(allMissing, transitMissing...)
					allWrongTerm = append(allWrongTerm, transitWrongTerm...)
				} else if completed.Grade == "F" {
					// Treat F as missing prerequisite
					allMissing = append(allMissing, prereqCode)
				} else if !s.isCourseTakenBefore(completed, currentCourse) {
					allWrongTerm = append(allWrongTerm, prereqCode)
				} else {
					// If this prerequisite is satisfied, also check its transitive prerequisites
					transitMissing, transitWrongTerm := s.validateTransitivePrerequisites(prereqCode, completedMap, currentCourse, visited, exemptions, curriculumID, admissionYear)
					allMissing = append(allMissing, transitMissing...)
					allWrongTerm = append(allWrongTerm, transitWrongTerm...)
				}
			}
		}
	}

	return allMissing, allWrongTerm
}

// validateTransitiveCorequisites checks all transitive corequisites for a course recursively
func (s *graduationService) validateTransitiveCorequisites(courseCode string, completedMap map[string]model.CompletedCourse, currentCourse model.CompletedCourse, visited map[string]bool, curriculumID uuid.UUID, admissionYear int) ([]string, []string) {
	if visited[courseCode] {
		return []string{}, []string{} // Avoid infinite loops
	}
	visited[courseCode] = true

	course, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(courseCode, curriculumID, admissionYear)
	if err != nil {
		return []string{}, []string{}
	}

	var allMissing []string
	var allWrongTerm []string

	// Check each corequisite group
	for _, group := range course.CorequisiteGroups {
		if group.IsOrGroup {
			// OR logic: at least one course in the group must be satisfied
			groupSatisfied := false
			var groupMissing []string
			var groupWrongTerm []string

			for _, link := range group.PrerequisiteCourses {
				coreqCode := link.PrerequisiteCourse.Code
				if completed, exists := completedMap[coreqCode]; exists {
					if s.isSameTerm(completed, currentCourse) {
						groupSatisfied = true
						// If this corequisite is satisfied, also check its transitive corequisites
						transitMissing, transitWrongTerm := s.validateTransitiveCorequisites(coreqCode, completedMap, currentCourse, visited, curriculumID, admissionYear)
						allMissing = append(allMissing, transitMissing...)
						allWrongTerm = append(allWrongTerm, transitWrongTerm...)
						break // OR condition satisfied, no need to check other options
					} else {
						groupWrongTerm = append(groupWrongTerm, coreqCode)
					}
				} else {
					groupMissing = append(groupMissing, coreqCode)
				}
			}

			if !groupSatisfied {
				// Add all individual courses from the OR group as missing
				allMissing = append(allMissing, groupMissing...)
				allWrongTerm = append(allWrongTerm, groupWrongTerm...)

				// Check transitive corequisites for all courses in the OR group
				for _, link := range group.PrerequisiteCourses {
					coreqCode := link.PrerequisiteCourse.Code
					transitMissing, transitWrongTerm := s.validateTransitiveCorequisites(coreqCode, completedMap, currentCourse, visited, curriculumID, admissionYear)
					allMissing = append(allMissing, transitMissing...)
					allWrongTerm = append(allWrongTerm, transitWrongTerm...)
				}
			}
		} else {
			// AND logic: all courses in the group must be satisfied
			for _, link := range group.PrerequisiteCourses {
				coreqCode := link.PrerequisiteCourse.Code
				if completed, exists := completedMap[coreqCode]; !exists {
					allMissing = append(allMissing, coreqCode)
					// Even if this corequisite is missing, check what its corequisites would be
					transitMissing, transitWrongTerm := s.validateTransitiveCorequisites(coreqCode, completedMap, currentCourse, visited, curriculumID, admissionYear)
					allMissing = append(allMissing, transitMissing...)
					allWrongTerm = append(allWrongTerm, transitWrongTerm...)
				} else if !s.isSameTerm(completed, currentCourse) {
					allWrongTerm = append(allWrongTerm, coreqCode)
				} else {
					// If this corequisite is satisfied, also check its transitive corequisites
					transitMissing, transitWrongTerm := s.validateTransitiveCorequisites(coreqCode, completedMap, currentCourse, visited, curriculumID, admissionYear)
					allMissing = append(allMissing, transitMissing...)
					allWrongTerm = append(allWrongTerm, transitWrongTerm...)
				}
			}
		}
	}

	return allMissing, allWrongTerm
}

func (s *graduationService) buildCompletedCourseMap(courses []model.CompletedCourse) map[string]model.CompletedCourse {
	courseMap := make(map[string]model.CompletedCourse)
	for _, course := range courses {
		courseMap[course.CourseCode] = course
	}
	return courseMap
}

func (s *graduationService) calculateTotalCredits(courses []model.CompletedCourse) int {
	total := 0
	for _, course := range courses {
		total += course.Credits
	}
	return total
}

func (s *graduationService) determineGraduationEligibility(totalCredits, requiredCredits int, categoryResults []model.CategoryCheckResult, violations []model.PrerequisiteViolation, creditViolations []model.CreditLimitViolation, unrecognizedCourses []string) bool {
	if totalCredits < requiredCredits {
		return false
	}

	if len(unrecognizedCourses) > 0 {
		return false
	}

	for _, result := range categoryResults {
		if !result.IsSatisfied {
			return false
		}
	}

	return len(violations) == 0 && len(creditViolations) == 0
}

func (s *graduationService) findMissingCourses(progress *model.StudentProgress, categoryResults []model.CategoryCheckResult) []string {
	var missing []string
	for _, result := range categoryResults {
		if !result.IsSatisfied {
			missing = append(missing, fmt.Sprintf("หมวด %s ขาด %d หน่วยกิต",
				result.CategoryName, result.RequiredCredits-result.EarnedCredits))
		}
	}
	return missing
}

// isFreeElectiveCategory returns true if the category name refers to a free elective.
func isFreeElectiveCategory(categoryName string) bool {
	return categoryName == "วิชาเสรี" || categoryName == "Free Elective"
}

// findUnrecognizedCourses returns course_codes that are expected to exist in the DB
// but cannot be found for the given curriculum and admission year.
//
// Rules:
//   - No category_name → must exist in DB for this year; if not → unrecognized.
//   - Has category_name, not free elective → custom GEN ED entry → skip entirely.
//   - Has category_name, is free elective → check DB:
//     found → validate normally (pre/co checked elsewhere).
//     not found → treated as a custom free elective entry → skip.
func (s *graduationService) findUnrecognizedCourses(progress *model.StudentProgress) []string {
	var unrecognized []string
	for _, course := range progress.Courses {
		catalogYear := s.resolvePrereqYearForCourse(progress.CurriculumID, course.Year)
		if course.CategoryName != "" {
			if !isFreeElectiveCategory(course.CategoryName) {
				// GEN ED or other manual category — skip entirely
				continue
			}
			// Free elective with a course code: check DB, skip if not found
			_, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(course.CourseCode, progress.CurriculumID, catalogYear)
			if err != nil {
				// Not in DB — treat as custom free elective entry, skip
				continue
			}
			// Found in DB — not unrecognized, pre/co validated elsewhere
			continue
		}
		// No category_name: must exist in DB for this course's catalog year (with fallback)
		_, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(course.CourseCode, progress.CurriculumID, catalogYear)
		if err != nil {
			unrecognized = append(unrecognized, course.CourseCode)
		}
	}
	if unrecognized == nil {
		return []string{}
	}
	return unrecognized
}

// resolvePrereqYearForCourse returns the catalog year to use for pre/co req checks for a specific course year.
// It uses courseYear if that catalog exists; otherwise falls back to the nearest earlier year that has catalog data.
// Returns courseYear as-is if no fallback found.
func (s *graduationService) resolvePrereqYearForCourse(curriculumID uuid.UUID, courseYear int) int {
	if s.courseRepo.CatalogYearExists(curriculumID, courseYear) {
		return courseYear
	}
	if year, ok := s.courseRepo.GetLatestAvailableCatalogYear(curriculumID, courseYear-1); ok {
		return year
	}
	return courseYear
}

// buildCatalogYearFallbacks returns a map of missing catalog year → fallback catalog year used for pre/co req checks.
func (s *graduationService) buildCatalogYearFallbacks(curriculumID uuid.UUID, missingYears []int) map[int]int {
	result := make(map[int]int)
	for _, y := range missingYears {
		result[y] = s.resolvePrereqYearForCourse(curriculumID, y)
	}
	return result
}

// findMissingCatalogYears returns academic years (course.Year) that the student took courses in
// but for which no catalog has been imported into the DB yet.
// These years need attention because pre/co req changes may not be reflected.
func (s *graduationService) findMissingCatalogYears(progress *model.StudentProgress) []int {
	checked := make(map[int]bool)
	var missing []int
	for _, course := range progress.Courses {
		y := course.Year
		if checked[y] {
			continue
		}
		checked[y] = true
		if !s.courseRepo.CatalogYearExists(progress.CurriculumID, y) {
			missing = append(missing, y)
		}
	}
	if missing == nil {
		return []int{}
	}
	return missing
}

func (s *graduationService) deduplicateStrings(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func (s *graduationService) calculateGPAX(courses []model.CompletedCourse) float64 {
	totalPoints := 0.0
	totalCredits := 0

	for _, course := range courses {
		points := 0.0
		isGraded := true

		switch course.Grade {
		case "A":
			points = 4.0
		case "B+":
			points = 3.5
		case "B":
			points = 3.0
		case "C+":
			points = 2.5
		case "C":
			points = 2.0
		case "D+":
			points = 1.5
		case "D":
			points = 1.0
		case "F":
			points = 0.0
		default:
			isGraded = false // W, S, U, etc. don't count towards GPA
		}

		if isGraded {
			totalPoints += points * float64(course.Credits)
			totalCredits += course.Credits
		}
	}

	if totalCredits == 0 {
		return 0.0
	}

	return totalPoints / float64(totalCredits)
}

// validateExemptions checks if all exempted courses are eligible for C.F.
// It checks the given admissionYear first, then falls back to any available year that has CF enabled.
func (s *graduationService) validateExemptions(exemptions []string, curriculumID uuid.UUID, admissionYear int) error {
	for _, courseCode := range exemptions {
		course, err := s.courseRepo.GetByCodeAndCurriculumIDAndYear(courseCode, curriculumID, admissionYear)
		if err != nil {
			// Course not found in this catalog year — check any year
			if s.courseRepo.CourseHasCFOptionInAnyCatalogYear(courseCode, curriculumID) {
				continue
			}
			return fmt.Errorf("exemption course '%s' not found in curriculum", courseCode)
		}
		if !course.HasCFOption {
			// Fallback: check if any catalog year has CF for this course
			if !s.courseRepo.CourseHasCFOptionInAnyCatalogYear(courseCode, curriculumID) {
				return fmt.Errorf("course '%s' does not allow C.F. exemption (no C.F. option in prerequisites or corequisites)", courseCode)
			}
		}
	}

	return nil
}
