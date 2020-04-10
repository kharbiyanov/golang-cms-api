package main

import (
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"log"
	"strings"
)

func SetSearch(tx *gorm.DB, params graphql.ResolveParams) *gorm.DB {
	if pSearch, ok := params.Args["search"].(string); ok {
		exact := false
		sentence := false
		if pExact, ok := params.Args["exact"].(bool); ok {
			exact = pExact
		}
		if pSentence, ok := params.Args["sentence"].(bool); ok {
			sentence = pSentence
		}

		query := "posts.title LIKE ? OR posts.content LIKE ? OR posts.excerpt LIKE ?"

		if !sentence {
			searchArray := strings.Split(pSearch, " ")
			for _, search := range searchArray {
				if exact {
					search = fmt.Sprintf("%%%s%%", search)
				}
				tx = tx.Where(query, search, search, search)
			}
		} else {
			if exact {
				pSearch = fmt.Sprintf("%%%s%%", pSearch)
			}
			tx = tx.Where(query, pSearch, pSearch, pSearch)
		}
	}

	return tx
}

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

func SetMetaQuery(tx *gorm.DB, params graphql.ResolveParams) (*gorm.DB, error) {
	if metaQueryParam, ok := params.Args["meta_query"].([]interface{}); ok {
		var metaQuery []models.MetaQuery
		if err := mapstructure.Decode(metaQueryParam, &metaQuery); err != nil {
			return nil, err
		}

		for i, query := range metaQuery {
			if len(query.Value) == 0 {
				continue
			}

			tName := fmt.Sprintf("pm_%d", i)
			tx = tx.Joins(fmt.Sprintf("LEFT JOIN post_meta %[1]s ON %[1]s.post_id = posts.id", tName))

			switch query.Compare {
			case "=", "!=", ">", ">=", "<", "<=", "LIKE", "NOT LIKE":
				value := query.Value[0]

				if query.Compare == "LIKE" || query.Compare == "NOT LIKE" {
					value = fmt.Sprintf("%%%v%%", value)
				}

				switch query.Compare {
				case "!=", "NOT LIKE":
					tx = tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value %[2]s ? OR %[1]s.value IS NULL", tName, query.Compare), query.Key, value)
				default:
					tx = tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value %[2]s ?", tName, query.Compare), query.Key, value)
				}
			case "NOT IN":
				var value interface{}
				if len(query.Value) == 1 {
					value = query.Value[0]
				} else {
					value = query.Value
				}
				tx = tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value NOT IN (?) OR %[1]s.value IS NULL", tName), query.Key, value)
			case "BETWEEN", "NOT BETWEEN":
				if len(query.Value) < 2 {
					continue
				}
				tx = tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value %[2]s ? AND ?", tName, query.Compare), query.Key, query.Value[0], query.Value[1])
			default:
				if len(query.Value) == 1 {
					tx = tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value = ?", tName), query.Key, query.Value[0])
				} else {
					tx = tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value IN (?)", tName), query.Key, query.Value)
				}
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
