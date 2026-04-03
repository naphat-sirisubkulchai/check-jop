package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type StudentProgress struct {
	CurriculumID  uuid.UUID         `json:"curriculum_id"`
	AdmissionYear int               `json:"admission_year"` // Year the student was admitted, used for course version lookup
	Courses       []CompletedCourse `json:"courses"`
	ManualCredits map[string]int    `json:"manual_credits,omitempty"` // category_name -> credits
	Exemptions    []string          `json:"exemptions,omitempty"`     // List of course codes with C.F. permission
}

type StudentProgressByName struct {
	NameTH        string            `json:"name_th"`
	AdmissionYear int               `json:"admission_year"`
	Courses       []CompletedCourse `json:"courses"`
	ManualCredits map[string]int    `json:"manual_credits,omitempty"` // category_name -> credits
	Exemptions    []string          `json:"exemptions,omitempty"`     // List of course codes with C.F. permission
}

type CompletedCourse struct {
	CourseCode   string `json:"course_code"`
	Year         int    `json:"year"`
	Semester     int    `json:"semester"`
	Grade        string `json:"grade,omitempty"`
	Credits      int    `json:"credits"`
	CategoryName string `json:"category_name,omitempty"` // For General Education courses
}

type GraduationCheckResult struct {
	CanGraduate            bool                    `json:"can_graduate"`
	GPAX                   float64                 `json:"gpax"`
	TotalCredits           int                     `json:"total_credits"`
	RequiredCredits        int                     `json:"required_credits"`
	CategoryResults        []CategoryCheckResult   `json:"category_results"`
	MissingCourses         []string                `json:"missing_courses"`
	UnrecognizedCourses    []string                `json:"unrecognized_courses"`
	MissingCatalogYears    []int                   `json:"missing_catalog_years"`
	CatalogYearFallbacks   map[int]int             `json:"catalog_year_fallbacks"` // missing year → catalog year actually used for pre/co req checks
	PrerequisiteViolations []PrerequisiteViolation `json:"prerequisite_violations"`
	CreditLimitViolations  []CreditLimitViolation  `json:"credit_limit_violations"`
}

type CategoryCheckResult struct {
	CategoryName    string `json:"category_name"`
	EarnedCredits   int    `json:"earned_credits"`
	RequiredCredits int    `json:"required_credits"`
	IsSatisfied     bool   `json:"is_satisfied"`
}

type PrerequisiteViolation struct {
	CourseCode              string   `json:"course_code"`
	MissingPrereqs          []string `json:"missing_prereqs"`
	PrereqsTakenInWrongTerm []string `json:"prereqs_taken_in_wrong_term"`
	TakenInWrongTerm        bool     `json:"taken_in_wrong_term"`
	MissingCoreqs           []string `json:"missing_coreqs"`
	CoreqsTakenInWrongTerm  []string `json:"coreqs_taken_in_wrong_term"`
}

// Custom JSON marshaling to ensure empty arrays instead of null
func (p PrerequisiteViolation) MarshalJSON() ([]byte, error) {
	type Alias PrerequisiteViolation
	if p.MissingPrereqs == nil {
		p.MissingPrereqs = []string{}
	}
	if p.PrereqsTakenInWrongTerm == nil {
		p.PrereqsTakenInWrongTerm = []string{}
	}
	if p.MissingCoreqs == nil {
		p.MissingCoreqs = []string{}
	}
	if p.CoreqsTakenInWrongTerm == nil {
		p.CoreqsTakenInWrongTerm = []string{}
	}
	return json.Marshal((Alias)(p))
}

type CreditLimitViolation struct {
	Year       int `json:"year"`
	Semester   int `json:"semester"`
	Credits    int `json:"credits"`
	MaxCredits int `json:"max_credits"`
}
