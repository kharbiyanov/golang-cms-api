package models

import (
	"time"
)

type (
	User struct {
		ID             int        `gorm:"type:bigserial; primary_key"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at"`
		DeletedAt      *time.Time `sql:"index" json:"deleted_at"`
		UserName       string     `json:"userName"`
		Password       string     `json:"-"`
		LastName       string     `json:"lastName"`
		FirstName      string     `json:"firstName"`
		MiddleName     string     `json:"middleName"`
		Avatar         string     `json:"avatar"`
		Phone          int        `json:"phone"`
		Email          string     `json:"email"`
		ActivationCode string     `gorm:"type:varchar(36)" json:"activation_code"`
		State          int
	}
)
