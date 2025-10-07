package models

import (
	"time"

	"github.com/google/uuid"
)

type ReadingList struct {
	ID        uint      `gorm:"primaryKey;not null"        json:"-"`
	Uuid      uuid.UUID `gorm:"type:varchar(36);not null"  json:"uuid"`
	Status    uint      `gorm:"not null;default:0"         json:"status"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Author    string    `gorm:"type:varchar(75);not null"  json:"author"`
	Image     string    `gorm:"type:varchar(255);not null" json:"image"`
	CreatedAt time.Time `gorm:"not null"                   json:"-"`
	UpdatedAt time.Time `gorm:"not null"                   json:"-"`
}

func (ReadingList) TableName() string {
	return "reading_list"
}
