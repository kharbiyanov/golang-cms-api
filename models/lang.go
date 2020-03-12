package models

import (
	"time"
)

type Lang struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	FullName  string     `gorm:"not null" json:"full_name"`
	Code      string     `gorm:"not null; type:varchar(2); unique;"`
	Flag      string
	Default   bool
}

func (l Lang) TableName() string {
	return "lang"
}
