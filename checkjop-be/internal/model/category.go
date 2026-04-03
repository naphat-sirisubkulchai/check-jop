package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CurriculumID uuid.UUID      `json:"curriculum_id" gorm:"type:uuid;not null;uniqueIndex:idx_category_names_curriculum"`
	NameTH       string         `json:"name_th" gorm:"not null;uniqueIndex:idx_category_names_curriculum"`
	NameEN       string         `json:"name_en" gorm:"not null;uniqueIndex:idx_category_names_curriculum"`
	MinCredits   int            `json:"min_credits" gorm:"not null"` // จำนวนหน่วยกิตขั้นต่ำในหมวดนี้
	SortOrder    int            `json:"sort_order" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Curriculum Curriculum `json:"-" gorm:"foreignKey:CurriculumID"`
	Courses    []Course   `json:"courses" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}
