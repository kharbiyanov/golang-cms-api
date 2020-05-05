package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

func CreateTerm(params graphql.ResolveParams, taxonomyConfig models.TaxonomyConfig) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)

	term := &models.Term{
		Name:     params.Args["name"].(string),
		Taxonomy: taxonomyConfig.Taxonomy,
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

func GetTerms(params graphql.ResolveParams, taxonomyConfig models.TaxonomyConfig) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)

	var terms []models.Term

	tx := utils.DB.Table("terms").
		Select("terms.*").
		Joins("left join translations on translations.element_id = terms.id").
		Where("terms.taxonomy = ? and translations.lang = ? and translations.element_type = ?", taxonomyConfig.Taxonomy, lang, fmt.Sprintf("tax_%s", taxonomyConfig.Taxonomy))

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

func UpdateTerm(params graphql.ResolveParams, taxonomyConfig models.TaxonomyConfig) (interface{}, error) {
	id, _ := params.Args["id"].(int)

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
		if utils.DB.Where(&models.Term{ID: parent, Taxonomy: taxonomyConfig.Taxonomy}).First(&models.Term{}).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.TermParentIDNotFoundMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
	}
	if slug, ok := params.Args["slug"].(string); ok {
		fields["slug"] = slug
		if !utils.DB.Where(&models.Term{Taxonomy: taxonomyConfig.Taxonomy, Slug: slug}).Not(&models.Term{ID: id}).First(&term).RecordNotFound() {
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

func DeleteTerm(params graphql.ResolveParams, taxonomyConfig models.TaxonomyConfig) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	term := models.Term{}

	if utils.DB.Where(&models.Term{ID: id, Taxonomy: taxonomyConfig.Taxonomy}).Find(&term).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.TermNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	return nil, utils.DB.Delete(term).Error
}

func GetMetaInTerm(params graphql.ResolveParams) (interface{}, error) {
	keys, keysExist := params.Args["keys"].([]interface{})
	term, termExist := params.Source.(models.Term)

	if !termExist {
		return nil, nil
	}

	tx := utils.DB

	if keysExist {
		tx = tx.Where("key in(?)", keys)
	}

	if err := tx.Model(&term).Association("Meta").Find(&term.Meta).Error; err != nil {
		return nil, err
	}

	return term.Meta, nil
}

func UpdateMeta(params graphql.ResolveParams) (interface{}, error) {
	termId, _ := params.Args["term_id"].(int)
	key, _ := params.Args["key"].(string)
	value, _ := params.Args["value"].(string)

	meta := models.TermMeta{
		TermID: termId,
		Key:    key,
		Value:  value,
	}

	if err := utils.DB.Model(&meta).Where(&models.TermMeta{TermID: termId, Key: key}).Update(&meta).Scan(&meta).Error; err != nil {
		if err := utils.DB.Save(&meta).First(&meta).Error; err != nil {
			return nil, err
		}
	}

	return meta, nil
}

func DeleteMeta(params graphql.ResolveParams) (interface{}, error) {
	termId, _ := params.Args["term_id"].(int)
	key, _ := params.Args["key"].(string)

	return nil, utils.DB.Delete(&models.TermMeta{}, &models.TermMeta{TermID: termId, Key: key}).Error
}

func SetTerms(params graphql.ResolveParams) (interface{}, error) {
	// TODO: проверка на существование термов

	postId, _ := params.Args["post_id"].(int)
	terms, _ := params.Args["terms"].([]interface{})

	var insertTerms []interface{}

	for _, termId := range terms {
		insertTerms = append(insertTerms, models.TermRelationship{
			PostID: postId,
			TermID: termId.(int),
		})
	}

	if err := utils.DB.Where("post_id = ?", postId).Delete(&models.TermRelationship{}).Error; err != nil {
		return nil, err
	}

	return nil, gormbulk.BulkInsert(utils.DB, insertTerms, 3000)
}

func GetTranslationsInTerm(params graphql.ResolveParams) (interface{}, error) {
	term, termExist := params.Source.(models.Term)

	if !termExist {
		return nil, nil
	}

	elementType := fmt.Sprintf("tax_%s", term.Taxonomy)

	innerSql := utils.DB.Table("translations t").Select("t.group_id").Where("t.element_id = ? AND t.element_type = ?", term.ID, elementType).QueryExpr()

	tx := utils.DB.Table("translations").
		Select("translations.*").
		Joins("LEFT JOIN lang l ON translations.lang = l.code").
		Where("l.deleted_at IS NULL AND translations.element_type = ? AND translations.group_id = (?)", elementType, innerSql)

	if err := tx.Find(&term.Translations).Error; err != nil {
		return nil, err
	}

	return term.Translations, nil
}
