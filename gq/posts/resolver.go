package main

import (
	"cms-api/config"
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

func GetPost(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	id, idExist := params.Args["id"].(int)
	if !idExist {
		return nil, nil
	}

	var post = models.Post{}
	post.ID = id

	if err := utils.DB.Where(&models.Post{Type: postConfig.Slug}).First(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetPosts(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)
	first, firstExist := params.Args["first"].(int)
	offset, offsetExist := params.Args["offset"].(int)

	var posts []models.Post

	tx := utils.DB.Table("posts").
		Select("posts.*").
		Joins("left join translations on translations.element_id = posts.id").
		Where("posts.type = ? and translations.lang = ? and translations.element_type = ?", postConfig.Slug, lang, fmt.Sprintf("post_%s", postConfig.Slug))

	if firstExist {
		tx = tx.Limit(first)
	} else {
		tx = tx.Limit(config.Get().DefaultPostsLimit)
	}
	if offsetExist {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func CreatePost(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	post := &models.Post{
		Title:   params.Args["title"].(string),
		Content: params.Args["content"].(string),
		Excerpt: params.Args["excerpt"].(string),
		Status:  params.Args["status"].(int),
		Slug:    params.Args["slug"].(string),
		Type:    postConfig.Slug,
	}

	lang, _ := params.Args["lang"].(string)

	if !utils.DB.Where(&models.Post{Slug: post.Slug, Type: postConfig.Slug}).First(&post).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.PostSlugExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Create(post).Scan(post).Error; err != nil {
		return nil, err
	}

	translation := models.Translation{
		ElementType: fmt.Sprintf("post_%s", postConfig.Slug),
		ElementID:   post.ID,
		Lang:        lang,
	}

	if err := utils.DB.Create(&translation).Scan(&translation).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func UpdatePost(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	fields := make(map[string]interface{})

	var post models.Post

	if title, ok := params.Args["title"].(string); ok {
		fields["title"] = title
	}
	if content, ok := params.Args["content"].(string); ok {
		fields["content"] = content
	}
	if excerpt, ok := params.Args["excerpt"].(string); ok {
		fields["excerpt"] = excerpt
	}
	if status, ok := params.Args["status"].(string); ok {
		fields["status"] = status
	}
	if slug, ok := params.Args["slug"].(string); ok {
		fields["slug"] = slug
		if !utils.DB.Where(&models.Post{Type: postConfig.Slug, Slug: slug}).Not(&models.Post{ID: id}).First(&post).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.PostSlugExistMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
	}

	post.ID = id

	if err := utils.DB.Model(&post).Updates(fields).Scan(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func DeletePost(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	id, idExist := params.Args["id"].(int)
	if !idExist {
		return nil, nil
	}

	var post = models.Post{}
	post.ID = id
	post.Type = postConfig.Slug

	if err := utils.DB.Delete(&post).Error; err != nil {
		return nil, err
	}
	if err := utils.DB.Delete(&models.PostMeta{}, &models.PostMeta{PostID: id}).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func GetMetaInPost(params graphql.ResolveParams) (interface{}, error) {
	keys, keysExist := params.Args["keys"].([]interface{})
	post, postExist := params.Source.(models.Post)

	if !postExist {
		return nil, nil
	}

	tx := utils.DB

	if keysExist {
		tx = tx.Where("key in(?)", keys)
	}

	if err := tx.Model(&post).Association("Meta").Find(&post.Meta).Error; err != nil {
		return nil, err
	}

	return post.Meta, nil
}

func GetMeta(params graphql.ResolveParams) (interface{}, error) {
	postId, postIdExist := params.Args["post_id"].(int)
	keys, keysExist := params.Args["keys"].([]interface{})

	if !postIdExist {
		return nil, nil
	}

	var meta []models.PostMeta

	tx := utils.DB.Where(&models.PostMeta{PostID: postId})

	if keysExist {
		tx = tx.Where("key in(?)", keys)
	}

	if err := tx.Find(&meta).Error; err != nil {
		return nil, err
	}
	return meta, nil
}

func UpdateMeta(params graphql.ResolveParams) (interface{}, error) {
	postId, postIdExist := params.Args["post_id"].(int)
	key, keyExist := params.Args["key"].(string)
	value, valueExist := params.Args["value"].(string)

	if !postIdExist || !keyExist || !valueExist {
		return nil, nil
	}

	meta := models.PostMeta{
		PostID: postId,
		Key:    key,
		Value:  value,
	}

	if err := utils.DB.Model(&meta).Where(&models.PostMeta{PostID: postId, Key: key}).Update(&meta).Scan(&meta).Error; err != nil {
		if err := utils.DB.Save(&meta).First(&meta).Error; err != nil {
			return nil, err
		}
	}

	return meta, nil
}

func DeleteMeta(params graphql.ResolveParams) (interface{}, error) {
	postId, postIdExist := params.Args["post_id"].(int)
	key, keyExist := params.Args["key"].(string)

	if !postIdExist || !keyExist {
		return nil, nil
	}

	if err := utils.DB.Delete(&models.PostMeta{}, &models.PostMeta{PostID: postId, Key: key}).Error; err != nil {
		return nil, err
	}
	return nil, nil
}
