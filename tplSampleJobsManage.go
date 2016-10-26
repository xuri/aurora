package main

import (
	"bytes"
)

// tplSampleJobsManage render sample job manage list page.
func tplSampleJobsManage(sampleList string, currentServer string) string {
	var isDisabledJobDataHighlight string
	if selfConf.IsDisabledJobDataHighlight != 1 {
		isDisabledJobDataHighlight = `<script src="./highlight/highlight.pack.js"></script><script>hljs.initHighlightingOnLoad();</script>`
	}
	buf := bytes.Buffer{}
	buf.WriteString(`<!DOCTYPE html><html lang="en-US">`)
	buf.WriteString(TplHead)
	buf.WriteString(`<body>`)
	buf.WriteString(TplNoScript)
	buf.WriteString(`<div class="navbar navbar-fixed-top navbar-default" role="navigation"><div class="container"><div class="navbar-header"><button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse"><span class="sr-only">Toggle navigation</span><span class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span></button><a class="navbar-brand" href="/">Beanstalk Console</a></div><div class="collapse navbar-collapse"><ul class="nav navbar-nav">`)
	buf.WriteString(dropDownServer(currentServer))
	buf.WriteString(`</ul><ul class="nav navbar-nav navbar-right"><li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown">Toolbox <span class="caret"></span></a><ul class="dropdown-menu"><li><a href="?action=manageSamples" role="button">Manage samples</a></li><li class="divider"></li><li><a href="#settings" role="button" data-toggle="modal">Edit settings</a></li></ul></li>`)
	buf.WriteString(TplLinks)
	buf.WriteString(`</ul></div><!--/.nav-collapse --></div></div><div class="container">`)
	buf.WriteString(sampleList)
	buf.WriteString(dropEditSettings())
	buf.WriteString(`</div><script>var url = "./sample"; var contentType = "";</script><script src='./assets/vendor/jquery/jquery.js'></script><script src="./js/jquery.color.js"></script><script src="./js/jquery.cookie.js"></script><script src="./js/jquery.regexp.js"></script><script src="./assets/vendor/bootstrap/js/bootstrap.min.js"></script>`)
	buf.WriteString(isDisabledJobDataHighlight)
	buf.WriteString(`<script src="./js/customer.js"></script></body></html>`)
	return buf.String()
}
