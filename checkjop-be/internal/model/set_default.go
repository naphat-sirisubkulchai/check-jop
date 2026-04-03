package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SetDefault struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CurriculumID uuid.UUID      `json:"curriculum_id" gorm:"type:uuid;not null;uniqueIndex:idx_set_default_unique"`
	CourseID     uuid.UUID      `json:"course_id" gorm:"type:uuid;not null;uniqueIndex:idx_set_default_unique"`
	Year         int            `json:"year" gorm:"not null"`
	Semester     int            `json:"semester" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Curriculum Curriculum `json:"curriculum" gorm:"foreignKey:CurriculumID"`
	Course     Course     `json:"course" gorm:"foreignKey:CourseID"`
}
