package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrerequisiteGroup represents a group of courses that can satisfy a prerequisite requirement
// IsOrGroup indicates if ANY course in the group satisfies the requirement (OR) or ALL courses are required (AND)
type PrerequisiteGroup struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CourseID       uuid.UUID `json:"course_id" gorm:"type:uuid;not null;index"`
	GroupType      string    `json:"group_type" gorm:"not null"`       // "prerequisite" or "corequisite"
	IsOrGroup      bool      `json:"is_or_group" gorm:"default:false"` // true for OR groups, false for AND groups
	HasCFCondition bool      `json:"has_cf_condition" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	Course              Course                   `json:"-" gorm:"foreignKey:CourseID"`
	PrerequisiteCourses []PrerequisiteCourseLink `json:"prerequisite_courses" gorm:"foreignKey:GroupID"`
}

// PrerequisiteCourseLink represents individual courses within a prerequisite group
type PrerequisiteCourseLink struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GroupID              uuid.UUID `json:"group_id" gorm:"type:uuid;not null;index"`
	PrerequisiteCourseID uuid.UUID `json:"prerequisite_course_id" gorm:"type:uuid;not null"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relations
	Group              PrerequisiteGroup `json:"-" gorm:"foreignKey:GroupID"`
	PrerequisiteCourse Course            `json:"prerequisite_course" gorm:"foreignKey:PrerequisiteCourseID"`
}

type Course struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CurriculumID uuid.UUID      `json:"curriculum_id" gorm:"type:uuid;not null;uniqueIndex:idx_course_code_year_curriculum_category"`
	CategoryID   uuid.UUID      `json:"category_id" gorm:"type:uuid;not null;uniqueIndex:idx_course_code_year_curriculum_category"`
	Code         string         `json:"code" gorm:"not null;uniqueIndex:idx_course_code_year_curriculum_category"`
	Year         int            `json:"year" gorm:"not null;uniqueIndex:idx_course_code_year_curriculum_category"`
	NameEN       string         `json:"name_en" gorm:"not null"`
	NameTH       string         `json:"name_th" gorm:"not null"`
	Credits      int            `json:"credits" gorm:"not null"`
	HasCFOption  bool           `json:"has_cf_option" gorm:"default:false"` // Indicates if this course allows C.F. exemption
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Curriculum         Curriculum          `json:"-" gorm:"foreignKey:CurriculumID"`
	Category           Category            `json:"-" gorm:"foreignKey:CategoryID"`
	PrerequisiteGroups []PrerequisiteGroup `json:"prerequisite_groups" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	CorequisiteGroups  []PrerequisiteGroup `json:"corequisite_groups" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}

// MarshalJSON customizes JSON marshaling to ensure PrerequisiteGroups and CorequisiteGroups are never null
func (c Course) MarshalJSON() ([]byte, error) {
	type Alias Course

	// Initialize empty slices if nil
	if c.PrerequisiteGroups == nil {
		c.PrerequisiteGroups = []PrerequisiteGroup{}
	}
	if c.CorequisiteGroups == nil {
		c.CorequisiteGroups = []PrerequisiteGroup{}
	}

	return json.Marshal((Alias)(c))
}
