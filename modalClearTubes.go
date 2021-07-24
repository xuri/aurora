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
	"strings"

	"github.com/xuri/aurora/beanstalk"
)

// modalClearTubes render modal popup for delete job in tubes.
func modalClearTubes(server string) string {
	var err error
	var buf, tubeList strings.Builder
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return ``
	}
	tubes, _ := bstkConn.ListTubes()
	for _, v := range tubes {
		tubeList.WriteString(`<div class="checkbox"><label class=""><input type="checkbox" name="`)
		tubeList.WriteString(v)
		tubeList.WriteString(`" value="1"><b>`)
		tubeList.WriteString(html.EscapeString(v))
		tubeList.WriteString(`</b></label></div>`)
	}
	buf.WriteString(`<div class="modal fade" id="clear-tubes" data-cookie="tubefilter" tabindex="-1" role="dialog" aria-labelledby="clear-tubes-label" aria-hidden="true"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button><h4 class="modal-title" id="clear-tubes-label">Clear multiple tubes</h4></div><div class="modal-body"><form><fieldset><div class="form-group"><label>Tube name <small class="text-muted">(supports <a href="http://james.padolsey.com/javascript/regex-selector-for-jquery/" target="_blank">jQuery regexp</a> syntax)</small></label><div class="input-group"><input class="form-control focused" id="tubeSelector" type="text" placeholder="prefix*" value="`)
	buf.WriteString(selfConf.TubeSelector)
	buf.WriteString(`"><div class="input-group-btn"><a href="javascript:void(0);" class="btn btn-info" id="clearTubesSelect">Select</a></div></div></div></fieldset><div><strong>Tube list</strong>`)
	buf.WriteString(tubeList.String())
	buf.WriteString(`</div></form></div><div class="modal-footer"><button type="button" class="btn btn-default" data-dismiss="modal">Close</button><a href="#" class="btn btn-success" id="clearTubes">Clear selected tubes</a><br/><br/><p class="text-muted text-right small">* Tube clear works by peeking to all jobs and deleting them in a loop.</p></div></div></div></div>`)
	return buf.String()
}
