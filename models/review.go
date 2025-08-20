package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID         uint        `gorm:"primaryKey;not null"                           json:"-"`
	Uuid       uuid.UUID   `gorm:"type:varchar(36);not null"                     json:"uuid"`
	SeriesID   *uint       `                                                     json:"-"`
	Series     *Series     `gorm:"constraint:OnUpdate:CASCADE,onDelete:RESTRICT" json:"-"`
	Title      string      `gorm:"type:varchar(255);not null"                    json:"title"`
	Slug       string      `gorm:"unique;type:varchar(255);not null"             json:"slug"`
	Author     string      `gorm:"type:varchar(75)"                              json:"author"`
	Image      *string     `gorm:"type:varchar(255)"                             json:"image"`
	Status     uint        `gorm:"not null;default:0"                            json:"status"`
	Content    string      `gorm:"type:text;not null"                            json:"content"`
	CreatedAt  time.Time   `gorm:"not null"                                      json:"created_at"`
	UpdatedAt  time.Time   `gorm:"not null"                                      json:"-"`
	Categories []*Category `gorm:"many2many:categories_reviews"                  json:"-"`
}

func (review *Review) BeforeCreate(transaction *gorm.DB) (err error) {
	review.Uuid = uuid.New()
	return
}
