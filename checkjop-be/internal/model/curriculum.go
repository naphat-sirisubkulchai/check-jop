package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// หลักสูตร
type Curriculum struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	NameTH          string         `json:"name_th" gorm:"not null;uniqueIndex:idx_curriculum_names"`
	NameEN          string         `json:"name_en" gorm:"not null;uniqueIndex:idx_curriculum_names"`
	Year            int            `json:"year" gorm:"not null"`              // ปีหลักสูตร
	MinTotalCredits int            `json:"min_total_credits" gorm:"not null"` // จำนวนหน่วยกิตรวมขั้นต่ำ
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Categories []Category `json:"categories" gorm:"foreignKey:CurriculumID;constraint:OnDelete:CASCADE"` // one to many
	Courses    []Course   `json:"courses" gorm:"foreignKey:CurriculumID;constraint:OnDelete:CASCADE"`    // one to many
}
