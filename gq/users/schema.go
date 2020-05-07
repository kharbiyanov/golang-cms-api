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

	plugin.QueryFields = queryFields
}

func setupQuery() {
	queryFields["userGet"] = &graphql.Field{
		Type:        UserType,
		Description: "Get user by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "user", "get"); err != nil {
				return nil, err
			}
			return GetUser(params)
		},
	}
}
