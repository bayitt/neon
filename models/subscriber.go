package models

import "time"

type Subscriber struct {
	ID        uint      `gorm:"primaryKey;not null"               json:"-"`
	Email     string    `gorm:"unique;type:varchar(100);not null" json:"email"`
	CreatedAt time.Time `gorm:"not null"                          json:"-"`
}
