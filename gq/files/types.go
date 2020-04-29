package main

import (
	"github.com/graphql-go/graphql"
)

var FileType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "File",
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
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"mime_type": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"file": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
