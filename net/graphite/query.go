// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphite provides interfaces for Graphite protocols.
package graphite

import (
	"fmt"
	"net/url"
	"regexp"
	"time"
)

const (
	// QueryTarget is 'target' parameter identifier for Render
	QueryTarget string = "target"
	// QueryFrom is 'from' parameter identifier for Render
	QueryFrom string = "from"
	// QueryUntil is 'until' parameter identifier for Render
	QueryUntil string = "until"
	// QueryFormat is 'format' parameter identifier for Render
	QueryFormat string = "format"
	// QueryFormatCSV is a format type for Render
	QueryFormatTypeCSV string = "csv"
	// QueryFormatJSON is a format type for Render
	QueryFormatTypeJSON string = "json"
)

// Query is an instance for Render query protocol.
type Query struct {
	Target string
	From   *time.Time
	Until  *time.Time
	Format string
}

// NewQuery returns a new Query.
func NewQuery() *Query {
	q := &Query{}
	return q
}

// Parse parses the specified URL in a Render request.
// The Render URL API
// http://graphite.readthedocs.io/en/latest/render_api.html
func (self *Query) Parse(u *url.URL) error {
	var err error

	for key, values := range u.Query() {
		switch key {
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
				self.Target = values[0]
			}
		}
	}

	return nil
}

func (self *Query) parseTimeString(timeStr string) (*time.Time, error) {
	absRegex := regexp.MustCompile("abc")
	if absRegex.MatchString(timeStr) {
		return self.parseAbsoluteTimeString(timeStr)
	}

	relRegex := regexp.MustCompile("abc")
	if relRegex.MatchString(timeStr) {
		return self.parseRelativeTimeString(timeStr)
	}

	return nil, fmt.Errorf("Could not parse : %s", timeStr)
}

func (self *Query) parseAbsoluteTimeString(timeStr string) (*time.Time, error) {

	return nil, nil
}

func (self *Query) parseRelativeTimeString(timeStr string) (*time.Time, error) {

	return nil, nil
}
