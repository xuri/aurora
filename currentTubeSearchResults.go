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
)

// currentTubeSearchResults constructs a search result table by given server,
// tube, search result limit and search content.
func currentTubeSearchResults(server string, tube string, limit string, searchStr string, result []SearchResult) string {
	var buf, tr strings.Builder
	if len(result) == 0 {
		buf.WriteString(`<br/>No results found for <b>`)
		buf.WriteString(html.EscapeString(searchStr))
		buf.WriteString(`</b> in tube: <b>`)
		buf.WriteString(tube)
		buf.WriteString(`</b>`)
		return buf.String()
	}
	for k, job := range result {
		tr.WriteString(`<tr><td>`)
		tr.WriteString(strconv.Itoa(int(job.ID)))
		tr.WriteString(`</td><td>`)
		tr.WriteString(job.State)
		tr.WriteString(`</td><td class="ellipsize">`)
		tr.WriteString(html.EscapeString(job.Data))
		tr.WriteString(`</td><td><div class="dropdown btn-group-xs"><button class="btn btn-default dropdown-toggle" type="button" id="dropdownMenu`)
		tr.WriteString(strconv.Itoa(k))
		tr.WriteString(`" data-toggle="dropdown" aria-expanded="true"> Actions <span class="caret"></span></button><ul class="dropdown-menu" role="menu" aria-labelledby="dropdownMenu`)
		tr.WriteString(strconv.Itoa(k))
		tr.WriteString(`"><li role="presentation"><a role="menuitem" class="addSample" data-jobid="`)
		tr.WriteString(strconv.Itoa(int(job.ID)))
		tr.WriteString(`" href="?server=`)
		tr.WriteString(server)
		tr.WriteString(`&tube=`)
		tr.WriteString(url.QueryEscape(tube))
		tr.WriteString(`&action=addSample"><i class="glyphicon glyphicon-plus glyphicon-white"></i> Add to samples </a></li><li role="presentation"><a role="menuitem" href="?server=`)
		tr.WriteString(server)
		tr.WriteString(`&tube=`)
		tr.WriteString(url.QueryEscape(tube))
		tr.WriteString(`&state=`)
		tr.WriteString(job.State)
		tr.WriteString(`&action=deleteJob&jobid=`)
		tr.WriteString(strconv.Itoa(int(job.ID)))
		tr.WriteString(`"><i class="glyphicon glyphicon-remove glyphicon-white"></i> Delete</a> </li><li role="presentation"><a role="menuitem" href="?server=`)
		tr.WriteString(server)
		tr.WriteString(`&tube=`)
		tr.WriteString(url.QueryEscape(tube))
		tr.WriteString(`&state=`)
		tr.WriteString(job.State)
		tr.WriteString(`&action=kickJob&jobid=`)
		tr.WriteString(strconv.Itoa(int(job.ID)))
		tr.WriteString(`"><i class="glyphicon glyphicon-forward glyphicon-white"></i> Kick </a></li></ul></div></td></tr>`)
	}
	buf.WriteString(`<section id="actionsRow"><a class="btn btn-default btn-sm" href="?server=`)
	buf.WriteString(server)
	buf.WriteString(`&tube=`)
	buf.WriteString(url.QueryEscape(tube))
	buf.WriteString(`"><i class="glyphicon glyphicon-backward"></i>  &nbsp;Back to tube</a></section><section id="searchResult"><div class="row"><div class="col-sm-12"><table class="table table-striped table-hover" style="table-layout:fixed;"><thead><tr><th class="col-md-1">id</th><th class="col-md-1">state</th><th>data</th><th class="col-md-1">action</th></tr></thead><tbody>`)
	buf.WriteString(tr.String())
	buf.WriteString(`</tbody></table></div></div>First `)
	buf.WriteString(limit)
	buf.WriteString(` rows are displayed for each state.<br/><br/></section>`)
	return buf.String()
}
