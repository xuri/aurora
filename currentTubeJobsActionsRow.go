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

// currentTubeJobsActionsRow render a section include kick, pause and unpause
// job button by given server and tube.
func currentTubeJobsActionsRow(server string, tube string) string {
	var err error
	var bstkConn *beanstalk.Conn
	var buf, pauseTimeLeft strings.Builder
	var pause = strconv.Itoa(selfConf.TubePauseSeconds)
	if pause == "-1" {
		pause = "3600"
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return ``
	}
	tubeStats := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	statsMap, _ := tubeStats.Stats()
	if statsMap["pause-time-left"] == "0" {
		pauseTimeLeft.WriteString(`<a class="btn btn-default btn-sm" href="?server=`)
		pauseTimeLeft.WriteString(server)
		pauseTimeLeft.WriteString(`&tube=`)
		pauseTimeLeft.WriteString(url.QueryEscape(tube))
		pauseTimeLeft.WriteString(`&action=pause&count=-1" title="Temporarily prevent jobs being reserved from the given tube. Pause for: `)
		pauseTimeLeft.WriteString(pause)
		pauseTimeLeft.WriteString(` seconds"><i class="glyphicon glyphicon-pause"></i> Pause tube</a>`)
	} else {
		pauseTimeLeft.WriteString(`<a class="btn btn-default btn-sm" href="?server=`)
		pauseTimeLeft.WriteString(server)
		pauseTimeLeft.WriteString(`&tube=`)
		pauseTimeLeft.WriteString(url.QueryEscape(tube))
		pauseTimeLeft.WriteString(`&action=pause&count=0" title="Pause seconds left: `)
		pauseTimeLeft.WriteString(statsMap["pause-time-left"])
		pauseTimeLeft.WriteString(`"><i class="glyphicon glyphicon-play"></i> Unpause tube</a>`)
	}
	bstkConn.Close()
	buf.WriteString(`<section id="actionsRow"><b>Actions:</b> &nbsp;<a class="btn btn-default btn-sm" href="?server=`)
	buf.WriteString(server)
	buf.WriteString(`&tube=`)
	buf.WriteString(url.QueryEscape(tube))
	buf.WriteString(`&action=kick&count=1"><i class="glyphicon glyphicon-forward"></i> Kick 1 job</a> <form method="GET"><div class="btn-group" role="group"><button type="submit" class="btn btn-default btn-sm" style="margin-right: -2px;"><i class="glyphicon glyphicon-fast-forward"></i> Kick more </button><input type="hidden" name="server" value="`)
	buf.WriteString(server)
	buf.WriteString(`"><input type="hidden" name="tube" value="`)
	buf.WriteString(url.QueryEscape(tube))
	buf.WriteString(`"><input type="hidden" name="action" value="kick"><input type="number" value="10" name="count" min="0" step="1" size="4" class="form-control input-sm" style="padding: 5px 2px 5px 12px; text-align: center;"></div></form> `)
	buf.WriteString(pauseTimeLeft.String())
	buf.WriteString(` &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <div class="btn-group"><a data-toggle="modal" class="btn btn-success btn-sm" href="#" id="addJob"><i class="glyphicon glyphicon-plus-sign glyphicon-white"></i> Add job</a><button class="btn btn-success btn-sm dropdown-toggle" data-toggle="dropdown"><span class="caret"></span></button><ul class="dropdown-menu">`)
	buf.WriteString(currentTubeJobsActionsRowSample(server, tube))
	buf.WriteString(`</ul></div></section>`)
	return buf.String()
}

// currentTubeJobsActionsRowSample render a dropdown sample list by given server
// and tube.
func currentTubeJobsActionsRowSample(server string, tube string) string {
	sample := strings.Builder{}
	for _, v := range sampleJobs.Tubes {
		if v.Name != tube {
			continue
		}
		if len(v.Keys) == 0 {
			continue
		}
		for _, k := range v.Keys {
			for _, j := range sampleJobs.Jobs {
				if j.Key != k {
					continue
				}
				sample.WriteString(`<li><a href="?server=`)
				sample.WriteString(server)
				sample.WriteString(`&tube=`)
				sample.WriteString(url.QueryEscape(tube))
				sample.WriteString(`&action=loadSample&key=`)
				sample.WriteString(j.Key)
				sample.WriteString(`">`)
				sample.WriteString(html.EscapeString(j.Name))
				sample.WriteString(`</a></li>`)
			}
		}
	}
	if sample.String() == "" {
		return `<li><a href="javascript:void(0);">There are no sample jobs</a></li>`
	}
	sample.WriteString(`<li class="divider"></li><li><a href="./sample?action=manageSamples">Manage samples</a></li>`)
	return sample.String()
}
