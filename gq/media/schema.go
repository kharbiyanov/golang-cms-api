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
	mutationFields["mediaUpload"] = &graphql.Field{
		Type:        MediaType,
		Description: "Upload new media.",
		Args: graphql.FieldConfigArgument{
			"file": &graphql.ArgumentConfig{
				Type: graphqlmultipart.Upload,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return UploadMedia(params)
		},
	}
}
