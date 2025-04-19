package data

import (
	"net/http"
	"slices"
	"strconv"
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

func (p *Pagination) GetOrderBy(allowedOrderBy []string, defaultOrderBy string) string {
	if !slices.Contains(allowedOrderBy, p.OrderBy) {
		return defaultOrderBy
	}

	return p.OrderBy
}
