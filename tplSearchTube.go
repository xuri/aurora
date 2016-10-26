package main

import (
	"bytes"
	"strconv"
)

// tplSearchTube rander navigation search box for search content in jobs by given tube.
func tplSearchTube(server string, tube string, state string) string {
	buf := bytes.Buffer{}
	buf.WriteString(`<form class="navbar-form navbar-right" style="margin-top:5px;margin-bottom:0px;" role="search" method="get"><input type="hidden" name="server" value="`)
	buf.WriteString(server)
	buf.WriteString(`"/><input type="hidden" name="tube" value="`)
	buf.WriteString(tube)
	buf.WriteString(`"/><input type="hidden" name="state" value="`)
	buf.WriteString(state)
	buf.WriteString(`"/><input type="hidden" name="action" value="search"/><input type="hidden" name="limit" value="`)
	buf.WriteString(strconv.Itoa(selfConf.SearchResultLimit))
	buf.WriteString(`"/><div class="form-group"><input type="text" class="form-control input-sm search-query" name="searchStr" placeholder="Search this tube"></div></form>`)
	return buf.String()
}
