package main

import (
	"cms-api/models"
	"cms-api/utils"
	graphqlmultipart "git.osg.uz/kharbiyanov/graphql-multipart-middleware"
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
	mutationFields["fileCreate"] = &graphql.Field{
		Type:        FileType,
		Description: "Create new file.",
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"file": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphqlmultipart.Upload),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "file", "create"); err != nil {
				return nil, err
			}

			return UploadFile(params)
		},
	}
}
