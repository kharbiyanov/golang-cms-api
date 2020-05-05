package main

import (
	"cms-api/models"
	"github.com/graphql-go/graphql"
	"strings"
)

var TranslationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TermTranslation",
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

func GetTaxonomyType(taxonomyConfig models.TaxonomyConfig) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: strings.Title(taxonomyConfig.Taxonomy),
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
						return GetMetaInTerm(params)
					},
				},
				"translations": &graphql.Field{
					Type: graphql.NewList(TranslationType),
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return GetTranslationsInTerm(params)
					},
				},
			},
		},
	)
}
