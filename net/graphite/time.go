// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphite provides interfaces for Graphite protocols.
package graphite

import (
	"fmt"
	"regexp"
	"time"
)

const (
	queryRelativeTimeSecondsRegex = "-[0-9]+s"
	queryRelativeTimeMinutesRegex = "-[0-9]+min"
	queryRelativeTimeHoursRegex   = "-[0-9]+h"
	queryRelativeTimeDaysRegex    = "-[0-9]+d"
	queryRelativeTimeWeeksRegex   = "-[0-9]+w"
	queryRelativeTimeMonthsRegex  = "-[0-9]+mon"
	queryRelativeTimeYearsRegex   = "-[0-9]+y"

	queryRelativeTimeFormat = "-%d%s"

	queryAbsoluteTimeNow    = "now"
	queryAbsoluteTimeRegex  = "[0-9]{2}:0-9]{2}_0-9]{8}"
	queryAbsoluteTimeFormat = "15:04_20060102"
)

// IsRelativeTimeString returns the specified string based on the releative time format.
func IsRelativeTimeString(timeStr string) bool {
	queryRelativeTimeRegexs := []string{
		queryRelativeTimeSecondsRegex,
		queryRelativeTimeMinutesRegex,
		queryRelativeTimeHoursRegex,
		queryRelativeTimeDaysRegex,
		queryRelativeTimeWeeksRegex,
		queryRelativeTimeMonthsRegex,
		queryRelativeTimeYearsRegex,
	}

	for _, regex := range queryRelativeTimeRegexs {
		matched, _ := regexp.MatchString(regex, timeStr)
		if matched {
			return true
		}
	}

	return false
}

// IsRelativeTimeString returns the specified string based on the releative time format.
func RelativeTimeStringToTime(timeStr string) (time.Time, error) {
	queryRelativeTimeRegexs := []string{
		queryRelativeTimeSecondsRegex,
		queryRelativeTimeMinutesRegex,
		queryRelativeTimeHoursRegex,
		queryRelativeTimeDaysRegex,
		queryRelativeTimeWeeksRegex,
		queryRelativeTimeMonthsRegex,
		queryRelativeTimeYearsRegex,
	}

	now := time.Now()

	for n, regex := range queryRelativeTimeRegexs {
		matched, _ := regexp.MatchString(regex, timeStr)
		if !matched {
			continue
		}

		switch n {
		case 0:
			return now.Add(-(1 * time.Second)), nil
		}

		return now, fmt.Errorf(errorQueryInvalidRelativeTimeFormat, timeStr)
	}

	return now, fmt.Errorf(errorQueryInvalidRelativeTimeFormat, timeStr)
}
