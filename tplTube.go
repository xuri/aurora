package main

import (
	"fmt"
)

// tplTube render job stats page in tube by given server.
func tplTube(content string, server string, tube string) string {
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
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse"><span class="sr-only">Toggle navigation</span><span class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span></button><a class="navbar-brand" href="/">Beanstalk console</a></div>
            <div class="collapse navbar-collapse">
                <ul class="nav navbar-nav">
                    %s
                    %s
                </ul>
                <ul class="nav navbar-nav navbar-right">
                    <li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown">Toolbox <span class="caret"></span></a>
                        <ul class="dropdown-menu">
                            <li><a href="#filter" role="button" data-toggle="modal">Filter columns</a></li>
                            <li><a href="./sample?action=manageSamples" role="button">Manage samples</a></li>
                            <li class="divider"></li>
                            <li><a href="#settings" role="button" data-toggle="modal">Edit settings</a></li>
                        </ul>
                    </li>
                    %s
                </ul>
                %s
            </div>
        </div>
    </div>
    <div class="container">
        %s
        %s
        %s
        <div id='idAllTubesCopy' style="display:none"></div>
        %s
        %s
    </div>
    <script>
        function getParameterByName(name, url) {
            if (!url) url = window.location.href;
            name = name.replace(/[\[\]]/g, "\\$&");
            var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
                results = regex.exec(url);
            if (!results) return null;
            if (!results[2]) return '';
            return decodeURIComponent(results[2].replace(/\+/g, " "));
        }
        var url = "./tube?server="+getParameterByName('server');
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
</html>
`, TplHead, dropDownServer(server), dropDownTube(server, tube), TplLinks, tplSearchTube(server, tube, ""), content, modalAddJob(tube), modalAddSample(server, tube), tplTubeFilter(), dropEditSettings(), isDisabledJobDataHighlight)
}
