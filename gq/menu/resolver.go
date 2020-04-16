package main

import (
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

func CreateMenu(params graphql.ResolveParams) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)

	menu := &models.Menu{
		Name: params.Args["name"].(string),
	}

	if err := utils.DB.Create(&menu).Scan(&menu).Error; err != nil {
		return nil, err
	}

	translation := models.Translation{
		ElementType: "menu",
		ElementID:   menu.ID,
		Lang:        lang,
	}

	if err := utils.DB.Create(&translation).Scan(&translation).Error; err != nil {
		return nil, err
	}

	return menu, nil
}

func CreateMenuItem(params graphql.ResolveParams) (interface{}, error) {
	menuID, _ := params.Args["menu_id"].(int)
	itemType, _ := params.Args["type"].(string)

	menuItem := &models.MenuItem{
		MenuID: menuID,
		Type:   itemType,
	}

	if title, ok := params.Args["title"].(string); ok {
		menuItem.Title = title
	}

	if object, ok := params.Args["object"].(string); ok {
		menuItem.Object = object
	}

	if objectID, ok := params.Args["object_id"].(int); ok {
		menuItem.ObjectID = objectID
	}

	if url, ok := params.Args["url"].(string); ok {
		menuItem.Url = url
	}

	if parent, ok := params.Args["parent"].(int); ok {
		menuItem.Parent = parent
	}

	if order, ok := params.Args["order"].(int); ok {
		menuItem.Order = order
	}

	if target, ok := params.Args["target"].(string); ok {
		menuItem.Target = target
	}

	if classes, ok := params.Args["classes"].(string); ok {
		menuItem.Classes = classes
	}

	if err := utils.DB.Create(&menuItem).Scan(&menuItem).Error; err != nil {
		return nil, err
	}

	return menuItem, nil
}
