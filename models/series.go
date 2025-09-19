package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Series struct {
	ID          uint      `gorm:"primaryKey;not null"               json:"-"`
	Uuid        uuid.UUID `gorm:"type:varchar(36);not null"         json:"-"`
	Name        string    `gorm:"unique;type:varchar(255);not null" json:"name"`
	Slug        string    `gorm:"unique;type:varchar(255);not null" json:"slug"`
	Author      string    `gorm:"type:varchar(75)"                  json:"author"`
	Description *string   `gorm:"type:text"                         json:"description"`
	CreatedAt   time.Time `gorm:"not null"                          json:"-"`
	UpdatedAt   time.Time `gorm:"not null"                          json:"-"`
	Reviews     []Review  `                                         json:"-"`

	Images []string `gorm:"-" json:"images"`
}

func (series *Series) BeforeCreate(transaction *gorm.DB) (err error) {
	series.Uuid = uuid.New()
	return
}
