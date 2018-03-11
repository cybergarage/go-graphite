// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"testing"
	"time"
)

func TestReletiveTimes(t *testing.T) {
	relTimeStrings := []string{
		"-30s",
		"-30min",
		"-30h",
		"-30d",
		"-30w",
		"-30mon",
		"-30y",
	}

	relTimeSeconds := []int64{
		30,
		30 * 60,
		30 * 60 * 60,
		30 * 60 * 60 * 24,
		30 * 60 * 60 * 24 * 7,
		30 * 60 * 60 * 24 * 30,
		30 * 60 * 60 * 24 * 365,
	}

	for n, timeStr := range relTimeStrings {
		ok := IsRelativeTimeString(timeStr)
		if !ok {
			t.Error(fmt.Errorf("Not relative time string : %s", timeStr))
		}

		ok = IsAbsoluteTimeString(timeStr)
		if !ok {
			t.Error(fmt.Errorf("Not relative time string : %s", timeStr))
		}

		now := time.Now()
		relTime, err := RelativeTimeStringToTime(timeStr)
		if err != nil {
			t.Error(err)
			continue
		}

		desiredTime := now.Unix() - relTimeSeconds[n]
		if relTime.Unix() != desiredTime {
			t.Error(fmt.Errorf("%s (%d) != %d", timeStr, relTime.Unix(), desiredTime))
		}
	}
}

func TestAbsoluteTimes(t *testing.T) {
	absTimeStrings := []string{
		"now",
		"15:04_20060102",
		"70815600",
	}

	absTimeSeconds := []int64{
		0,
		1136181840,
		70815600,
	}

	for n, timeStr := range absTimeStrings {
		ok := IsRelativeTimeString(timeStr)
		if ok {
			t.Error(fmt.Errorf("Not absolute time string : %s", timeStr))
		}

		ok = IsAbsoluteTimeString(timeStr)
		if !ok {
			t.Error(fmt.Errorf("Not absolute time string : %s", timeStr))
		}

		absTimeSeconds[0] = time.Now().Unix()
		absTime, err := AbsouleteTimeStringToTime(timeStr)
		if err != nil {
			t.Error(err)
			continue
		}

		if absTime.UTC().Unix() != absTimeSeconds[n] {
			t.Error(fmt.Errorf("%s (%d) != %d", timeStr, absTime.UTC().Unix(), absTimeSeconds[n]))
		}
	}
}
