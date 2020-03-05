package posts

import (
	"cms-api/errors"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

func GetPost(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	id, idExist := params.Args["id"].(int)
	if ! idExist {
		return nil, nil
	}

	var post = Post{}
	post.ID = id

	if err := utils.DB.Where(&Post{Type: postConfig.Slug}).First(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetPosts(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	first, firstExist := params.Args["first"].(int)
	offset, offsetExist := params.Args["offset"].(int)

	var posts []Post

	tx := utils.DB.Where(&Post{Type: postConfig.Slug})
	if firstExist {
		tx = tx.Limit(first)
	} else {
		tx = tx.Limit(utils.Config.DefaultPostsLimit)
	}
	if offsetExist {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func CreatePost(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	post := &Post{
		Title:   params.Args["title"].(string),
		Content: params.Args["content"].(string),
		Excerpt: params.Args["excerpt"].(string),
		Status:  params.Args["status"].(int),
		Slug:    params.Args["slug"].(string),
		Type:    postConfig.Slug,
	}

	if ! utils.DB.Where(&Post{Slug: post.Slug, Type: postConfig.Slug}).First(&post).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.PostSlugExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Create(post).Scan(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func UpdatePost(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	fields := make(map[string]interface{})

	var post Post

	if title, titleExist := params.Args["title"].(string); titleExist {
		fields["title"] = title
	}
	if content, contentExist := params.Args["content"].(string); contentExist {
		fields["content"] = content
	}
	if excerpt, excerptExist := params.Args["excerpt"].(string); excerptExist {
		fields["excerpt"] = excerpt
	}
	if status, statusExist := params.Args["status"].(string); statusExist {
		fields["status"] = status
	}
	if slug, slugExist := params.Args["slug"].(string); slugExist {
		fields["slug"] = slug
		if ! utils.DB.Where(&Post{Type: postConfig.Slug, Slug: slug}).Not(&Post{ID: id}).First(&post).RecordNotFound() {
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

func DeletePost(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	id, idExist := params.Args["id"].(int)
	if ! idExist {
		return nil, nil
	}

	var post = Post{}
	post.ID = id

	if err := utils.DB.Delete(&post).Error; err != nil {
		return nil, err
	}
	if err := utils.DB.Where(&PostMeta{PostID: id}).Delete(&PostMeta{}).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func GetMeta(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	keys, keysExist := params.Args["keys"].([]interface{})
	post, postExist := params.Source.(Post)

	if ! postExist {
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
