package main

import "fmt"

// tplMain render server list.
func tplMain(serverList string, currentServer string) string {
	var isDisabledJobDataHighlight string
	if selfConf.IsDisabledJobDataHighlight != 1 {
		isDisabledJobDataHighlight = `<script src="./highlight/highlight.pack.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>`
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en-US">
%s
<body>
    <div class="navbar navbar-fixed-top navbar-default" role="navigation">
        <div class="container">
            <div class="navbar-header">
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
                    <span class="sr-only">Toggle navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
                <a class="navbar-brand" href="/">Beanstalk Console</a>
            </div>
            <div class="collapse navbar-collapse">
                <ul class="nav navbar-nav">
                    %s
                </ul>
                <ul class="nav navbar-nav navbar-right">
                    <li class="dropdown">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown">Toolbox <span class="caret"></span></a>
                        <ul class="dropdown-menu">
                            <li><a href="#filterServer" role="button" data-toggle="modal">Filter columns</a></li>
                            <li><a href="./sample?action=manageSamples" role="button">Manage samples</a></li>
                            <li class="divider"></li>
                            <li><a href="#settings" role="button" data-toggle="modal">Edit settings</a></li>
                        </ul>
                    </li>
                    %s
                    <li>
                        <button type="button" id="autoRefreshSummary" class="btn btn-default btn-small">
                            <span class="glyphicon glyphicon-refresh"></span>
                        </button>
                    </li>
                </ul>
            </div>
            <!--/.nav-collapse -->
        </div>
    </div>
    <div class="container">
        <div id="idServers">
        %s
        </div>
        %s
        <div id="idServersCopy" style="display:none"></div>
        <div id="servers-add" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="servers-add-label" aria-hidden="true">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                        <h4 class="modal-title" id="servers-add-labal">Add Server</h4>
                    </div>
                    <div class="modal-body">
                        <form class="form-horizontal">
                            <div class="form-group">
                                <label class="control-label col-sm-2" for="host">Host</label>
                                <div class="col-sm-10">
                                    <input type="text" id="host" value="localhost" class="form-control">
                                </div>
                            </div>
                            <div class="form-group">
                                <label class="control-label col-sm-2" for="port">Port</label>
                                <div class="col-sm-10">
                                    <input type="text" id="port" value="11300" class="form-control">
                                </div>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-info">Add server</button>
                        <button class="btn" data-dismiss="modal" aria-hidden="true">Cancel</button>
                    </div>
                </div>
            </div>
        </div>
        %s
        %s
    </div>
    <script>
        var url = "./index?server=";
        var contentType = "";
    </script>
    <script src='./assets/vendor/jquery/jquery.js'></script>
    <script src="./js/jquery.color.js"></script>
    <script src="./js/jquery.cookie.js"></script>
    <script src="./js/jquery.regexp.js"></script>
    <script src="./assets/vendor/bootstrap/js/bootstrap.min.js"></script>
    %s
    <script src="./js/customer.js"></script>
</body>
</html>`, TplHead, dropDownServer(currentServer), TplLinks, serverList, checkUpdate(), tplServerFilter(), dropEditSettings(), isDisabledJobDataHighlight)
}
