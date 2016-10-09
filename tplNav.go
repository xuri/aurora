package main

import (
	"fmt"

	"github.com/kr/beanstalk"
)

// getServerStatus render a server stats table.
func getServerStatus() string {
	var err error
	var td, th string
	for _, addr := range selfConf.Servers {
		var bstkConn *beanstalk.Conn
		if bstkConn, err = beanstalk.Dial("tcp", addr); err != nil {
			td += fmt.Sprintf(`<tr> <td>%s</td><td colspan="7" class="row-full">&nbsp;</td><td> <a class="btn btn-xs btn-danger" title="Remove from list" href="serversRemove?action=serversRemove&removeServer=%s"><span class="glyphicon glyphicon-minus"> </span></a> </td></tr>`, addr, addr)
			continue
		}
		s, _ := bstkConn.Stats()
		bstkConn.Close()

		td += fmt.Sprintf(`<tr><td><a href="/server?server=%s">%s</a></td>`, addr, addr)
		for _, v := range selfConf.Filter {
			td += fmt.Sprintf(`<td class="" name="%s">%s</td>`, v, s[v])
		}
		td += fmt.Sprintf(`<td><a class="btn btn-xs btn-danger" title="Remove from list" href="serversRemove?action=serversRemove&removeServer=%s"><span class="glyphicon glyphicon-minus"> </span></a></td></tr>`, addr)
	}
	for _, v := range selfConf.Filter {
		th += fmt.Sprintf(`<th class="" name="%s">%s</th>`, v, v)
	}
	template := fmt.Sprintf(`<div class="row"> <div class="col-sm-12"> <table class="table table-striped table-hover" id="servers-index"> <thead> <tr> <th>name</th> %s <th>&nbsp;</th> </tr></thead> <tbody>%s </tbody> </table> <a href="#servers-add" role="button" class="btn btn-info" id="addServer">Add server</a> </div></div>`, th, td)
	return template
}

// getServerTubes render a tubes stats table by given server.
func getServerTubes(server string) string {
	var err error
	var th, tr string
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		for _, v := range selfConf.TubeFilters {
			th += fmt.Sprintf(`<th name="%s">%s</th>`, v, v)
		}
		return fmt.Sprintf(`<div id="idAllTubes"><section id="summaryTable"> <div class="row"> <div class="col-sm-12"> <table class="table table-striped table-hover"> <thead> <tr> <th>name</th>%s</tr> </thead> <tbody> </tbody> </table> </div> </div> </section></div>`, th)
	}
	tubes, _ := bstkConn.ListTubes()
	for _, v := range selfConf.TubeFilters {
		th += fmt.Sprintf(`<th name="%s">%s</th>`, v, v)
	}
	for _, v := range tubes {
		tubeStats := &beanstalk.Tube{bstkConn, v}
		statsMap, err := tubeStats.Stats()
		if err != nil {
			continue
		}
		var td string
		for _, stats := range selfConf.TubeFilters {
			td += fmt.Sprintf(`<td>%s</td>`, statsMap[stats])
		}
		tr += fmt.Sprintf(`<tr><td name="pause-time-left"><a href="tube?server=%s&tube=%s">%s</a></td>%s</tr>`, server, v, v, td)
	}
	bstkConn.Close()
	template := fmt.Sprintf(`<div id="idAllTubes"><section id="summaryTable"> <div class="row"> <div class="col-sm-12"> <table class="table table-striped table-hover"> <thead> <tr> <th>name</th>%s</tr> </thead> <tbody> %s </tbody> </table> </div> </div> </section></div>`, th, tr)
	return template
}

// dropDownServer render a navigation dropdown menu for server list.
func dropDownServer(currentServer string) string {
	if currentServer == "" {
		currentServer = `All servers`
	}
	ul := fmt.Sprintf(`<li class="dropdown">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown">
                            %s <span class="caret"></span>
                        </a>
                        <ul class="dropdown-menu">`, currentServer)
	for _, addr := range selfConf.Servers {
		if addr == currentServer {
			continue
		}
		ul += fmt.Sprintf(`<li><a href="./server?server=%s">%s</a></li>`, addr, addr)
	}
	if currentServer != "All servers" {
		ul += `<li><a href="./public">All servers</a></li>`
	}
	ul += `</ul></li>`
	return ul
}

// dropDownTube render a navigation dropdown menu for tube list.
func dropDownTube(server string, currentTube string) string {
	if currentTube == "" {
		currentTube = `All tubes`
	}
	ul := fmt.Sprintf(`<li class="dropdown">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown">
                            %s <span class="caret"></span>
                        </a>
                        <ul class="dropdown-menu">`, currentTube)

	var bstkConn *beanstalk.Conn
	var err error
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		if currentTube != "" {
			ul += `<li><a href="./public">All tubes</a></li>`
		}
		ul += `</ul></li>`
		return ul
	}
	tubes, _ := bstkConn.ListTubes()
	for _, v := range tubes {
		ul += fmt.Sprintf(`<li><a href="./tube?server=%s&tube=%s">%s</a></li>`, server, v, v)
	}
	bstkConn.Close()
	if currentTube != "All tubes" {
		ul += fmt.Sprintf(`<li><a href="./server?server=%s">All tubes</a></li>`, server)
	}
	ul += `</ul></li>`
	return ul
}

// dropEditSettings render a navigation dropdown menu for set preference.
func dropEditSettings() string {
	var isDisabledJSONDecode, isDisabledJobDataHighlight string
	if selfConf.IsDisabledJSONDecode != 1 {
		isDisabledJSONDecode = `checked="checked"`
	}
	if selfConf.IsDisabledJobDataHighlight != 1 {
		isDisabledJobDataHighlight = `checked="checked"`
	}
	return fmt.Sprintf(`<div id="settings" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="settings-label" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                <h4 class="modal-title" id="settings-label">Settings</h4>
            </div>
            <div class="modal-body">
                <fieldset>
                    <div class="form-group">
                        <label for="tubePauseSeconds"><b>Tube pause seconds</b> (<i>-1</i> means the default: <i>3600</i>, <i>0</i> is reserved for
                            un-pause)</label>

                        <input class="form-control focused" id="tubePauseSeconds" type="text" value="%d">
                    </div>
                    <div class="form-group">
                        <label for="focusedInput"><b>Auto-refresh interval in milliseconds</b> (Default: <i>500</i>)</label>
                        <input class="form-control focused" id="autoRefreshTimeoutMs" type="text" value="%d">
                    </div>
                    <div class="form-group">
                        <label for="focusedInput"><b>Search result limits</b> (Default: <i>25</i>)</label>
                        <input class="form-control focused" id="searchResultLimit" type="text" value="%d">
                    </div>
                    <div class="form-group">
                        <label for="focusedInput"><b>Preferred way to deal with job data</b></label>

                        <div class="checkbox">
                            <label>
                                <input type="checkbox" id="isDisabledJsonDecode" value="1" %s>
                                before display: json_decode()
                            </label>
                        </div>

                        <div class="checkbox">
                            <label>
                                <input type="checkbox" id="isDisabledJobDataHighlight" value="1" %s>
                                after display: enable highlight
                            </label>
                        </div>

                    </div>
                </fieldset>
            </div>
            <div class="modal-footer">
                <button class="btn" data-dismiss="modal" aria-hidden="true">Close</button>
            </div>

        </div>
    </div>
</div>`, selfConf.TubePauseSeconds, selfConf.AutoRefreshTimeoutMs, selfConf.SearchResultLimit, isDisabledJSONDecode, isDisabledJobDataHighlight)
}
