package server

import (
	"cms-api/config"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"path/filepath"
	"plugin"
)

var (
	queryFields    = graphql.Fields{}
	mutationFields = graphql.Fields{}
)

func GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := handler.New(&handler.Config{
			Schema:     getSchema(),
			Pretty:     true,
			Playground: config.Get().Debug,
			RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
				return map[string]interface{}{
					"context": c,
				}
			},
		})
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func getSchema() *graphql.Schema {
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   "Query",
					Fields: queryFields,
				},
			),
			Mutation: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   "Mutation",
					Fields: mutationFields,
				},
			),
		},
	)

	return &schema
}

func SetupPlugins() {
	plugins, err := filepath.Glob("plugins/*.so")

	if err != nil {
		panic(err)
	}

	for _, filename := range plugins {
		query, mutation := getPluginFields(filename)

		for key, val := range query {
			queryFields[key] = val
		}
		for key, val := range mutation {
			mutationFields[key] = val
		}

		if config.Get().Debug {
			log.Printf("Plugin '%s' loaded", filename)
		}
	}
}

func getPluginFields(filename string) (graphql.Fields, graphql.Fields) {
	p, err := plugin.Open(filename)
	if err != nil {
		panic(err)
	}
	symbol, err := p.Lookup("GetSchemaFields")
	if err != nil {
		panic(err)
	}

	getSchemaFields, ok := symbol.(func() (graphql.Fields, graphql.Fields))

	if !ok {
		panic("Plugin has no 'GetSchemaFields() (graphql.Fields, graphql.Fields)' function")
	}

	return getSchemaFields()
}
