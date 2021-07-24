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
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/aurora/beanstalk"
)

// getServerStatus render a server stats table.
func getServerStatus() string {
	var err error
	var buf, td, th strings.Builder
	for _, addr := range selfConf.Servers {
		var bstkConn *beanstalk.Conn
		if bstkConn, err = beanstalk.Dial("tcp", addr); err != nil {
			td.WriteString(`<tr><td>`)
			td.WriteString(addr)
			td.WriteString(`</td><td colspan="`)
			td.WriteString(strconv.Itoa(len(selfConf.Filter)))
			td.WriteString(`" class="row-full">&nbsp;</td><td><a class="btn btn-xs btn-danger" title="Remove from list" href="serversRemove?action=serversRemove&removeServer=`)
			td.WriteString(addr)
			td.WriteString(`"><span class="glyphicon glyphicon-minus"> </span></a></td></tr>`)
			continue
		}
		s, _ := bstkConn.Stats()
		bstkConn.Close()
		td.WriteString(`<tr><td><a href="server?server=`)
		td.WriteString(addr)
		td.WriteString(`">`)
		td.WriteString(addr)
		td.WriteString(`</a></td>`)
		for _, v := range selfConf.Filter {
			td.WriteString(`<td>`)
			td.WriteString(s[v])
			td.WriteString(`</td>`)
		}
		td.WriteString(`<td><a class="btn btn-xs btn-danger" title="Remove from list" href="serversRemove?action=serversRemove&removeServer=`)
		td.WriteString(addr)
		td.WriteString(`"><span class="glyphicon glyphicon-minus"> </span></a></td></tr>`)
	}
	for _, v := range selfConf.Filter {
		th.WriteString(`<th>`)
		th.WriteString(v)
		th.WriteString(`</th>`)
	}
	buf.WriteString(`<div class="row"><div class="col-sm-12"><table class="table table-striped table-hover" id="servers-index"><thead><tr><th>name</th>`)
	buf.WriteString(th.String())
	buf.WriteString(`<th>&nbsp;</th></tr></thead><tbody>`)
	buf.WriteString(td.String())
	buf.WriteString(`</tbody></table><a href="#servers-add" role="button" class="btn btn-info" id="addServer">Add server</a></div></div>`)
	return buf.String()
}

// getServerTubes render a tubes stats table by given server.
func getServerTubes(server string) string {
	var err error
	var buf, th, tr, td strings.Builder
	var bstkConn *beanstalk.Conn
	for _, v := range selfConf.TubeFilters {
		th.WriteString(`<th>`)
		th.WriteString(v)
		th.WriteString(`</th>`)
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		buf.WriteString(`<div id="idAllTubes"><section id="summaryTable"><div class="row"><div class="col-sm-12"><table class="table table-striped table-hover"><thead><tr><th>name</th>`)
		buf.WriteString(th.String())
		buf.WriteString(`</tr></thead><tbody></tbody></table></div></div></section></div>`)
		return buf.String()
	}
	tubes, _ := bstkConn.ListTubes()
	sort.Strings(tubes)
	for _, v := range tubes {
		tubeStats := &beanstalk.Tube{
			Conn: bstkConn,
			Name: v,
		}
		statsMap, err := tubeStats.Stats()
		if err != nil {
			continue
		}
		for _, stats := range selfConf.TubeFilters {
			td.WriteString(`<td>`)
			td.WriteString(statsMap[stats])
			td.WriteString(`</td>`)
		}
		tr.WriteString(`<tr><td><a href="tube?server=`)
		tr.WriteString(server)
		tr.WriteString(`&tube=`)
		tr.WriteString(url.QueryEscape(v))
		tr.WriteString(`">`)
		tr.WriteString(v)
		tr.WriteString(`</a></td>`)
		tr.WriteString(td.String())
		tr.WriteString(`</tr>`)
		td.Reset()
	}
	bstkConn.Close()
	buf.WriteString(`<div id="idAllTubes"><section id="summaryTable"><div class="row"><div class="col-sm-12"><table class="table table-striped table-hover"><thead><tr><th>name</th>`)
	buf.WriteString(th.String())
	buf.WriteString(`</tr></thead><tbody>`)
	buf.WriteString(tr.String())
	buf.WriteString(`</tbody></table></div></div></section></div>`)
	return buf.String()
}

