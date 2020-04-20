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
	queryFields["menuList"] = &graphql.Field{
		Type:        graphql.NewList(MenuType),
		Description: "Get menu list.",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetMenuList()
		},
	}
	queryFields["menuItemList"] = &graphql.Field{
		Type:        graphql.NewList(MenuItemType),
		Description: "Get menu item list.",
		Args: graphql.FieldConfigArgument{
			"menu_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetMenuItemList(params)
		},
	}
}

func setupMutation() {
	mutationFields["menuCreate"] = &graphql.Field{
		Type:        MenuType,
		Description: "Create new menu.",
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "menu", "create"); err != nil {
				return nil, err
			}

			return CreateMenu(params)
		},
	}
	mutationFields["menuDelete"] = &graphql.Field{
		Type:        MenuType,
		Description: "Delete menu.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "menu", "delete"); err != nil {
				return nil, err
			}

			return DeleteMenu(params)
		},
	}
	mutationFields["menuUpdate"] = &graphql.Field{
		Type:        MenuType,
		Description: "Update menu.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "menu", "update"); err != nil {
				return nil, err
			}

			return UpdateMenu(params)
		},
	}
	mutationFields["menuItemCreate"] = &graphql.Field{
		Type:        MenuItemType,
		Description: "Create new menu item.",
		Args: graphql.FieldConfigArgument{
			"menu_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"type": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"object": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"object_id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"url": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"parent": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"order": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"target": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"classes": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "menuItem", "create"); err != nil {
				return nil, err
			}

			return CreateMenuItem(params)
		},
	}
	mutationFields["menuItemDelete"] = &graphql.Field{
		Type:        MenuItemType,
		Description: "Delete menu item.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "menuItem", "delete"); err != nil {
				return nil, err
			}

			return DeleteMenuItem(params)
		},
	}
	mutationFields["menuItemUpdate"] = &graphql.Field{
		Type:        MenuItemType,
		Description: "Update menu item.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"type": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"object": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"object_id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"url": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"parent": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"order": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"target": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"classes": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "menuItem", "update"); err != nil {
				return nil, err
			}

			return UpdateMenuItem(params)
		},
	}
}
