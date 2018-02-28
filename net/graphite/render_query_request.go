// Copyright (C) 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphite

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// handleRenderRequest handles requests for Render API.
// The Render URL API
// http://readthedocs.io/en/latest/render_api.html
func (self *Render) handleRenderRequest(httpWriter http.ResponseWriter, httpReq *http.Request) {
	query := NewQuery()
	err := query.ParseHTTPRequest(httpReq)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	if self.RenderListener == nil {
		self.responseInternalServerError(httpWriter, httpReq)
		return
	}

	metrics, err := self.RenderListener.QueryMetricsRequestReceived(query, nil)
	if err != nil {
		self.responseBadRequest(httpWriter, httpReq)
		return
	}

	self.responseQueryMetrics(httpWriter, httpReq, query, metrics)
}

func (self *Render) responseQueryMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	switch query.Format {
	case QueryFormatTypeRaw:
		self.responseQueryRawMetrics(httpWriter, httpReq, query, metrics)
		return
	case QueryFormatTypeCSV:
		self.responseQueryCSVMetrics(httpWriter, httpReq, query, metrics)
		return
	case QueryFormatTypeJSON:
		self.responseQueryJSONMetrics(httpWriter, httpReq, query, metrics)
		return
	}

	self.responseBadRequest(httpWriter, httpReq)
}

func (self *Render) responseQueryRawMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeRaw)
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	httpWriter.WriteHeader(http.StatusOK)

	for _, m := range metrics {
		err := m.SortDataPoints()
		if err != nil {
			continue
		}
		dpCount := m.GetDataPointCount()
		if dpCount <= 0 {
			continue
		}

		var from, until time.Time
		var step int64

		switch dpCount {
		case 0:
			continue
		case 1:
			firstDp, err := m.GetDataPoint(0)
			if err != nil {
				continue
			}
			from = firstDp.Timestamp
			until = from
			step = 0
		default:
			firstDp, err := m.GetDataPoint(0)
			if err != nil {
				continue
			}
			from = firstDp.Timestamp

			lastDp, err := m.GetDataPoint((dpCount - 1))
			if err != nil {
				continue
			}
			until = lastDp.Timestamp

			secondDp, err := m.GetDataPoint(1)
			if err != nil {
				continue
			}
			// FIXME : Step is calculated only using the first few data points.
			step = secondDp.Timestamp.Unix() - firstDp.Timestamp.Unix()
		}

		msg := fmt.Sprintf("%s,%d,%d,%d|", m.Name, from.Unix(), until.Unix(), step)
		httpWriter.Write([]byte(msg))

		for n, dp := range m.DataPoints {
			var value string
			switch n {
			case (dpCount - 1):
				value = fmt.Sprintf("%f", dp.Value)
			default:
				value = fmt.Sprintf("%f,", dp.Value)
			}
			httpWriter.Write([]byte(value))
		}

		httpWriter.Write([]byte("\n"))
	}
}

func (self *Render) responseQueryCSVMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeCSV)
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	httpWriter.WriteHeader(http.StatusOK)

	for _, m := range metrics {
		for _, dp := range m.DataPoints {
			mRow := fmt.Sprintf("%s,%s\n", m.Name, dp.RenderCSVString())
			httpWriter.Write([]byte(mRow))
		}
	}
}

func (self *Render) responseQueryJSONMetrics(httpWriter http.ResponseWriter, httpReq *http.Request, query *Query, metrics []*Metrics) {
	httpWriter.Header().Set(httpHeaderContentType, QueryContentTypeJSON)
	httpWriter.Header().Set(httpHeaderAccessControlAllowOrigin, httpHeaderAccessControlAllowOriginAll)
	httpWriter.WriteHeader(http.StatusOK)

	httpWriter.Write([]byte("[\n"))

	mCount := len(metrics)
	for i, m := range metrics {
		httpWriter.Write([]byte("{\n"))

		// Output the target name
		httpWriter.Write([]byte(fmt.Sprintf("\"target\": \"%s\",\n", m.Name)))

		// Output the datapoint array
		dpCount := m.GetDataPointCount()
		httpWriter.Write([]byte("\"datapoints\": [\n"))
		for j, dp := range m.DataPoints {
			if !math.IsNaN(dp.Value) {
				httpWriter.Write([]byte(fmt.Sprintf("[%f,", dp.Value)))
			} else {
				httpWriter.Write([]byte("[null,"))
			}
			httpWriter.Write([]byte(fmt.Sprintf("%d", dp.UnixTimestamp())))
			if j < (dpCount - 1) {
				httpWriter.Write([]byte("],\n"))
			} else {
				httpWriter.Write([]byte("]\n"))
			}
		}
		httpWriter.Write([]byte("]\n"))

		if i < (mCount - 1) {
			httpWriter.Write([]byte("},\n"))
		} else {
			httpWriter.Write([]byte("}\n"))
		}
	}

	httpWriter.Write([]byte("]\n"))
}
