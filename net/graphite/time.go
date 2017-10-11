// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package graphite provides interfaces for Graphite protocols.
package graphite

import (
	"fmt"
	"regexp"
	"strconv"
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

	queryAbsoluteTimeNowRegex  = "now"
	queryAbsoluteTimeRegex     = "[0-9]{2}:[0-9]{2}_[0-9]{8}"
	queryAbsoluteTimeUnixRegex = "[0-9]*"
	queryAbsoluteTimeFormat    = "15:04_20060102"
)

// IsAbsoluteTimeString returns the specified string whether it is based on the absolute time formats.
func IsAbsoluteTimeString(timeStr string) bool {
	queryAbsTimeRegexs := []string{
		queryAbsoluteTimeNowRegex,
		queryAbsoluteTimeRegex,
		queryAbsoluteTimeUnixRegex,
	}

	for _, regex := range queryAbsTimeRegexs {
		matched, _ := regexp.MatchString(regex, timeStr)
		if matched {
			return true
		}
	}

	return false
}

// AbsouleteTimeStringToTime returns a time based on the specified relative time string.
func AbsouleteTimeStringToTime(timeStr string) (*time.Time, error) {
	queryRelativeTimeRegexs := []string{
		queryAbsoluteTimeNowRegex,
		queryAbsoluteTimeRegex,
		queryAbsoluteTimeUnixRegex,
	}

	for n, regex := range queryRelativeTimeRegexs {
		matched, _ := regexp.MatchString(regex, timeStr)
		if !matched {
			continue
		}

		switch n {
		case 0: // queryAbsoluteTimeNowRegex
			now := time.Now()
			return &now, nil
		case 1: // queryAbsoluteTimeRegex
			time, err := time.Parse(queryAbsoluteTimeFormat, timeStr)
			if err != nil {
				return nil, err
			}
			return &time, nil
		case 2: // queryAbsoluteTimeUnixRegex
			unixTime, err := strconv.ParseInt(timeStr, 10, 64)
			if err != nil {
				break
			}
			time := time.Unix(unixTime, 0)
			return &time, nil
		}

		return nil, fmt.Errorf(errorQueryInvalidTimeFormat, timeStr)
	}

	return nil, fmt.Errorf(errorQueryInvalidTimeFormat, timeStr)
}

// IsRelativeTimeString returns the specified string whether it is based on the releative time formats.
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

// RelativeTimeStringToTime returns a time based on the specified relative time string.
func RelativeTimeStringToTime(timeStr string) (*time.Time, error) {
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

		var timeNum int
		var timeUnit string
		_, err := fmt.Sscanf(timeStr, queryRelativeTimeFormat, &timeNum, &timeUnit)
		if err != nil {
			break
		}

		now = time.Now()
		switch n {
		case 0: // queryRelativeTimeSecondsRegex
			time := now.Add(-(time.Duration(timeNum) * time.Second))
			return &time, nil
		case 1: // queryRelativeTimeMinutesRegex
			time := now.Add(-(time.Duration(timeNum) * time.Minute))
			return &time, nil
		case 2: // queryRelativeTimeHoursRegex
			time := now.Add(-(time.Duration(timeNum) * time.Hour))
			return &time, nil
		case 3: // queryRelativeTimeDaysRegex
			time := now.Add(-(time.Duration(timeNum) * time.Hour * 24))
			return &time, nil
		case 4: // queryRelativeTimeWeeksRegex
			time := now.Add(-(time.Duration(timeNum) * time.Hour * 24 * 7))
			return &time, nil
		case 5: // queryRelativeTimeMonthsRegex
			time := now.Add(-(time.Duration(timeNum) * time.Hour * 24 * 30))
			return &time, nil
		case 6: // queryRelativeTimeYearsRegex
			time := now.Add(-(time.Duration(timeNum) * time.Hour * 24 * 365))
			return &time, nil
		}

		return nil, fmt.Errorf(errorQueryInvalidTimeFormat, timeStr)
	}

	return nil, fmt.Errorf(errorQueryInvalidTimeFormat, timeStr)
}
