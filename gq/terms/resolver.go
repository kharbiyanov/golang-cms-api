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

func GetTerms(params graphql.ResolveParams) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)
	tax, _ := params.Args["taxonomy"].(string)

	var terms []models.Term

	tx := utils.DB.Table("terms").
		Select("terms.*").
		Joins("left join translations on translations.element_id = terms.id").
		Where("terms.taxonomy = ? and translations.lang = ? and translations.element_type = ?", tax, lang, fmt.Sprintf("tax_%s", tax))

	if parent, ok := params.Args["parent"].(int); ok {
		tx = tx.Where("terms.parent = ?", parent)
	}
	if first, ok := params.Args["first"].(int); ok {
		tx = tx.Limit(first)
	}
	if offset, ok := params.Args["offset"].(int); ok {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&terms).Error; err != nil {
		return nil, err
	}
	return terms, nil
}

func UpdateTerm(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	tax, _ := params.Args["taxonomy"].(string)

	fields := make(map[string]interface{})

	var term models.Term

	if name, ok := params.Args["name"].(string); ok {
		fields["name"] = name
	}
	if description, ok := params.Args["description"].(string); ok {
		fields["description"] = description
	}
	if parent, ok := params.Args["parent"].(int); ok {
		fields["parent"] = parent
		if utils.DB.Where(&models.Term{ID: parent, Taxonomy: tax}).First(&models.Term{}).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.TermParentIDNotFoundMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
	}
	if slug, ok := params.Args["slug"].(string); ok {
		fields["slug"] = slug
		if !utils.DB.Where(&models.Term{Taxonomy: tax, Slug: slug}).Not(&models.Term{ID: id}).First(&term).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.TermSlugExistMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
	}

	term.ID = id

	if err := utils.DB.Model(&term).Updates(fields).Scan(&term).Error; err != nil {
		return nil, err
	}

	return term, nil
}

func DeleteTerm(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)

	term := models.Term{}

	if utils.DB.First(&term, id).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.TermNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	return nil, utils.DB.Delete(term).Error
}
