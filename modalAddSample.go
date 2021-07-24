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

// modalAddSample render add sample popup window.
func modalAddSample(server string, tube string) string {
	var err error
	var tubeList, buf strings.Builder
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return ``
	}
	tubes, _ := bstkConn.ListTubes()
	for _, v := range tubes {
		var checked string
		if v == tube {
			checked = `checked="checked"`
		}
		tubeList.WriteString(`<div class="form-group"><div class="checkbox"><label><input type="checkbox" name="tubes[`)
		tubeList.WriteString(v)
		tubeList.WriteString(`]" value="1" `)
		tubeList.WriteString(checked)
		tubeList.WriteString(`>`)
		tubeList.WriteString(html.EscapeString(v))
		tubeList.WriteString(`</label></div></div>`)
	}
	buf.WriteString(`<div id="modalAddSample" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="addsamples-label" aria-hidden="true"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button><h4 id="addsamples-label" class="modal-title">Add to samples</h4></div><div class="modal-body"><input type="hidden" name="tube" value="`)
	buf.WriteString(tube)
	buf.WriteString(`"/><fieldset><div class="alert alert-danger" id="sampleSaveAlert" style="display: none;"><button type="button" class="close" onclick="$('#sampleSaveAlert').fadeOut('fast');">×</button><span><strong>Error!</strong> Required fields are marked *</span></div><input type="hidden" name="addsamplejobid" id="addsamplejobid"><div class="form-group"><label for="addsamplename" title="You can highlight text inside the job, then hit the Add button, it will be automatically populated here."><b>Name *</b><i>(highlighted text is auto populated)</i></label><input class="form-control focused" id="addsamplename" name="addsamplename" type="text" value="" autocomplete="off"></div></fieldset><div><label class="control-label"><b>Available on tubes *</b></label>`)
	buf.WriteString(tubeList.String())
	buf.WriteString(`</div></div><div class="modal-footer"><button class="btn" data-dismiss="modal" aria-hidden="true">Close</button><a href="#" class="btn btn-success" id="sampleSave">Save</a></div></div></div></div>`)
	return buf.String()
}
