package main

import (
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"log"
)

func SetOrder(tx *gorm.DB, params graphql.ResolveParams) *gorm.DB {
	pOrderBy, pOrderByExist := params.Args["order_by"].(string)
	pOrder, _ := params.Args["order"].(string)

	orderBy := "posts.created_at"
	order := "desc"

	if pOrderByExist {
		switch pOrderBy {
		case "date":
			orderBy = "posts.created_at"
		case "updated":
			orderBy = "posts.updated_at"
		case "author":
			orderBy = "posts.author_id"
		case "title":
			orderBy = "posts.title"
		case "content":
			orderBy = "posts.content"
		case "status":
			orderBy = "posts.status"
		case "slug":
			orderBy = "posts.slug"
		}
	}

	if pOrder == "asc" {
		order = "asc"
	}

	return tx.Order(fmt.Sprintf("%s %s", orderBy, order))
}

func SetTaxQuery(tx *gorm.DB, params graphql.ResolveParams) (*gorm.DB, error) {
	if taxQueryParam, ok := params.Args["tax_query"].([]interface{}); ok {
		var taxQuery []models.TaxQuery
		if err := mapstructure.Decode(taxQueryParam, &taxQuery); err != nil {
			return nil, err
		}

		for i, query := range taxQuery {
			if len(query.Terms) == 0 {
				continue
			}

			termIDs, termsErr := getTermIDs(query)

			if termsErr != nil {
				return nil, termsErr
			}

			if len(termIDs) == 0 {
				continue
			}

			switch query.Operator {
			case "NOT IN":
				innerSql := utils.DB.Table("term_relationships tr").Select("tr.post_id").Where("tr.term_id IN (?)", termIDs).QueryExpr()
				tx = tx.Where("posts.id NOT IN (?)", innerSql)
			case "AND":
				innerSql := utils.DB.Table("term_relationships tr").Select("COUNT(1)").Where("tr.term_id IN (?) AND tr.post_id = posts.id", termIDs).QueryExpr()
				tx = tx.Where("(?) = ?", innerSql, len(termIDs))
			default:
				tName := fmt.Sprintf("tr_%d", i)
				tx = tx.Joins(fmt.Sprintf("LEFT JOIN term_relationships %[1]s ON %[1]s.post_id = posts.id", tName))
				tx = tx.Where(fmt.Sprintf("%s.term_id IN (?)", tName), termIDs)
			}
		}
	}

	return tx, nil
}

func getTermIDs(query models.TaxQuery) ([]int, error) {
	taxonomy := "category"
	field := "id"

	if query.Taxonomy != "" {
		taxonomy = query.Taxonomy
	}

	if query.Field == "slug" {
		field = "slug"
	}

	termRows, termErr := utils.DB.Model(&models.Term{}).Where(fmt.Sprintf("taxonomy = ? AND %s IN (?)", field), taxonomy, query.Terms).Select("id").Rows()
	defer termRows.Close()

	if termErr != nil {
		return nil, termErr
	}

	var termIDs []int

	for termRows.Next() {
		var termID int
		if err := termRows.Scan(&termID); err != nil {
			log.Println("err:", err)
			return nil, err
		}
		termIDs = append(termIDs, termID)
	}

	return termIDs, nil
}
