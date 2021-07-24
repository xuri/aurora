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

// tplServerFilterStatsGroups render server filter stats groups checkbox.
func tplServerFilterStatsGroups() []string {
	stats := []string{"", "", "", ""}
	buf := strings.Builder{}
	statsGroupsFilter := [][]map[string]string{binlogStatsGroups, cmdStatsGroups, currentStatsGroups, otherStatsGroups}
	for k, statsGroups := range statsGroupsFilter {
		for _, statsGroup := range statsGroups {
			for property, description := range statsGroup {
				status := ""
				if checkInSlice(selfConf.Filter, property) {
					status = `checked`
				}
				buf.Reset()
				buf.WriteString(`<div class="control-group"><div class="controls"><div class="checkbox"><label><input type="checkbox" name="`)
				buf.WriteString(property)
				buf.WriteString(`" `)
				buf.WriteString(status)
				buf.WriteString(`><b>`)
				buf.WriteString(property)
				buf.WriteString(`</b><br/>`)
				buf.WriteString(description)
				buf.WriteString(`</label></div></div></div>`)
				stats[k] += buf.String()
			}
		}
	}
	return stats
}

// tplServerFilter render modal popup for select server tube stats column.
func tplServerFilter() string {
	filter := strings.Builder{}
	stats := tplServerFilterStatsGroups()
	filter.WriteString(`<div id="filterServer" data-cookie="filter" class="modal fade" tabindex="-1" role="dialog"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button><h3 id="filter-label" class="text-info">Filter columns</h3></div><div class="modal-body"><form class="form-group"><div class="tabbable"><ul class="nav nav-tabs"><li class="active"><a href="#binlog" data-toggle="tab">binlog</a></li><li><a href="#cmd" data-toggle="tab">cmd</a></li><li><a href="#current" data-toggle="tab">current</a></li><li><a href="#other" data-toggle="tab">other</a></li></ul><div class="tab-content"><div class="tab-pane active" id="binlog">`)
	filter.WriteString(stats[0])
	filter.WriteString(`</div><div class="tab-pane" id="cmd">`)
	filter.WriteString(stats[1])
	filter.WriteString(`</div><div class="tab-pane" id="current">`)
	filter.WriteString(stats[2])
	filter.WriteString(`</div><div class="tab-pane" id="other">`)
	filter.WriteString(stats[3])
	filter.WriteString(`</div></div></div></form></div><div class="modal-footer"><button class="btn" data-dismiss="modal" aria-hidden="true">Close</button></div></div></div></div>`)
	return filter.String()
}

// tplTubeFilter render a modal popup for select job stats of tube.
func tplTubeFilter() string {
	var buf, currents, others strings.Builder
	for k, current := range tubeStatFields {
		if k > 7 {
			continue
		}
		for property, description := range current {
			status := ""
			if checkInSlice(selfConf.TubeFilters, property) {
				status = `checked`
			}
			currents.WriteString(`<div class="form-group"><div class="checkbox"><label class="checkbox"><input type="checkbox" name="`)
			currents.WriteString(property)
			currents.WriteString(`" `)
			currents.WriteString(status)
			currents.WriteString(`><b>`)
			currents.WriteString(property)
			currents.WriteString(`</b><br/>`)
			currents.WriteString(description)
			currents.WriteString(`</label></div></div>`)
		}
	}

	for k, other := range tubeStatFields {
		if k < 8 {
			continue
		}
		for property, description := range other {
			status := ""
			if checkInSlice(selfConf.TubeFilters, property) {
				status = `checked`
			}
			others.WriteString(`<div class="form-group"><div class="checkbox"><label class="checkbox"><input type="checkbox" name="`)
			others.WriteString(property)
			others.WriteString(`" `)
			others.WriteString(status)
			others.WriteString(`><b>`)
			others.WriteString(property)
			others.WriteString(`</b><br/>`)
			others.WriteString(description)
			others.WriteString(`</label></div></div>`)
		}
	}
	buf.WriteString(`<div id="filter" data-cookie="tubefilter" class="modal fade" tabindex="-1" role="dialog"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button><h4 class="modal-title" id="filter-columns-label">Filter columns</h4></div><div class="modal-body"><form><div class="tabbable"><ul class="nav nav-tabs"><li class="active"><a href="#current" data-toggle="tab">current</a></li><li><a href="#other" data-toggle="tab">other</a></li></ul><div class="tab-content"><div class="tab-pane active" id="current">`)
	buf.WriteString(currents.String())
	buf.WriteString(`</div><div class="tab-pane" id="other">`)
	buf.WriteString(others.String())
	buf.WriteString(`</div></div></div></form></div><div class="modal-footer"><button class="btn btn-success" data-dismiss="modal" aria-hidden="true">Close</button></div></div></div></div>`)
	return buf.String()
}
