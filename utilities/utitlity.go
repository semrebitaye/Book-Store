package utilities

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PaginationParam struct {
	Page         string `form:"page"`
	PerPage      string `form:"per_page"`
	Sort         string `form:"sort"`
	Search       string `form:"search"`
	Filter       string `form:"filter"`
	LinkOperator string `form:"link_operator"`
}
type Filter struct {
	ColumnName string `json:"column_name"`
	Operator   string `json:"operator"`
	Value      any    `json:"value"`
}

type Sort struct {
	ColumnName string `json:"column_name"`
	Value      string `json:"value"`
}

type FilterParam struct {
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Sort    Sort     `json:"sort"`
	Search  string   `json:"search"`
	Filters []Filter `json:"filters"`
}

func ExtractPagination(param PaginationParam) (FilterParam, error) {

	page, err := strconv.Atoi(param.Page)
	if err != nil || page <= 0 {
		page = 1
	}
	fmt.Println("page int:", page)
	per_page, err := strconv.Atoi(param.PerPage)
	if err != nil || per_page <= 0 {
		per_page = 10 // Default limit
	}
	fmt.Println("perpage int:", per_page)
	var sort Sort
	if param.Sort == "" {
		sort.ColumnName = "created_at"
		sort.Value = "asc"

	} else {
		err = json.Unmarshal([]byte(param.Sort), &sort)
		if err != nil {
			return FilterParam{}, err
		}
	}

	var filter []Filter
	if param.Filter != "" {
		err = json.Unmarshal([]byte(param.Filter), &filter)
		if err != nil {
			return FilterParam{}, err
		}
	}

	return FilterParam{
		Page:    page,
		PerPage: per_page,
		Search:  param.Search,
		Sort:    sort,
		Filters: filter,
	}, nil
}
