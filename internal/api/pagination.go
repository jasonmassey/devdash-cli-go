package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// PaginatedResponse represents a response with cursor-based pagination.
type PaginatedResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"nextCursor"`
}

// FetchAll retrieves all pages for a paginated endpoint.
// The path should include query params; pagination cursor is appended.
func FetchAll[T any](c *Client, path string) ([]T, error) {
	var all []T
	cursor := ""

	for {
		p := path
		if cursor != "" {
			sep := "&"
			if !containsQuery(p) {
				sep = "?"
			}
			p += sep + "cursor=" + url.QueryEscape(cursor)
		}

		data, err := c.Get(p)
		if err != nil {
			return nil, err
		}

		var page PaginatedResponse[T]
		if err := json.Unmarshal(data, &page); err != nil {
			// Try unmarshalling as plain array (some endpoints return arrays directly)
			var items []T
			if err2 := json.Unmarshal(data, &items); err2 == nil {
				return append(all, items...), nil
			}
			return nil, fmt.Errorf("failed to parse paginated response: %w", err)
		}

		all = append(all, page.Data...)

		if page.NextCursor == "" {
			break
		}
		cursor = page.NextCursor
	}

	return all, nil
}

func containsQuery(path string) bool {
	for _, c := range path {
		if c == '?' {
			return true
		}
	}
	return false
}
