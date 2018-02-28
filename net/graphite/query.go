// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"net/url"
	"time"
)

const (
	// QueryTargetRegexp is 'query' parameter identifier for Metrics API
	QueryTargetRegexp string = "query"
	// QueryTarget is 'target' parameter identifier for Render API
	QueryTarget string = "target"
	// QueryFrom is 'from' parameter identifier for Render API
	QueryFrom string = "from"
	// QueryUntil is 'until' parameter identifier for Render API
	QueryUntil string = "until"
	// QueryFormat is 'format' parameter identifier for Render API
	QueryFormat string = "format"
	// QueryFormatTypeCompleter is a format type for Metrics API
	QueryFormatTypeCompleter string = "completer"
	// QueryFormatTypeTreeJSON is a format type for Metrics API
	QueryFormatTypeTreeJSON string = "treejson"
	// QueryFormatTypeRaw is a format type for Render API
	QueryFormatTypeRaw string = "raw"
	// QueryFormatTypeCSV is a format type for Render API
	QueryFormatTypeCSV string = "csv"
	// QueryFormatTypeJSON is a format type for Render API
	QueryFormatTypeJSON string = "json"
	// QueryContentTypeRaw is a content type for the CSV format
	QueryContentTypeRaw string = "text/plain"
	// QueryContentTypeCSV is a content type for the CSV format
	QueryContentTypeCSV string = "text/csv"
	// QueryContentTypeJSON is a content type for the JSON format
	QueryContentTypeJSON string = "application/json"
)

// Query is an instance for Render query protocol.
type Query struct {
	Target string
	From   *time.Time
	Until  *time.Time
	Format string
}

// NewQuery returns a new Query.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func NewQuery() *Query {
	now := time.Now()
	from := now.Add(-(time.Duration(24) * time.Hour))
	q := &Query{
		Target: "",
		From:   &from, // it defaults to 24 hours ago.
		Until:  &now,  // it defaults to the current time (now).
	}
	return q
}

// Parse parses the specified URL in a Render request.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Query) Parse(u *url.URL) error {
	var err error

	for key, values := range u.Query() {
		switch key {
		// For Metrics API
		case QueryTargetRegexp:
			if 0 < len(values) {
				self.Target = values[0]
			}
		// For Render API
		case QueryTarget:
			if 0 < len(values) {
				self.Target = values[0]
			}
		case QueryFrom:
			if 0 < len(values) {
				self.From, err = self.parseTimeString(values[0])
				if err != nil {
					return err
				}
			}
		case QueryUntil:
			if 0 < len(values) {
				self.Until, err = self.parseTimeString(values[0])
				if err != nil {
					return err
				}
			}
		case QueryFormat:
			if 0 < len(values) {
				self.Format = values[0]
			}
		}
	}

	return nil
}

func (self *Query) parseTimeString(timeStr string) (*time.Time, error) {
	if IsRelativeTimeString(timeStr) {
		return RelativeTimeStringToTime(timeStr)
	}

	if IsAbsoluteTimeString(timeStr) {
		return AbsouleteTimeStringToTime(timeStr)
	}

	return nil, fmt.Errorf(errorQueryInvalidTimeFormat, timeStr)
}

// RenderURLString returns a path for Render API
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Query) RenderURLString(host string, port int) (string, error) {
	if len(self.Target) <= 0 {
		return "", fmt.Errorf("%s is not specified", QueryTarget)
	}

	params := make(map[string]string)

	params[QueryTarget] = self.Target

	if self.From != nil {
		params[QueryFrom] = self.From.Format(queryAbsoluteTimeFormat)
	}

	if self.Until != nil {
		params[QueryUntil] = self.Until.Format(queryAbsoluteTimeFormat)
	}

	if 0 < len(self.Format) {
		params[QueryFormat] = self.Format
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
