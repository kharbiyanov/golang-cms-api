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
	mutationFields["langCreate"] = &graphql.Field{
		Type:        LangType,
		Description: "Update meta by post_id and meta_key. If the meta field for the post does not exist, it will be added.",
		Args: graphql.FieldConfigArgument{
			"fullName": &graphql.ArgumentConfig{
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
}
