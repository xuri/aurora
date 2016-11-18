package main

import "bytes"

// modalAddJob render modal popup for add a job to tube.
func modalAddJob(tube string) string {
	buf := bytes.Buffer{}
	buf.WriteString(`<div class="modal fade" id="modalAddJob" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button type="button" class="close" data-dismiss="modal">×</button><h4 class="modal-title">Add new job</h4></div><div class="modal-body"><form class="form-horizontal"><fieldset><div class="alert alert-danger" id="tubeSaveAlert" style="display: none;"><button type="button" class="close" onclick="$('#tubeSaveAlert').fadeOut('fast');">×</button><strong>Error!</strong> Required fields are marked * </div><div class="form-group"><label class="control-label col-xs-3">*Tube name</label><div class="col-xs-9"><input class="form-control focused" id="tubeName" type="text" value="`)
	buf.WriteString(tube)
	buf.WriteString(`"></div></div><div class="form-group"><label class="control-label col-xs-3">*Data</label><div class="col-xs-9"><textarea id="tubeData" rows="3" class="form-control"></textarea></div></div><div class="form-group"><label class="control-label col-xs-3">Priority</label><div class="col-xs-9"><input class="form-control focused" id="tubePriority" type="text" value=""></div></div><div class="form-group"><label class="control-label col-xs-3">Delay</label><div class="col-xs-9"><input class="form-control focused" id="tubeDelay" type="text" value=""></div></div><div class="form-group"><label class="control-label col-xs-3">Ttr</label><div class="col-xs-9"><input class="form-control focused" id="tubeTtr" type="text" value=""></div></div><div class="modal-footer"><a href="#" class="btn" data-dismiss="modal">Close</a><a href="#" class="btn btn-success" id="tubeSave">Save changes</a></div></fieldset></form></div></div></div></div>`)
	return buf.String()
}
