package response

import (
	"github.com/Gustcat/people-info-service/internal/lib/urlbuilder"
	"strconv"
)

const (
	DefaultOffset uint64 = 0
	DefaultLimit  uint64 = 5
)

type Pagination struct {
	Limit    uint64  `json:"limit"`
	Offset   uint64  `json:"offset"`
	Total    uint64  `json:"total"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
}

func NewPagination(limit uint64, offset uint64, total uint64, url string) (*Pagination, error) {
	var prev, next *string
	var err error
	var prevOffset uint64
	limitStr := strconv.FormatUint(limit, 10)

	if offset > 0 {
		if limit < offset {
			prevOffset = offset - limit
		}
		prevOffsetStr := strconv.FormatUint(prevOffset, 10)
		prev = new(string)
		*prev, err = urlbuilder.BuildWithQueryParams(url, map[string]string{"limit": limitStr, "offset": prevOffsetStr})
		if err != nil {
			return nil, err
		}

	}

	if offset+limit < total {
		nextOffsetStr := strconv.FormatUint(offset+limit, 10)
		next = new(string)
		*next, err = urlbuilder.BuildWithQueryParams(url, map[string]string{"limit": limitStr, "offset": nextOffsetStr})
		if err != nil {
			return nil, err
		}
	}

	return &Pagination{
		Limit:    limit,
		Offset:   offset,
		Total:    total,
		Next:     next,
		Previous: prev,
	}, nil
}
