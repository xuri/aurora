package main

import (
	"fmt"
	"html"
)

// currentTubeSearchResults constructs a search result table by given server, tube, search result limit and search content.
func currentTubeSearchResults(server string, tube string, limit string, searchStr string, result []SearchResult) string {
	if len(result) == 0 {
		return fmt.Sprintf(`<br/>No results found for <b>%s</b> in tube: <b>%s</b>`, html.EscapeString(searchStr), tube)
	}
	var tr string
	for _, job := range result {
		tr += fmt.Sprintf(`<tr><td>%d</td>
                                    <td>%s</td>
                                    <td class="ellipsize">%s</td>
                                    <td>
                                        <div class="dropdown btn-group-xs">
                                            <button class="btn btn-default dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown" aria-expanded="true">
                                                Actions
                                                <span class="caret"></span>
                                            </button>
                                            <ul class="dropdown-menu" role="menu" aria-labelledby="dropdownMenu1">
                                                <li role="presentation"><a role="menuitem" class="addSample" data-jobid="%d"
                                                                           href="?server=%s&tube=%s&action=addSample">
                                                        <i class="glyphicon glyphicon-plus glyphicon-white"></i>
                                                        Add to samples</a>
                                                </li>
                                                <li role="presentation"><a role="menuitem"
                                                                           href="?server=%s&tube=%s&state=%s&action=deleteJob&jobid=%d"><i
                                                            class="glyphicon glyphicon-remove glyphicon-white"></i>
                                                        Delete</a>
                                                </li>
                                                <li role="presentation"><a role="menuitem"
                                                                           href="?server=%s&tube=%s&state=%s&action=kickJob&jobid=%d"><i
                                                            class="glyphicon glyphicon-forward glyphicon-white"></i>
                                                        Kick</a>
                                                </li>
                                            </ul>
                                        </div>
                                    </td></tr>`, int(job.ID), job.State, html.EscapeString(job.Data), int(job.ID), server, tube, server, tube, job.State, int(job.ID), server, tube, job.State, int(job.ID))
	}
	return fmt.Sprintf(`<section id="actionsRow">
    <a class="btn btn-default btn-sm" href="?server=%s&tube=%s"><i class="glyphicon glyphicon-backward"></i>  &nbsp;Back to tube</a>
</section>
    <section id="searchResult">
        <div class="row">
            <div class="col-sm-12">
                <table class="table table-striped table-hover" style="table-layout:fixed;">
                    <thead>
                        <tr>
                            <th class="col-md-1">id</th>
                            <th class="col-md-1">state</th>
                            <th>data</th>
                            <th class="col-md-1">action</th>
                        </tr>
                    </thead>
                    <tbody>
                        %s
                    </tbody>
                </table>
            </div>
        </div>
        First %s rows are displayed for each state.
        <br/>
        <br/>
    </section>`, server, tube, tr, limit)
}
