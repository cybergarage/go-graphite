// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// QueryTargetRegexp is 'query' parameter identifier for Metrics API.
	QueryTargetRegexp string = "query"
	// QueryTarget is 'target' parameter identifier for Render API.
	QueryTarget string = "target"
	// QueryFrom is 'from' parameter identifier for Render API.
	QueryFrom string = "from"
	// QueryUntil is 'until' parameter identifier for Render API.
	QueryUntil string = "until"
	// QueryFormat is 'format' parameter identifier for Render API.
	QueryFormat string = "format"
	// QueryFormatTypeCompleter is a format type for Metrics API.
	QueryFormatTypeCompleter string = "completer"
	// QueryFormatTypeTreeJSON is a format type for Metrics API.
	QueryFormatTypeTreeJSON string = "treejson"
	// QueryFormatTypeRaw is a format type for Render API.
	QueryFormatTypeRaw string = "raw"
	// QueryFormatTypeCSV is a format type for Render API.
	QueryFormatTypeCSV string = "csv"
	// QueryFormatTypeJSON is a format type for Render API.
	QueryFormatTypeJSON string = "json"
	// QueryContentTypeRaw is a content type for the CSV format.
	QueryContentTypeRaw string = "text/plain"
	// QueryContentTypeCSV is a content type for the CSV format.
	QueryContentTypeCSV string = "text/csv"
	// QueryContentTypeJSON is a content type for the JSON format.
	QueryContentTypeJSON string = "application/json"
)

// Query is an instance for Render query protocol.
type Query struct {
	Target string
	From   *time.Time
	Until  *time.Time
	Format string
}

// NewQuery returns a new query.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func NewQuery() *Query {
	now := time.Now()
	from := now.Add(-(time.Duration(24) * time.Hour))
	q := &Query{
		Target: "",
		From:   &from, // it defaults to 24 hours ago.
		Until:  &now,  // it defaults to the current time (now).
		Format: QueryContentTypeCSV,
	}
	return q
}

// NewQueryWithQuery copies a query.
func NewQueryWithQuery(oq *Query) *Query {
	q := &Query{
		Target: oq.Target,
		From:   oq.From,  // FIXME : Shallow copy
		Until:  oq.Until, // FIXME : Shallow copy
		Format: oq.Format,
	}
	return q
}

// ParseHTTPRequest parses the specified Render request.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (q *Query) ParseHTTPRequest(httpReq *http.Request) error {
	err := httpReq.ParseForm()
	if err != nil {
		return err
	}

	return q.ParseURLValues(httpReq.Form)
}

// ParseURLValues parses the specified parameters in a Render request.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (q *Query) ParseURLValues(urlValues url.Values) error {
	var err error

	for key, values := range urlValues {
		switch key {
		// For Metrics API
		case QueryTargetRegexp:
			if 0 < len(values) {
				q.Target = values[0]
			}
		// For Render API
		case QueryTarget:
			if 0 < len(values) {
				q.Target = values[0]
			}
		case QueryFrom:
			if 0 < len(values) {
				q.From, err = q.parseTimeString(values[0])
				if err != nil {
					return err
				}
			}
		case QueryUntil:
			if 0 < len(values) {
				q.Until, err = q.parseTimeString(values[0])
				if err != nil {
					return err
				}
			}
		case QueryFormat:
			if 0 < len(values) {
				q.Format = values[0]
			}
		}
	}

	return nil
}

func (q *Query) parseTimeString(timeStr string) (*time.Time, error) {
	if IsRelativeTimeString(timeStr) {
		return RelativeTimeStringToTime(timeStr)
	}

	if IsAbsoluteTimeString(timeStr) {
		return AbsoluteTimeStringToTime(timeStr)
	}

	return nil, fmt.Errorf(errorQueryInvalidTimeFormat, timeStr)
}

// FindMetricsURL returns a path for Metrics API
// The Metrics API
// https://graphite-api.readthedocs.io/en/latest/api.html
func (q *Query) FindMetricsURL(host string, port int) (string, error) {
	if len(q.Target) == 0 {
		return "", fmt.Errorf("%s is not specified", QueryTarget)
	}

	url := fmt.Sprintf("http://%s:%d%s?%s=%s",
		host,
		port,
		renderDefaultFindRequestPath,
		QueryTargetRegexp,
		q.Target)

	return url, nil
}

// RenderURLString returns a path for Render API
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (q *Query) RenderURLString(host string, port int) (string, error) {
	if len(q.Target) == 0 {
		return "", fmt.Errorf("%s is not specified", QueryTarget)
	}

	params := make(map[string]string)

	params[QueryTarget] = q.Target

	if q.From != nil {
		params[QueryFrom] = q.From.Format(queryAbsoluteTimeFormat)
	}

	if q.Until != nil {
		params[QueryUntil] = q.Until.Format(queryAbsoluteTimeFormat)
	}

	if 0 < len(q.Format) {
		params[QueryFormat] = q.Format
	} else {
		params[QueryFormat] = QueryFormatTypeCSV
	}

	query := ""
	for key, value := range params {
		if 0 < len(query) {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", key, value)
	}

	url := fmt.Sprintf("http://%s:%d%s?%s", host, port, renderDefaultQueryRequestPath, query)

	return url, nil
}
