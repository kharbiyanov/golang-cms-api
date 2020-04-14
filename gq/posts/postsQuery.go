package main

import (
	"cms-api/config"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"log"
	"strings"
)

type PostQuery struct {
	tx         *gorm.DB
	params     graphql.ResolveParams
	postConfig models.PostConfig
	posts      []models.Post
}

func (query *PostQuery) init() error {

	query.setTables()
	query.setSearch()
	query.setStatus()
	query.setOrder()
	query.setLimit()
	query.setGroup()

	if err := query.setTaxQuery(); err != nil {
		return err
	}

	if err := query.setMetaQuery(); err != nil {
		return err
	}

	if err := query.setDateQuery(); err != nil {
		return err
	}

	if err := query.tx.Find(&query.posts).Error; err != nil {
		return err
	}

	return nil
}

func (query *PostQuery) setSearch() {
	if pSearch, ok := query.params.Args["search"].(string); ok {
		exact := false
		sentence := false
		if pExact, ok := query.params.Args["exact"].(bool); ok {
			exact = pExact
		}
		if pSentence, ok := query.params.Args["sentence"].(bool); ok {
			sentence = pSentence
		}

		sql := "posts.title LIKE ? OR posts.content LIKE ? OR posts.excerpt LIKE ?"

		if !sentence {
			searchArray := strings.Split(pSearch, " ")
			for _, search := range searchArray {
				if exact {
					search = fmt.Sprintf("%%%s%%", search)
				}
				query.tx = query.tx.Where(sql, search, search, search)
			}
		} else {
			if exact {
				pSearch = fmt.Sprintf("%%%s%%", pSearch)
			}
			query.tx = query.tx.Where(sql, pSearch, pSearch, pSearch)
		}
	}
}

func (query *PostQuery) setTables() {
	lang, _ := query.params.Args["lang"].(string)
	query.tx = query.tx.Table("posts").
		Select("posts.*").
		Joins("LEFT JOIN translations t ON t.element_id = posts.id").
		Where("posts.type = ? AND t.lang = ? AND t.element_type = ?", query.postConfig.Slug, lang, fmt.Sprintf("post_%s", query.postConfig.Slug))
}

func (query *PostQuery) setOrder() {
	pOrderBy, pOrderByExist := query.params.Args["order_by"].(string)
	pOrder, _ := query.params.Args["order"].(string)

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

	query.tx = query.tx.Order(fmt.Sprintf("%s %s", orderBy, order))
}

func (query *PostQuery) setStatus() {
	status := "publish"
	if pStatus, ok := query.params.Args["status"].(string); ok {
		switch pStatus {
		case "publish", "draft", "pending", "trash", "any":
			status = pStatus
		}
	}
	if status == "any" {
		query.tx = query.tx.Where("posts.status <> ?", "trash")
	} else {
		query.tx = query.tx.Where("posts.status = ?", status)
	}
}

func (query *PostQuery) setTaxQuery() error {
	if taxQueryParam, ok := query.params.Args["tax_query"].([]interface{}); ok {
		var taxQuery []models.TaxQuery
		if err := mapstructure.Decode(taxQueryParam, &taxQuery); err != nil {
			return err
		}

		for i, tQuery := range taxQuery {
			if len(tQuery.Terms) == 0 {
				continue
			}

			termIDs, termsErr := getTermIDs(tQuery)

			if termsErr != nil {
				return termsErr
			}

			if len(termIDs) == 0 {
				continue
			}

			switch tQuery.Operator {
			case "NOT IN":
				innerSql := utils.DB.Table("term_relationships tr").Select("tr.post_id").Where("tr.term_id IN (?)", termIDs).QueryExpr()
				query.tx = query.tx.Where("posts.id NOT IN (?)", innerSql)
			case "AND":
				innerSql := utils.DB.Table("term_relationships tr").Select("COUNT(1)").Where("tr.term_id IN (?) AND tr.post_id = posts.id", termIDs).QueryExpr()
				query.tx = query.tx.Where("(?) = ?", innerSql, len(termIDs))
			default:
				tName := fmt.Sprintf("tr_%d", i)
				query.tx = query.tx.Joins(fmt.Sprintf("LEFT JOIN term_relationships %[1]s ON %[1]s.post_id = posts.id", tName))
				query.tx = query.tx.Where(fmt.Sprintf("%s.term_id IN (?)", tName), termIDs)
			}
		}
	}

	return nil
}

