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
	queryFields["roleList"] = &graphql.Field{
		Type:        graphql.NewList(RoleType),
		Description: "Get role list.",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "role", "list"); err != nil {
				return nil, err
			}
			return GetRoleList()
		},
	}
}

func setupMutation() {

}
