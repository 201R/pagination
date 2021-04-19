package pagination

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Pagination [pagination struct ]
type Pagination struct {
	Limit        int         `json:"limit"`
	Page         int         `json:"page"`
	TotalRows    int64       `json:"totalRows"`
	TotalPages   int         `json:"totalPages"`
	Sort         string      `json:"sort"`
	FirstPage    string      `json:"firstPage"`
	LastPage     string      `json:"lastPage"`
	PreviousPage string      `json:"previousPage"`
	NextPage     string      `json:"nextPage"`
	FromRow      int         `json:"fromRow"`
	ToRow        int         `json:"toRows"`
	Searchs      []Search    `json:"searchs"`
	Rows         interface{} `json:"rows"`
}

// Search [filed you want to search]
type Search struct {
	Column string `json:"column"`
	Action string `json:"action"`
	Query  string `json:"query"`
}

// Offset return number of record to skip before
// starting to return the record from Database
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

// SetTotalRows set total Rows to pagination :: totalRows is anint from DB
func (p *Pagination) SetTotalRows(totalRows int64) {
	p.TotalRows = totalRows
}

// Paginate will init [FromRows, ToRows, TotalPage]
// on the pagination struct
func (p *Pagination) Paginate() *Pagination {

	var (
		fromRows, toRows int
		totalPage        = int(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))
	)

	if p.Page == 1 {
		fromRows = 1
		toRows = p.Limit
	} else if p.Page <= totalPage {
		fromRows = p.Offset() + 1
		toRows = p.Page * p.Limit
	}

	p.FromRow = fromRows
	p.ToRow = toRows
	p.TotalPages = totalPage

	return p
}

// PageLink return an URL for a given page
func (p *Pagination) PageLink(page int) string {
	link := fmt.Sprintf("?limit=%d&page=%d&sort=%s", p.Limit, page, p.Sort) + p.searchQueryParams()
	return link
}

// PageLinkPrev return URL to the previous page
func (p *Pagination) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.Page - 1)
	}
	return
}

// PageLinkNext return URL to the next page
func (p *Pagination) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.Page + 1)
	}
	return
}

// PageLinkFirst return the first page link
func (p *Pagination) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

// PageLinkLast return the last page link
func (p *Pagination) PageLinkLast() (link string) {
	return p.PageLink(p.TotalPages)
}

// HasPrev return true if the current page has a predecessor
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// HasNext return true if the current page has a successor
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

// Return all search query in a given url
func (p *Pagination) searchQueryParams() string {
	searchQueryParams := ""

	for _, search := range p.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}
	return searchQueryParams
}

// NewPaginator instantiates a paginator struct
func New(ctx *gin.Context) *Pagination {
	// Default values
	var (
		limit   = 10
		page    = 1
		sort    = "id desc"
		searchs []Search
	)

	query := ctx.Request.URL.Query()

	for k, v := range query {
		queryValue := v[len(v)-1]

		switch k {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break
		}

		if strings.Contains(k, ".") {
			searchKeys := strings.Split(k, ".")
			search := Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}

			searchs = append(searchs, search)
		}
	}

	return &Pagination{
		Limit:   limit,
		Page:    page,
		Sort:    sort,
		Searchs: searchs,
	}
}