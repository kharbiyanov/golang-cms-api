package main

import (
	"cms-api/config"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

var (
	queryFields    = graphql.Fields{}
	mutationFields = graphql.Fields{}
)

func InitSchema(plugin *models.Plugin) {
	for _, taxonomyConfig := range config.Get().Taxonomies {
		taxonomyType := GetTaxonomyType(taxonomyConfig)
		setupTaxonomiesQuery(taxonomyType, taxonomyConfig)
		setupTaxonomiesMutation(taxonomyType, taxonomyConfig)
	}
	setupMutations()

	plugin.QueryFields = queryFields
	plugin.MutationFields = mutationFields
}

func setupTaxonomiesQuery(taxonomyType *graphql.Object, taxonomyConfig models.TaxonomyConfig) {
	queryFields[fmt.Sprintf("%sList", taxonomyConfig.Taxonomy)] = &graphql.Field{
		Type:        graphql.NewList(taxonomyType),
		Description: fmt.Sprintf("Get %s list.", taxonomyConfig.Taxonomy),
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
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
			return GetTerms(params, taxonomyConfig)
		},
	}
}

func setupTaxonomiesMutation(taxonomyType *graphql.Object, taxonomyConfig models.TaxonomyConfig) {
	mutationFields[fmt.Sprintf("%sCreate", taxonomyConfig.Taxonomy)] = &graphql.Field{
		Type:        taxonomyType,
		Description: fmt.Sprintf("Create new %s.", taxonomyConfig.Taxonomy),
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
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
			if err := utils.ValidateUser(params, taxonomyConfig.Taxonomy, "create"); err != nil {
				return nil, err
			}

			return CreateTerm(params, taxonomyConfig)
		},
	}
	mutationFields[fmt.Sprintf("%sUpdate", taxonomyConfig.Taxonomy)] = &graphql.Field{
		Type:        taxonomyType,
		Description: fmt.Sprintf("Update %s by id.", taxonomyConfig.Taxonomy),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
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
			if err := utils.ValidateUser(params, taxonomyConfig.Taxonomy, "update"); err != nil {
				return nil, err
			}

			return UpdateTerm(params, taxonomyConfig)
		},
	}
	mutationFields[fmt.Sprintf("%sDelete", taxonomyConfig.Taxonomy)] = &graphql.Field{
		Type:        taxonomyType,
		Description: fmt.Sprintf("Delete %s by id.", taxonomyConfig.Taxonomy),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, taxonomyConfig.Taxonomy, "delete"); err != nil {
				return nil, err
			}
			return DeleteTerm(params, taxonomyConfig)
		},
	}
}

func setupMutations() {
	mutationFields["termMetaUpdate"] = &graphql.Field{
		Type:        MetaType,
		Description: "Update term meta by term_id and meta_key. If the meta field for the post does not exist, it will be added.",
		Args: graphql.FieldConfigArgument{
			"term_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"value": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "termMeta", "update"); err != nil {
				return nil, err
			}

			return UpdateMeta(params)
		},
	}
	mutationFields["termMetaDelete"] = &graphql.Field{
		Type:        MetaType,
		Description: "Delete %s by term_id and meta_key.",
		Args: graphql.FieldConfigArgument{
			"term_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "termMeta", "delete"); err != nil {
				return nil, err
			}

			return DeleteMeta(params)
		},
	}
	mutationFields["setTerms"] = &graphql.Field{
		Type:        MetaType,
		Description: "Set terms by post id.",
		Args: graphql.FieldConfigArgument{
			"post_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"terms": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.NewList(graphql.Int)),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "term", "set"); err != nil {
				return nil, err
			}

			return SetTerms(params)
		},
	}
}
