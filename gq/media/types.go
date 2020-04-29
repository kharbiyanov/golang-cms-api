package main

import (
	"github.com/graphql-go/graphql"
)

var MediaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Media",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"created_at": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"updated_at": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
