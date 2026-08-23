package contract

import (
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"strconv"
)

type PageRequest struct {
	Limit, Offset int
	Cursor        string
}

func ParsePage(limitRaw, offsetRaw string) (PageRequest, error) {
	limit, offset := 50, 0
	var err error
	if limitRaw != "" {
		limit, err = strconv.Atoi(limitRaw)
		if err != nil {
			return PageRequest{}, fmt.Errorf("%w: limit", domain.ErrValidation)
		}
	}
	if offsetRaw != "" {
		offset, err = strconv.Atoi(offsetRaw)
		if err != nil {
			return PageRequest{}, fmt.Errorf("%w: offset", domain.ErrValidation)
		}
	}
	if limit < 1 || limit > 200 || offset < 0 {
		return PageRequest{}, domain.ErrValidation
	}
	return PageRequest{Limit: limit, Offset: offset}, nil
}
func TotalPages(total, limit int) int {
	if total <= 0 || limit <= 0 {
		return 0
	}
	pages := total / limit
	if total%limit != 0 {
		pages++
	}
	return pages
}
