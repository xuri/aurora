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
	"strconv"
	"strings"
)

// modalAddJob render modal popup for add a job to tube.
func modalAddJob(tube string) string {
	buf := strings.Builder{}
	buf.WriteString(`<div class="modal fade" id="modalAddJob" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal">×</button><h4 class="modal-title">Add new job</h4></div><div class="modal-body"><form class="form-horizontal"><fieldset><div class="alert alert-danger" id="tubeSaveAlert" style="display: none;"><button type="button" class="close" onclick="$('#tubeSaveAlert').fadeOut('fast');">×</button><strong>Error!</strong> Required fields are marked * </div><div class="form-group"><label class="control-label col-xs-3">*Tube name</label><div class="col-xs-9"><input class="form-control focused" id="tubeName" type="text" value="`)
	buf.WriteString(tube)
	buf.WriteString(`"></div></div><div class="form-group"><label class="control-label col-xs-3">*Data</label><div class="col-xs-9"><textarea id="tubeData" rows="3" class="form-control"></textarea></div></div><div class="form-group"><label class="control-label col-xs-3">Priority</label><div class="col-xs-9"><input class="form-control focused" id="tubePriority" type="number" value="`)
	buf.WriteString(strconv.Itoa(DefaultPriority))
	buf.WriteString(`"></div></div><div class="form-group"><label class="control-label col-xs-3">Delay</label><div class="col-xs-9"><input class="form-control focused" id="tubeDelay" type="number" value="`)
	buf.WriteString(strconv.Itoa(DefaultDelay))
	buf.WriteString(`"></div></div><div class="form-group"><label class="control-label col-xs-3">TTR</label><div class="col-xs-9"><input class="form-control focused" id="tubeTtr" type="number" value="`)
	buf.WriteString(strconv.Itoa(DefaultTTR))
	buf.WriteString(`"></div></div><div class="modal-footer"><a href="#" class="btn" data-dismiss="modal">Close</a><a href="#" class="btn btn-success" id="tubeSave">Save changes</a></div></fieldset></form></div></div></div></div>`)
	return buf.String()
}
