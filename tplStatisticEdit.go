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
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/aurora/beanstalk"
)

// tplStatisticEdit provide method to render the statistics preference page.
func tplStatisticEdit(alert string) string {
	var err error
	var buf, savedTo, ST, tubeList strings.Builder
	frequency := selfConf.StatisticsFrequency
	if frequency < 1 {
		frequency = 300
	}
	for _, server := range selfConf.Servers {
		var bstkConn *beanstalk.Conn
		tubeList.Reset()
		if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
			continue
		}
		tubes, _ := bstkConn.ListTubes()
		sort.Strings(tubes)
		bstkConn.Close()
		for _, v := range tubes {
			var checked string
			s, ok := statisticsData.Server[server]
			if ok {
				_, ok := s[v]
				if ok {
					checked = `checked="checked"`
				}
			}
			tubeList.WriteString(`<div class="control-group"><div class="controls"><label class="checkbox-inline"><input type="checkbox" name="tubes[`)
			tubeList.WriteString(server)
			tubeList.WriteString(`:`)
			tubeList.WriteString(v)
			tubeList.WriteString(`]" value="1" `)
			tubeList.WriteString(checked)
			tubeList.WriteString(`>`)
			tubeList.WriteString(v)
			tubeList.WriteString(`</label></div></div>`)
		}
		ST.WriteString(`<div class="pull-left" style="padding-right: 35px;">`)
		ST.WriteString(server)
		ST.WriteString(`<blockquote>`)
		ST.WriteString(tubeList.String())
		ST.WriteString(`</blockquote></div>`)
	}

	buf.WriteString(`<form name="statisticsPreference" action="./statistics?action=save" method="POST"><div class="clearfix form-group"><div class="pull-left"><h4 class="text-info">Statistics preference</h4></div></div><div class="form-group"><fieldset>`)
	buf.WriteString(alert)
	buf.WriteString(`<div class="control-group"><label class="control-label"><b>Collection record number of each server or tube (Default: <i>0</i>, reserved for not statistics, recommended value: <i>300</i>) *</b></label><div class="controls form-group"><input class="form-control input-sm focused" name="collection" type="number" min="0" style="width: 15em;" required="" value="`)
	buf.WriteString(strconv.Itoa(selfConf.StatisticsCollection))
	buf.WriteString(`" autocomplete="off"></div></div><div class="control-group"><label class="control-label"><b>Acquisition frequency seconds of each server or tube (Default: <i>300</i>, minimum: <i>1</i>) *</b></label><div class="controls form-group"><input class="form-control input-sm" name="frequency" type="number" min="1" style="width: 15em;" required="" value="`)
	buf.WriteString(strconv.Itoa(frequency))
	buf.WriteString(`" autocomplete="off"></div></div></fieldset><div class="clearfix"><label class="control-label"><b>Available on tubes *</b></label><br/>`)
	buf.WriteString(savedTo.String())
	buf.WriteString(ST.String())
	buf.WriteString(`</div></div><div><input type="submit" class="btn btn-success" value="Save"/></div></form>`)
	return buf.String()
}
