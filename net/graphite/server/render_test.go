// Copyright 2017 The go-graphite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"testing"

	"github.com/cybergarage/go-graphite/net/graphite"
	"github.com/cybergarage/go-graphite/net/graphite/client"
)

type TestRender struct {
	*Render
	QueryCount int
}

func NewTestRender() *TestRender {
	render := &TestRender{NewRender(), 0}
	render.RenderListener = render
	return render
}

func (self *TestRender) QueryRequestReceived(query *graphite.Query, err error) ([]*graphite.Metric, error) {
	if err != nil {
		return nil, nil
	}
	self.QueryCount++
	return nil, nil
}

func TestNewRender(t *testing.T) {
	NewRender()
}

func TestRenderQuery(t *testing.T) {
	render := NewTestRender()
	render.RenderListener = render
	err := render.Start()
	if err != nil {
		t.Error(err)
	}

	cli := client.NewClient()

	loopCount := 0
	for n := 0; n < 10; n++ {
		q := graphite.NewQuery()
		q.Target = fmt.Sprintf("path%d", n)
		_, err := cli.PostQuery(q)
		if err != nil {
			t.Error(err)
		}
		loopCount++
	}

	if render.QueryCount != loopCount {
		t.Error(fmt.Errorf("%d != %d", render.QueryCount, loopCount))
	}

	err = render.Stop()
	if err != nil {
		t.Error(err)
	}
}