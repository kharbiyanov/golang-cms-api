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
			if err := utils.ValidateUser(params, "terms", "create"); err != nil {
				return nil, err
			}

			return CreateTerm(params)
		},
	}
}
