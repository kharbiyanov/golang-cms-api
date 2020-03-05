package posts

import (
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

func GetPost(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	id, idOK := params.Args["id"].(int)
	if ! idOK {
		return nil, nil
	}
	var post Post
	if err := utils.DB.Where(&Post{ID: id, Type: fmt.Sprintf("%s", postConfig.Slug)}).First(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetPosts(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	var posts []Post
	if err := utils.DB.Where(&Post{Type: fmt.Sprintf("%s", postConfig.Slug)}).Find(&posts).Error; err != nil {
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
	if err := utils.DB.Create(post).Scan(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetMeta(params graphql.ResolveParams, postConfig PostConfig) (interface{}, error) {
	keys, keysOK := params.Args["keys"].([]interface{})
	post, postOK := params.Source.(Post)
	metaQuery := utils.DB
	if ! postOK {
		return nil, nil
	}
	if keysOK {
		metaQuery = metaQuery.Where("key in(?)", keys)
	}
	if err := metaQuery.Model(&post).Association("Meta").Find(&post.Meta).Error; err != nil {
		return nil, err
	}
	return post.Meta, nil
}
