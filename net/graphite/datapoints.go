// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

// DataPoints is a slice of DataPoint.
type DataPoints []*DataPoint

// NewDataPoint returns a new datapoint slice.
func NewDataPoints(size int) []*DataPoint {
	return make([]*DataPoint, size)
}

func (dps DataPoints) Len() int {
	return len(dps)
}

func (dps DataPoints) Swap(i, j int) {
	dps[i], dps[j] = dps[j], dps[i]
}

func (dps DataPoints) Less(i, j int) bool {
	return dps[i].Timestamp.Unix() < dps[j].Timestamp.Unix()
}
