package models

import "time"

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
