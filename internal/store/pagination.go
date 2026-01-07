package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationParams struct {
	Limit  int      `json:"limit" validate:"min=1,max=100"`
	Offset int      `json:"offset" validate:"min=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (params *PaginationParams) Parse(r *http.Request) (*PaginationParams, error) {
	query := r.URL.Query()

	limitStr := query.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}
		params.Limit = limit
	}
	offsetStr := query.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, err
		}
		params.Offset = offset
	}
	sort := query.Get("sort")
	if sort != "" {
		params.Sort = sort
	}

	tags := query.Get("tags")
	if tags != "" {
		params.Tags = strings.Split(tags, ",")
	}

	search := query.Get("search")
	if search != "" {
		params.Search = search
	}

	since := query.Get("since")
	if since != "" {
		params.Since = parseTime(since)
	}

	until := query.Get("until")
	if until != "" {
		params.Until = parseTime(until)
	}
	return params, nil
}

func parseTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Format(time.DateTime)
	}
	if t, err := time.Parse(time.DateTime, s); err == nil {
		return t.Format(time.DateTime)
	}
	return ""
}
