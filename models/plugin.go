package models

import "github.com/graphql-go/graphql"

type Plugin struct {
	ID             string
	Name           string
	Enable         bool
	QueryFields    graphql.Fields
	MutationFields graphql.Fields
}
