package pagination

import (
	"fmt"
	"net/url"
)

type Page struct{ Limit, Offset int }

func FromQuery(values url.Values) Page {
	limit, offset := 50, 0
	if raw := values.Get("limit"); raw != "" {
		if _, err := fmt.Sscanf(raw, "%d", &limit); err != nil {
			limit = 50
		}
	}
	if raw := values.Get("offset"); raw != "" {
		if _, err := fmt.Sscanf(raw, "%d", &offset); err != nil {
			offset = 0
		}
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return Page{Limit: limit, Offset: offset}
}
