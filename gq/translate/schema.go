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
	queryFields["langList"] = &graphql.Field{
		Type:        graphql.NewList(LangType),
		Description: "Get lang list.",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetLangList()
		},
	}
}

func setupMutation() {
	mutationFields["langCreate"] = &graphql.Field{
		Type:        LangType,
		Description: "Create new lang.",
		Args: graphql.FieldConfigArgument{
			"full_name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"code": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"flag": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"default": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "lang", "create"); err != nil {
				return nil, err
			}

			return CreateLang(params)
		},
	}
	mutationFields["langDelete"] = &graphql.Field{
		Type:        LangType,
		Description: "Delete lang by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "lang", "delete"); err != nil {
				return nil, err
			}
			return DeleteLang(params)
		},
	}
	mutationFields["langUpdate"] = &graphql.Field{
		Type:        LangType,
		Description: "Update lang by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"full_name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"code": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"flag": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"default": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "lang", "update"); err != nil {
				return nil, err
			}

			return UpdateLang(params)
		},
	}
}
