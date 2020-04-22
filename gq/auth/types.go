package main

import (
	"github.com/graphql-go/graphql"
)

var TokenType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Token",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"expires_at": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
