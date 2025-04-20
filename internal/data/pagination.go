package data

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

type Pagination struct {
	Page    int
	PerPage int
	OrderBy string
}

type AllowedField struct {
	FieldName          string
	FormattedFieldName string
}

func GetPagination(r *http.Request) Pagination {
	param := r.URL.Query()
	paramPage := param.Get("page")
	paramPerPage := param.Get("per_page")
	paramOrderBy := param.Get("order_by")

	if paramPage == "" {
		paramPage = strconv.Itoa(DefaultPage)
	}

	if paramPerPage == "" {
		paramPerPage = strconv.Itoa(DefaultPerPage)
	}

	page, _ := strconv.Atoi(paramPage)
	perPage, _ := strconv.Atoi(paramPerPage)

	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return Pagination{
		Page:    page,
		PerPage: perPage,
		OrderBy: paramOrderBy,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *Pagination) GetLimit() int {
	return p.PerPage
}

func (p *Pagination) GetOrderByField(allowedOrderBy []AllowedField, defaultOrderBy string) string {
	parsedOrderBy := p.OrderBy

	if strings.Contains(parsedOrderBy, ",") {
		parsedOrderBy = strings.Split(parsedOrderBy, ",")[0]
	}

	for _, allowedField := range allowedOrderBy {
		if allowedField.FieldName == strings.ToLower(parsedOrderBy) {
			return allowedField.FormattedFieldName
		}
	}

	return defaultOrderBy
}

func (p *Pagination) GetOrderByDirection() string {
	if strings.Contains(p.OrderBy, ",") {
		orderBy := strings.Split(p.OrderBy, ",")
		if len(orderBy) > 1 {
			if orderBy[1] == "desc" {
				return "desc"
			}
		}
	}

	return "asc"
}
