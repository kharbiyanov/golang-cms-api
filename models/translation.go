package models

import (
	"time"
)

type Translation struct {
	ID          int        `gorm:"type:bigserial; primary_key"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
	ElementType string     `gorm:"not null"`
	ElementID   int        `gorm:"not null"`
	GroupID     int        `gorm:"auto_increment; not null"`
	Lang        string     `gorm:"not null; type:varchar(2);"`
}
