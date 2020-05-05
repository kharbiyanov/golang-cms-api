package main

import (
	"github.com/graphql-go/graphql"
)

var TranslationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MenuTranslation",
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

var MenuType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Menu",
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
			"translations": &graphql.Field{
				Type: graphql.NewList(TranslationType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return GetTranslationsInMenu(params)
				},
			},
		},
	},
)

var MenuItemType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MenuItem",
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
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"object": &graphql.Field{
				Type: graphql.String,
			},
			"object_id": &graphql.Field{
				Type: graphql.Int,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"parent": &graphql.Field{
				Type: graphql.Int,
			},
			"order": &graphql.Field{
				Type: graphql.Int,
			},
			"target": &graphql.Field{
				Type: graphql.String,
			},
			"classes": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
