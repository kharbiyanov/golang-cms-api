package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

func CreateTerm(params graphql.ResolveParams) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)

	term := &models.Term{
		Name:     params.Args["name"].(string),
		Taxonomy: params.Args["taxonomy"].(string),
		Slug:     params.Args["slug"].(string),
	}

	if description, ok := params.Args["description"].(string); ok {
		term.Description = description
	}

	if parent, ok := params.Args["parent"].(int); ok {
		term.Parent = parent
		if utils.DB.Where(&models.Term{ID: parent, Taxonomy: term.Taxonomy}).First(&models.Term{}).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.TermParentIDNotFoundMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
	}

	if !utils.DB.Where(&models.Term{Slug: term.Slug, Taxonomy: term.Taxonomy}).First(&models.Term{}).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.TermSlugExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Create(&term).Scan(&term).Error; err != nil {
		return nil, err
	}

	translation := models.Translation{
		ElementType: fmt.Sprintf("tax_%s", term.Taxonomy),
		ElementID:   term.ID,
		Lang:        lang,
	}

	if err := utils.DB.Create(&translation).Scan(&translation).Error; err != nil {
		return nil, err
	}

	return term, nil
}
