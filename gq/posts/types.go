package main

import (
	"cms-api/models"
	"github.com/graphql-go/graphql"
	"strings"
)

var MetaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PostMeta",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"key": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func GetPostType(postConfig models.PostConfig) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: strings.Title(postConfig.Slug),
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
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
				"content": &graphql.Field{
					Type: graphql.String,
				},
				"excerpt": &graphql.Field{
					Type: graphql.String,
				},
				"status": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"slug": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"meta": &graphql.Field{
					Type: graphql.NewList(MetaType),
					Args: graphql.FieldConfigArgument{
						"keys": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return GetMetaInPost(params)
					},
				},
			},
		},
	)
}
