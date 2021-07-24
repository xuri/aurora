// Copyright 2016 - 2021 The aurora Authors. All rights reserved. Use of this
// source code is governed by a MIT license that can be found in the LICENSE
// file.
//
// The aurora is a web-based beanstalkd queue server console written in Go
// and works on macOS, Linux and Windows machines. Main idea behind using Go
// for backend development is to utilize ability of the compiler to produce
// zero-dependency binaries for multiple platforms. aurora was created as an
// attempt to build very simple and portable application to work with local or
// remote beanstalkd server.

package main

import (
	"html"
	"net/url"
	"strconv"
	"strings"

	"github.com/xuri/aurora/beanstalk"
)

// currentTubeJobsShowcase return a section include three stats of job, call
// currentTubeJobsShowcaseSections function and get that return value based on
// the given server and tube config.
func currentTubeJobsShowcase(server string, tube string) string {
	var buf strings.Builder
	buf.WriteString(`<section class="jobsShowcase">`)
	buf.WriteString(currentTubeJobsShowcaseSections(server, tube))
	buf.WriteString(`</section>`)
	return buf.String()
}

// currentTubeJobsShowcaseSections constructs a tube job in ready, delayed and
// buried stats table based on the given server and tube config.
func currentTubeJobsShowcaseSections(server string, tube string) string {
	stats := []string{"ready", "delayed", "buried"}
	var err error
	var buf, s, j, b, m, r strings.Builder
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return `<hr><div class="pull-left"><h3>Next job in "ready" state</h3></div><div class="clearfix"></div><i>empty</i><hr><div class="pull-left"><h3>Next job in "delayed" state</h3></div><div class="clearfix"></div><i>empty</i><hr><div class="pull-left"><h3>Next job in "buried" state</h3></div><div class="clearfix"></div><i>empty</i>`
	}
	tubeStats := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	peek := []func() (id uint64, body []byte, err error){tubeStats.PeekReady, tubeStats.PeekDelayed, tubeStats.PeekBuried}
	for k, stat := range stats {
		s.Reset()
		j.Reset()
		b.Reset()
		m.Reset()
		r.Reset()
		tubes, _ := bstkConn.ListTubes()
		jobID, jobBody, err := peek[k]()
		if err != nil {
			buf.WriteString(`<hr><div class="pull-left"><h3>Next job in "`)
			buf.WriteString(stat)
			buf.WriteString(`" state</h3></div><div class="clearfix"></div><i>empty</i>`)
			continue
		}
		statsJob, err := bstkConn.StatsJob(jobID)
		if err != nil {
			continue
		}
		for _, v := range jobStatsOrder {
			s.WriteString(`<tr><td>`)
			s.WriteString(v)
			s.WriteString(`</td><td>`)
			s.WriteString(statsJob[v])
			s.WriteString(`</td></tr>`)
		}
		for _, v := range tubes {
			m.WriteString(`<li><a href="?server=`)
			m.WriteString(server)
			m.WriteString(`&tube=`)
			m.WriteString(url.QueryEscape(tube))
			m.WriteString(`&action=moveJobsTo&destTube=`)
			m.WriteString(url.QueryEscape(v))
			m.WriteString(`&state=`)
			m.WriteString(stat)
			m.WriteString(`">`)
			m.WriteString(html.EscapeString(v))
			m.WriteString(`</a></li>`)
		}
		if jobBody != nil {
			b.WriteString(`<div class="pull-right"><div style="margin-bottom: 3px;"><a class="btn btn-sm btn-info addSample" data-jobid="`)
			b.WriteString(strconv.Itoa(int(jobID)))
			b.WriteString(`" href="?server=`)
			b.WriteString(server)
			b.WriteString(`&tube=`)
			b.WriteString(url.QueryEscape(tube))
			b.WriteString(`&action=addSample"><i class="glyphicon glyphicon-plus glyphicon-white"></i> Add to samples</a> <div class="btn-group"> <button class="btn btn-info btn-sm dropdown-toggle" data-toggle="dropdown"> <i class="glyphicon glyphicon-arrow-right glyphicon-white"></i> Move all `)
			b.WriteString(stat)
			b.WriteString(` to </button><ul class="dropdown-menu"><li><input class="moveJobsNewTubeName input-medium" type="text" data-href="?server=`)
			b.WriteString(server)
			b.WriteString(`&tube=`)
			b.WriteString(url.QueryEscape(tube))
			b.WriteString(`&action=moveJobsTo&state=`)
			b.WriteString(stat)
			b.WriteString(`&destTube=" placeholder="New tube name"/></li>`)
			b.WriteString(m.String())
			b.WriteString(`<li class="divider"></li><li><a href="?server=`)
			b.WriteString(server)
			b.WriteString(`&tube=`)
			b.WriteString(url.QueryEscape(tube))
			b.WriteString(`&action=moveJobsTo&destState=buried&state=`)
			b.WriteString(stat)
			b.WriteString(`">Buried</a></li></ul></div> <a class="btn btn-sm btn-danger" href="?server=`)
			b.WriteString(server)
			b.WriteString(`&tube=`)
			b.WriteString(url.QueryEscape(tube))
			b.WriteString(`&state=`)
			b.WriteString(stat)
			b.WriteString(`&action=deleteAll&count=1" onclick="return confirm('This process might hang a while on tubes with lots of jobs. Are you sure you want to continue?');"><i class="glyphicon glyphicon-trash glyphicon-white"></i> Delete all `)
			b.WriteString(stat)
			b.WriteString(` jobs</a> <a class="btn btn-sm btn-danger" href="?server=`)
			b.WriteString(server)
			b.WriteString(`&tube=`)
			b.WriteString(url.QueryEscape(tube))
			b.WriteString(`&state=`)
			b.WriteString(stat)
			b.WriteString(`&action=deleteJob&jobid=`)
			b.WriteString(strconv.Itoa(int(jobID)))
			b.WriteString(`"><i class="glyphicon glyphicon-remove glyphicon-white"></i> Delete</a></div></div>`)
		}
		if jobBody != nil {
			j.WriteString(preformat(jobBody))
		}
		if jobBody != nil {
			r.WriteString(`<hr><div class="pull-left"><h3>Next job in "`)
			r.WriteString(stat)
			r.WriteString(`" state</h3></div><div class="clearfix"></div><div class="row show-grid"><div class="col-sm-3"><table class="table"><thead><tr><th>Stats:</th><th>&nbsp;</th></tr></thead><tbody>`)
			r.WriteString(s.String())
			r.WriteString(`</tbody></table></div><div class="col-sm-9"><div class="clearfix"><div class="pull-left"><b>Job data:</b></div>`)
			r.WriteString(b.String())
			r.WriteString(`</div><pre><code>`)
			r.WriteString(j.String())
			r.WriteString(`</code></pre></div></div>`)
		} else {
			r.WriteString(`<hr><div class="pull-left"><h3>Next job in "`)
			r.WriteString(stat)
			r.WriteString(`" state</h3></div><div class="clearfix"></div><i>empty</i>`)
		}
		buf.WriteString(r.String())
	}
	bstkConn.Close()
	return buf.String()
}
