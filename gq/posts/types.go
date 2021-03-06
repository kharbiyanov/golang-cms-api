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

var TranslationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PostTranslation",
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

var TermType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TermType",
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
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"taxonomy": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"slug": &graphql.Field{
				Type: graphql.String,
			},
			"parent": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var StateEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "State",
	Values: graphql.EnumValueConfigMap{
		"TRASH": &graphql.EnumValueConfig{
			Value: models.PostStateTrash,
		},
		"PUBLISH": &graphql.EnumValueConfig{
			Value: models.PostStatePublish,
		},
		"DRAFT": &graphql.EnumValueConfig{
			Value: models.PostStateDraft,
		},
	},
})

func GetPostType(postConfig models.PostConfig) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: strings.Title(postConfig.Type),
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
				"state": &graphql.Field{
					Type: graphql.NewNonNull(StateEnum),
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
				"terms": &graphql.Field{
					Type: graphql.NewList(TermType),
					Args: graphql.FieldConfigArgument{
						"taxonomies": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return GetTermsInPost(params)
					},
				},
				"translations": &graphql.Field{
					Type: graphql.NewList(TranslationType),
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return GetTranslationsInPost(params)
					},
				},
			},
		},
	)
}
