package models

import (
	"time"
)

type Translation struct {
	ID          int        `gorm:"type:bigserial; primary_key"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `sql:"index" json:"deletedAt"`
	ElementType string     `gorm:"not null"`
	ElementID   int        `gorm:"not null"`
	GroupID     int        `gorm:"auto_increment; not null"`
	Lang        string     `gorm:"not null; type:varchar(2);"`
}
