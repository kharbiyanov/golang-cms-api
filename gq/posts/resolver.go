package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

func GetPost(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	id, _ := params.Args["id"].(int)

	var post = models.Post{}
	post.ID = id

	if err := utils.DB.Where(&models.Post{Type: postConfig.Type}).First(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetPosts(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	var query = PostQuery{
		tx:         utils.DB,
		params:     params,
		postConfig: postConfig,
	}

	if err := query.init(); err != nil {
		return nil, err
	}

	return query.results, nil
}

func CreatePost(params graphql.ResolveParams, postConfig models.PostConfig) (interface{}, error) {
	status := "publish"

	if pStatus, ok := params.Args["status"].(string); ok {
		switch pStatus {
		case "publish", "draft", "pending":
			status = pStatus
		}
	}

	post := &models.Post{
		Title:  params.Args["title"].(string),
		Status: status,
		Slug:   params.Args["slug"].(string),
		Type:   postConfig.Type,
	}

	if content, ok := params.Args["content"].(string); ok {
		post.Content = content
	}

	if excerpt, ok := params.Args["excerpt"].(string); ok {
		post.Excerpt = excerpt
	}

	lang, _ := params.Args["lang"].(string)

	if !utils.DB.Where(&models.Post{Slug: post.Slug, Type: postConfig.Type}).First(&post).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.PostSlugExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Create(post).Scan(post).Error; err != nil {
		return nil, err
	}

	translation := models.Translation{
		ElementType: fmt.Sprintf("post_%s", postConfig.Type),
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
		switch status {
		case "publish", "draft", "pending", "trash":
			fields["status"] = status
		}
	}
	if slug, ok := params.Args["slug"].(string); ok {
		fields["slug"] = slug
		if !utils.DB.Where(&models.Post{Type: postConfig.Type, Slug: slug}).Not(&models.Post{ID: id}).First(&post).RecordNotFound() {
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
	id, _ := params.Args["id"].(int)

	var post = models.Post{}
	post.ID = id
	post.Type = postConfig.Type

	if err := utils.DB.Delete(&post).Error; err != nil {
		return nil, err
	}
	return nil, utils.DB.Delete(&models.PostMeta{}, &models.PostMeta{PostID: id}).Error
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

func GetTermsInPost(params graphql.ResolveParams) (interface{}, error) {
	taxonomies, taxonomiesExist := params.Args["taxonomies"].([]interface{})
	post, postExist := params.Source.(models.Post)

	if !postExist {
		return nil, nil
	}

	tx := utils.DB

	if taxonomiesExist {
		tx = tx.Where("terms.taxonomy in(?)", taxonomies)
	}

	tx = tx.Table("terms").Select("terms.*").Joins("LEFT JOIN term_relationships r ON terms.id = r.term_id").Where("r.post_id = ?", post.ID)

	if err := tx.Find(&post.Terms).Error; err != nil {
		return nil, err
	}

	return post.Terms, nil
}

func GetMeta(params graphql.ResolveParams) (interface{}, error) {
	postId, _ := params.Args["post_id"].(int)
	keys, keysExist := params.Args["keys"].([]interface{})

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
	postId, _ := params.Args["post_id"].(int)
	key, _ := params.Args["key"].(string)
	value, _ := params.Args["value"].(string)

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
	postId, _ := params.Args["post_id"].(int)
	key, _ := params.Args["key"].(string)

	return nil, utils.DB.Delete(&models.PostMeta{}, &models.PostMeta{PostID: postId, Key: key}).Error
}
