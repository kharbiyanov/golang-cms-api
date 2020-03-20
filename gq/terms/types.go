package main

import (
	"github.com/graphql-go/graphql"
)

var MetaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TermMeta",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"key": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"value": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var TermType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Term",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"taxonomy": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"slug": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"parent": &graphql.Field{
				Type: graphql.Int,
			},
			"meta": &graphql.Field{
				Type: graphql.NewList(MetaType),
				Args: graphql.FieldConfigArgument{
					"keys": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	},
)
