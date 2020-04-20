package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
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

func DeleteMenu(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	menu := models.Menu{}

	if utils.DB.First(&menu, id).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.MenuNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Delete(menu).Error; err != nil {
		return nil, err
	}

	if err := utils.DB.Delete(&models.MenuItem{}, "menu_id = ?", id).Error; err != nil {
		return nil, err
	}

	return menu, nil
}

func UpdateMenu(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	name, _ := params.Args["name"].(string)

	var menu models.Menu

	if err := utils.DB.Model(&menu).Where("id = ?", id).Update("name", name).Find(&menu).Error; err != nil {
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

func DeleteMenuItem(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	menuItem := models.MenuItem{}

	if utils.DB.First(&menuItem, id).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.MenuItemNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Delete(menuItem).Error; err != nil {
		return nil, err
	}

	return menuItem, nil
}

func UpdateMenuItem(params graphql.ResolveParams) (interface{}, error) {
	menuID, _ := params.Args["menu_id"].(int)
	fields := make(map[string]interface{})

	var menuItem = models.MenuItem{
		ID: menuID,
	}

	if title, ok := params.Args["title"].(string); ok {
		fields["title"] = title
	}

	if itemType, ok := params.Args["type"].(string); ok {
		fields["type"] = itemType
	}

	if object, ok := params.Args["object"].(string); ok {
		fields["object"] = object
	}

	if objectID, ok := params.Args["object_id"].(int); ok {
		fields["object_id"] = objectID
	}

	if url, ok := params.Args["url"].(string); ok {
		fields["url"] = url
	}

	if parent, ok := params.Args["parent"].(int); ok {
		fields["parent"] = parent
	}

	if order, ok := params.Args["order"].(int); ok {
		fields["order"] = order
	}

	if target, ok := params.Args["target"].(string); ok {
		fields["target"] = target
	}

	if classes, ok := params.Args["classes"].(string); ok {
		fields["classes"] = classes
	}

	if err := utils.DB.Model(&menuItem).Updates(fields).Scan(&menuItem).Error; err != nil {
		return nil, err
	}

	return menuItem, nil
}

func GetMenuList() (interface{}, error) {
	var menuList []models.Menu

	if err := utils.DB.Find(&menuList).Error; err != nil {
		return nil, err
	}
	return menuList, nil
}

func GetMenuItemList(params graphql.ResolveParams) (interface{}, error) {
	menuID, _ := params.Args["menu_id"].(int)

	var items []models.MenuItem

	rows, err := utils.DB.Table("menu_items i").
		Select("i.id, i.created_at, i.updated_at, i.menu_id, i.author_id, i.title, i.type, i.object, i.object_id, i.url, i.parent, i.order, i.target, i.classes").
		Where("i.deleted_at IS NULL AND i.menu_id = ?", menuID).
		Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.MenuItem

		if err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.MenuID,
			&item.AuthorID,
			&item.Title,
			&item.Type,
			&item.Object,
			&item.ObjectID,
			&item.Url,
			&item.Parent,
			&item.Order,
			&item.Target,
			&item.Classes,
		); err != nil {
			return nil, err
		}
		if item.Object != "" && item.ObjectID > 0 {
			switch item.Type {
			case "post":
				var post = models.Post{}
				post.ID = item.ObjectID

				if err := utils.DB.Where(&models.Post{Type: item.Object}).First(&post).Error; err == nil {
					item.Url = fmt.Sprintf("http://%s", post.Slug)
				}
			case "taxonomy":
				var term = models.Term{}
				term.ID = item.ObjectID

				if err := utils.DB.Where(&models.Term{Taxonomy: item.Object}).First(&term).Error; err == nil {
					item.Url = fmt.Sprintf("http://%s", term.Slug)
				}
			}
		}
		items = append(items, item)
	}

	return items, nil
}
