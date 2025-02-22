package util

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultSize = 10
)

// PaginationQuery holds pagination parameters
type PaginationQuery struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

// SetSize parses and sets the size parameter
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n
	return nil
}

// SetPage parses and sets the page parameter
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n
	return nil
}

// SetOrderBy sets the orderBy parameter
func (q *PaginationQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// GetOffset calculates the offset for database queries
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit returns the size limit
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// GetOrderBy returns the orderBy parameter
func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}

// GetPage returns the current page
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// GetSize returns the page size
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

// GetQueryString constructs a query string from the pagination parameters
func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

// GetPaginationFromCtx extracts pagination parameters from the Fiber context
func GetPaginationFromCtx(c *fiber.Ctx) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.Query("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.Query("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.Query("orderBy"))
	return q, nil
}

// GetTotalPages calculates the total number of pages
func GetTotalPages(totalCount int, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

// GetHasMore determines if there are more pages available
func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	if pageSize == 0 {
		return false
	}
	return currentPage < int(math.Ceil(float64(totalCount)/float64(pageSize)))
}
