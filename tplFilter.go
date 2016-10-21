package main

import (
	"fmt"
)

// tplServerFilterStatsGroups render server filter stats groups checkbox.
func tplServerFilterStatsGroups() []string {
	stats := []string{"", "", "", ""}
	statsGroupsFilter := [][]map[string]string{binlogStatsGroups, cmdStatsGroups, currentStatsGroups, otherStatsGroups}
	for k, statsGroups := range statsGroupsFilter {
		for _, statsGroup := range statsGroups {
			for property, description := range statsGroup {
				status := ""
				if checkInSlice(selfConf.Filter, property) {
					status = `checked="true"`
				}
				stats[k] += fmt.Sprintf(`<div class="control-group">
                 <div class="controls">
                     <div class="checkbox">
                         <label>
                             <input type="checkbox" name="%s" %s>
                             <b>%s</b>
                             <br/>%s</label>
                     </div>
                 </div>
             </div>`, property, status, property, description)
			}
		}
	}
	return stats
}

// tplServerFilter render modal popup for select server tube stats column.
func tplServerFilter() string {
	stats := tplServerFilterStatsGroups()
	filter := fmt.Sprintf(`<div id="filterServer" data-cookie="filter" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="servers-add-label" aria-hidden="true">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                        <h3 id="filter-label" class="text-info">Filter columns</h3>
                    </div>
                    <div class="modal-body">
                        <form class="form-group">
                            <div class="tabbable">
                                <ul class="nav nav-tabs">
                                    <li class="active"><a href="#binlog" data-toggle="tab">binlog</a></li>
                                    <li><a href="#cmd" data-toggle="tab">cmd</a></li>
                                    <li><a href="#current" data-toggle="tab">current</a></li>
                                    <li><a href="#other" data-toggle="tab">other</a></li>
                                </ul>
                                <div class="tab-content">
                                    <div class="tab-pane active" id="binlog">
                                        %s
                                    </div>
                                    <div class="tab-pane" id="cmd">
                                        %s
                                    </div>
                                    <div class="tab-pane" id="current">
                                        %s
                                    </div>
                                    <div class="tab-pane" id="other">
                                        %s
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn" data-dismiss="modal" aria-hidden="true">Close</button>
                    </div>
                </div>
            </div>
        </div>`, stats[0], stats[1], stats[2], stats[3])
	return filter
}

// tplTubeFilter render a modal popup for select job stats of tube.
func tplTubeFilter() string {
	var currents, others string
	for k, current := range tubeStatFields {
		if k > 7 {
			continue
		}
		for property, description := range current {
			status := ""
			if checkInSlice(selfConf.TubeFilters, property) {
				status = `checked="true"`
			}
			currents += fmt.Sprintf(`<div class="form-group">
                    <div class="checkbox">
                        <label class="checkbox">
                            <input type="checkbox" name="%s" %s><b>%s</b>
                            <br/>%s</label>
                    </div>
                </div>`, property, status, property, description)
		}
	}

	for k, other := range tubeStatFields {
		if k < 8 {
			continue
		}
		for property, description := range other {
			status := ""
			if checkInSlice(selfConf.TubeFilters, property) {
				status = `checked="true"`
			}
			others += fmt.Sprintf(`<div class="form-group">
                    <div class="checkbox">
                        <label class="checkbox">
                            <input type="checkbox" name="%s" %s><b>%s</b>
                            <br/>%s</label>
                    </div>
                </div>`, property, status, property, description)
		}
	}

	return fmt.Sprintf(`<div id="filter" data-cookie="tubefilter" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="filter-label" aria-hidden="true">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                        <h4 class="modal-title" id="filter-label">Filter columns</h4></div>
                    <div class="modal-body">
                        <form>
                            <div class="tabbable">
                                <ul class="nav nav-tabs">
                                    <li class="active"><a href="#current" data-toggle="tab">current</a></li>
                                    <li><a href="#other" data-toggle="tab">other</a></li>
                                </ul>
                                <div class="tab-content">
                                    <div class="tab-pane active" id="current">
                                        %s
                                    </div>
                                    <div class="tab-pane" id="other">
                                        %s
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-success" data-dismiss="modal" aria-hidden="true">Close</button>
                    </div>
                </div>
            </div>
        </div>`, currents, others)
}
