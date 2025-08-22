package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uint      `gorm:"primaryKey;not null"               json:"-"`
	Uuid      uuid.UUID `gorm:"type:varchar(36);not null"         json:"-"`
	Name      string    `gorm:"unique;type:varchar(255);not null" json:"name"`
	Slug      string    `gorm:"unique;type:varchar(255);not null" json:"slug"`
	CreatedAt time.Time `gorm:"not null"                          json:"-"`
	UpdatedAt time.Time `gorm:"not null"                          json:"-"`
	Reviews   []*Review `gorm:"many2many:categories_reviews"      json:"-"`
}

func (category *Category) BeforeCreate(transaction *gorm.DB) (err error) {
	category.Uuid = uuid.New()
	return
}