// dropDownServer render a navigation dropdown menu for server list.
func dropDownServer(currentServer string) string {
	var ul strings.Builder
	if currentServer == "" {
		currentServer = `All servers`
	}
	ul.WriteString(`<li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown">`)
	ul.WriteString(currentServer)
	ul.WriteString(` <span class="caret"></span></a><ul class="dropdown-menu">`)
	for _, addr := range selfConf.Servers {
		if addr == currentServer {
			continue
		}
		ul.WriteString(`<li><a href="./server?server=`)
		ul.WriteString(addr)
		ul.WriteString(`">`)
		ul.WriteString(addr)
		ul.WriteString(`</a></li>`)
	}
	if currentServer != "All servers" {
		ul.WriteString(`<li><a href="./public">All servers</a></li>`)
	}
	ul.WriteString(`</ul></li>`)
	return ul.String()
}

// dropDownTube render a navigation dropdown menu for tube list.
func dropDownTube(server string, currentTube string) string {
	var ul strings.Builder
	if currentTube == "" {
		currentTube = `All tubes`
	}
	ul.WriteString(`<li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown">`)
	ul.WriteString(currentTube)
	ul.WriteString(` <span class="caret"></span></a><ul class="dropdown-menu">`)
	var bstkConn *beanstalk.Conn
	var err error
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		if currentTube != "" {
			ul.WriteString(`<li><a href="./public">All tubes</a></li>`)
		}
		ul.WriteString(`</ul></li>`)
		return ul.String()
	}
	tubes, _ := bstkConn.ListTubes()
	sort.Strings(tubes)
	for _, v := range tubes {
		ul.WriteString(`<li><a href="./tube?server=`)
		ul.WriteString(server)
		ul.WriteString(`&tube=`)
		ul.WriteString(url.QueryEscape(v))
		ul.WriteString(`">`)
		ul.WriteString(v)
		ul.WriteString(`</a></li>`)
	}
	bstkConn.Close()
	if currentTube != "All tubes" {
		ul.WriteString(`<li><a href="./server?server=`)
		ul.WriteString(server)
		ul.WriteString(`">All tubes</a></li>`)
	}
	ul.WriteString(`</ul></li>`)
	return ul.String()
}

// dropEditSettings render a navigation dropdown menu for set preference.
func dropEditSettings() string {
	var buf strings.Builder
	var isDisabledJSONDecode, isDisabledJobDataHighlight, isEnabledBase64Decode string
	if selfConf.IsDisabledJSONDecode != 1 {
		isDisabledJSONDecode = `checked="checked"`
	}
	if selfConf.IsDisabledJobDataHighlight != 1 {
		isDisabledJobDataHighlight = `checked="checked"`
	}
	if selfConf.IsEnabledBase64Decode != 0 {
		isEnabledBase64Decode = `checked="checked"`
	}
	buf.WriteString(`<div id="settings" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="settings-label" aria-hidden="true"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button><h4 class="modal-title" id="settings-label">Settings</h4></div><div class="modal-body"><fieldset><div class="form-group"><label for="tubePauseSeconds"><b>Tube pause seconds</b> (<i>-1</i> means the default: <i>3600</i>, <i>0</i> is reserved for un-pause)</label><input class="form-control focused" id="tubePauseSeconds" type="number" value="`)
	buf.WriteString(strconv.Itoa(selfConf.TubePauseSeconds))
	buf.WriteString(`"></div><div class="form-group"><label><b>Auto-refresh interval in milliseconds</b> (Default: <i>500</i>)</label><input class="form-control focused" id="autoRefreshTimeoutMs" type="number" value="`)
	buf.WriteString(strconv.Itoa(selfConf.AutoRefreshTimeoutMs))
	buf.WriteString(`"></div><div class="form-group"><label><b>Search result limits</b> (Default: <i>25</i>)</label><input class="form-control focused" id="searchResultLimit" type="number" value="`)
	buf.WriteString(strconv.Itoa(selfConf.SearchResultLimit))
	buf.WriteString(`"></div><div class="form-group"><label><b>Preferred way to deal with job data</b></label><div class="checkbox"><label><input type="checkbox" id="isDisabledJsonDecode" value="1" `)
	buf.WriteString(isDisabledJSONDecode)
	buf.WriteString(`>before display: JSON decode</label></div><div class="checkbox"><label><input type="checkbox" id="isEnabledBase64Decode" value="1" `)
	buf.WriteString(isEnabledBase64Decode)
	buf.WriteString(`>before display: Base64 decode</label></div><div class="checkbox"><label><input type="checkbox" id="isDisabledJobDataHighlight" value="1" `)
	buf.WriteString(isDisabledJobDataHighlight)
	buf.WriteString(`>after display: enable highlight</label></div></div></fieldset></div><div class="modal-footer"><button class="btn" data-dismiss="modal" aria-hidden="true">Close</button></div></div></div></div>`)
	return buf.String()
}
