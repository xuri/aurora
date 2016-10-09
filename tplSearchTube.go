package main

import (
	"fmt"
)

// tplSearchTube rander navigation search box for search content in jobs by given tube.
func tplSearchTube(server string, tube string, state string) string {
	return fmt.Sprintf(`<form  class="navbar-form navbar-right" style="margin-top:5px;margin-bottom:0px;" role="search" action="" method="get">
                                <input type="hidden" name="server" value="%s"/>
                                <input type="hidden" name="tube" value="%s"/>
                                <input type="hidden" name="state" value="%s"/>
                                <input type="hidden" name="action" value="search"/>
                                <input type="hidden" name="limit" value="%d"/>
                                <div class="form-group">
                                    <input type="text" class="form-control input-sm search-query" name="searchStr" placeholder="Search this tube">
                                </div>
                            </form>`, server, tube, state, selfConf.SearchResultLimit)
}
