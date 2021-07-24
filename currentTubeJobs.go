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

import "strings"

// currentTubeJobs call currentTubeJobsSummaryTable, currentTubeJobsActionsRow
// and currentTubeJobsShowcase functions by given server and tube config, and
// merge these functions return value.
func currentTubeJobs(server string, tube string) string {
	var table = currentTubeJobsSummaryTable(server, tube)
	buf := strings.Builder{}
	if table == `` {
		buf.WriteString(`Tube "`)
		buf.WriteString(tube)
		buf.WriteString(`" not found or it is empty <br><br><a href="./server?server=`)
		buf.WriteString(server)
		buf.WriteString(`"> &lt;&lt; back </a>`)
		return buf.String()
	}
	buf.WriteString(table)
	buf.WriteString(currentTubeJobsActionsRow(server, tube))
	buf.WriteString(currentTubeJobsShowcase(server, tube))
	return buf.String()
}
