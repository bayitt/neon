package models

import "time"

type CategoryReview struct {
	CategoryID uint      `gorm:"primaryKey;not null"`
	ReviewID   uint      `gorm:"primaryKey;not null"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (CategoryReview) TableName() string {
	return "categories_reviews"
}
