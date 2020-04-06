package models

import (
	"time"
)

type Term struct {
	ID          int        `gorm:"type:bigserial; primary_key"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `sql:"index" json:"deletedAt"`
	Name        string     `gorm:"not null"`
	Taxonomy    string     `gorm:"not null"`
	Description string
	Slug        string `gorm:"not null; type:varchar(255);"`
	Parent      int
	Meta        []TermMeta
}

type TermMeta struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	TermID    int
	Key       string `gorm:"type:varchar(255)"`
	Value     string
}

type TermRelationship struct {
	PostID int
	TermID int
}

type TaxQuery struct {
	Taxonomy string
	Terms    []interface{}
	Operator string
	Field    string
}
