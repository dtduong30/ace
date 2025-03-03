package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// Pagination struct to hold pagination info
type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

// GetPagination extracts pagination params from query and returns Pagination
func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

// CalculateTotalPage calculates total pages based on total items and limit
func (p *Pagination) CalculateTotalPage(total int64) {
	p.Total = total
	p.TotalPage = (total + int64(p.Limit) - 1) / int64(p.Limit)
}