func (query *PostQuery) setMetaQuery() error {
	if metaQueryParam, ok := query.params.Args["meta_query"].([]interface{}); ok {
		var metaQuery []models.MetaQuery
		if err := mapstructure.Decode(metaQueryParam, &metaQuery); err != nil {
			return err
		}

		for i, tQuery := range metaQuery {
			if len(tQuery.Value) == 0 {
				continue
			}

			tName := fmt.Sprintf("pm_%d", i)
			query.tx = query.tx.Joins(fmt.Sprintf("LEFT JOIN post_meta %[1]s ON %[1]s.post_id = posts.id", tName))

			switch tQuery.Compare {
			case "=", "!=", ">", ">=", "<", "<=", "LIKE", "NOT LIKE":
				value := tQuery.Value[0]

				if tQuery.Compare == "LIKE" || tQuery.Compare == "NOT LIKE" {
					value = fmt.Sprintf("%%%v%%", value)
				}

				switch tQuery.Compare {
				case "!=", "NOT LIKE":
					query.tx = query.tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value %[2]s ? OR %[1]s.value IS NULL", tName, tQuery.Compare), tQuery.Key, value)
				default:
					query.tx = query.tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value %[2]s ?", tName, tQuery.Compare), tQuery.Key, value)
				}
			case "NOT IN":
				var value interface{}
				if len(tQuery.Value) == 1 {
					value = tQuery.Value[0]
				} else {
					value = tQuery.Value
				}
				query.tx = query.tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value NOT IN (?) OR %[1]s.value IS NULL", tName), tQuery.Key, value)
			case "BETWEEN", "NOT BETWEEN":
				if len(tQuery.Value) < 2 {
					continue
				}
				query.tx = query.tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value %[2]s ? AND ?", tName, tQuery.Compare), tQuery.Key, tQuery.Value[0], tQuery.Value[1])
			default:
				if len(tQuery.Value) == 1 {
					query.tx = query.tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value = ?", tName), tQuery.Key, tQuery.Value[0])
				} else {
					query.tx = query.tx.Where(fmt.Sprintf("%[1]s.key = ? AND %[1]s.value IN (?)", tName), tQuery.Key, tQuery.Value)
				}
			}
		}
	}

	return nil
}

func (query *PostQuery) setDateQuery() error {
	if dateQueryParam, ok := query.params.Args["date_query"].([]interface{}); ok {
		var dateQuery []models.DateQuery
		if err := mapstructure.Decode(dateQueryParam, &dateQuery); err != nil {
			return err
		}

		for _, dQuery := range dateQuery {
			column := "posts.created_at"
			compare := "="
			beforeCompare := "<"
			afterCompare := ">"

			if dQuery.Column == "updated_at" {
				column = "posts.updated_at"
			}

			if dQuery.Compare != "" {
				compare = dQuery.Compare
			}

			if dQuery.Inclusive {
				beforeCompare = "<="
				afterCompare = ">="
			}

			if dQuery.Before != "" {
				query.tx = query.tx.Where(fmt.Sprintf("%s %s ?", column, beforeCompare), dQuery.Before)
			}

			if dQuery.After != "" {
				query.tx = query.tx.Where(fmt.Sprintf("%s %s ?", column, afterCompare), dQuery.After)
			}

			var dateQueries []*models.DateQueries

			if len(dQuery.Year) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "year",
					Values:   dQuery.Year,
				})
			}

			if len(dQuery.DayOfYear) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "doy",
					Values:   dQuery.DayOfYear,
				})
			}

			if len(dQuery.Month) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "month",
					Values:   dQuery.Month,
				})
			}

			if len(dQuery.Week) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "week",
					Values:   dQuery.Week,
				})
			}

			if len(dQuery.Day) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "day",
					Values:   dQuery.Day,
				})
			}

			if len(dQuery.DayOfWeek) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "dow",
					Values:   dQuery.DayOfWeek,
				})
			}

			if len(dQuery.Hour) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "hour",
					Values:   dQuery.Hour,
				})
			}

			if len(dQuery.Minute) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "minute",
					Values:   dQuery.Minute,
				})
			}

			if len(dQuery.Second) > 0 {
				dateQueries = append(dateQueries, &models.DateQueries{
					DataPart: "second",
					Values:   dQuery.Second,
				})
			}

			if len(dateQueries) > 0 {
				for _, dateQuery := range dateQueries {
					if compare == "IN" || compare == "NOT IN" {
						query.tx = query.tx.Where(fmt.Sprintf("date_part('%s',%s) %s (?)", dateQuery.DataPart, column, compare), dateQuery.Values)
					} else if compare == "BETWEEN" || compare == "NOT BETWEEN" {
						query.tx = query.tx.Where(fmt.Sprintf("date_part('%s',%s) %s ? AND ?", dateQuery.DataPart, column, compare), dateQuery.Values[0], dateQuery.Values[1])
					} else {
						for _, val := range dateQuery.Values {
							query.tx = query.tx.Where(fmt.Sprintf("date_part('%s',%s) %s ?", dateQuery.DataPart, column, compare), val)
						}
					}
				}
			}
		}
	}

	return nil
}

func (query *PostQuery) setLimit() {
	first := config.Get().DefaultPostsLimit

	if pFirst, ok := query.params.Args["first"].(int); ok {
		first = pFirst
	}
	if offset, ok := query.params.Args["offset"].(int); ok {
		query.tx = query.tx.Offset(offset)
	}

	query.tx = query.tx.Limit(first)
}

func (query *PostQuery) setGroup() {
	query.tx = query.tx.Group("posts.id")
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
