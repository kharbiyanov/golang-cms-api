package models

type (
	Role struct {
		Name   string
		Access []Access
	}
	Access struct {
		Object string
		Action string
	}
)
