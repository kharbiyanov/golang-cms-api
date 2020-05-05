package main

import (
	"github.com/graphql-go/graphql"
)

var TranslationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "FileTranslation",
		Fields: graphql.Fields{
			"element_id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"lang": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
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
			"translations": &graphql.Field{
				Type: graphql.NewList(TranslationType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return GetTranslationsInFile(params)
				},
			},
		},
	},
)
