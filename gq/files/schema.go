package main

import (
	"cms-api/models"
	"git.osg.uz/kharbiyanov/graphql-multipart-middleware"
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
			return UploadFile(params)
		},
	}
}
