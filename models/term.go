package models

import (
	"time"
)

type TaxonomyConfig struct {
	Taxonomy  string `json:"taxonomy"`
	SingleUrl string `json:"singleUrl"`
}

type Term struct {
	ID           int        `gorm:"type:bigserial; primary_key"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
	Name         string     `gorm:"not null"`
	Taxonomy     string     `gorm:"not null"`
	Description  string
	Slug         string `gorm:"not null; type:varchar(255);"`
	Parent       int
	Meta         []TermMeta
	Translations []Translation
}

type TermMeta struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
