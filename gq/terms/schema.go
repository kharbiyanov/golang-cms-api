package main

import (
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

var (
	queryFields    = graphql.Fields{}
	mutationFields = graphql.Fields{}
)

func InitSchema(plugin *models.Plugin) {
	setupQuery()
	setupMutation()

	plugin.QueryFields = queryFields
	plugin.MutationFields = mutationFields
}

func setupQuery() {
	queryFields["termGetList"] = &graphql.Field{
		Type:        graphql.NewList(TermType),
		Description: "Get term list.",
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"taxonomy": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"parent": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"first": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetTerms(params)
		},
	}
}

func setupMutation() {
	mutationFields["termCreate"] = &graphql.Field{
		Type:        TermType,
		Description: "Create new lang.",
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"taxonomy": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"slug": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"parent": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "term", "create"); err != nil {
				return nil, err
			}

			return CreateTerm(params)
		},
	}
	mutationFields["termUpdate"] = &graphql.Field{
		Type:        TermType,
		Description: "Update term by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"taxonomy": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"slug": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"parent": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "term", "update"); err != nil {
				return nil, err
			}

			return UpdateTerm(params)
		},
	}
	mutationFields["termDelete"] = &graphql.Field{
		Type:        TermType,
		Description: "Delete term by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "term", "delete"); err != nil {
				return nil, err
			}
			return DeleteTerm(params)
		},
	}
}
