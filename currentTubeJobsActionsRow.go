package main

import (
	"fmt"
	"html"

	"github.com/kr/beanstalk"
)

func currentTubeJobsActionsRow(server string, tube string) string {
	var err error
	var pauseTimeLeft string
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return ``
	}
	tubeStats := &beanstalk.Tube{bstkConn, tube}
	statsMap, _ := tubeStats.Stats()
	if statsMap["pause-time-left"] == "0" {
		pauseTimeLeft = fmt.Sprintf(`<a class="btn btn-default btn-sm" href="?server=%s&tube=%s&action=pause&count=-1"
           title="Temporarily prevent jobs being reserved from the given tube. Pause for: %d seconds"><i class="glyphicon glyphicon-pause"></i>
            Pause tube</a>`, server, tube, selfConf.TubePauseSeconds)
	} else {
		pauseTimeLeft = fmt.Sprintf(`<a class="btn btn-default btn-sm" href="?server=%s&tube=%s&action=pause&count=0"
           title="Pause seconds left: %d"><i class="glyphicon glyphicon-play"></i> Unpause tube</a>`, server, tube, statsMap["pause-time-left"])
	}
	bstkConn.Close()

	return fmt.Sprintf(`<section id="actionsRow">
    <b>Actions:</b>&nbsp;
    <a class="btn btn-default btn-sm" href="?server=%s&tube=%s&action=kick&count=1"><i class="glyphicon glyphicon-forward"></i> Kick 1 job</a>
    <a class="btn btn-default btn-sm" href="?server=%s&tube=%s&action=kick&count=10"
       title="To kick more jobs, edit the "count" parameter"><i class="glyphicon glyphicon-fast-forward"></i> Kick 10 job</a>
       %s
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    <div class="btn-group">
        <a data-toggle="modal" class="btn btn-success btn-sm" href="#" id="addJob"><i class="glyphicon glyphicon-plus-sign glyphicon-white"></i> Add job</a>
        <button class="btn btn-success btn-sm dropdown-toggle" data-toggle="dropdown">
            <span class="caret"></span>
        </button>
        <ul class="dropdown-menu">
            %s
        </ul>
    </div>
</section>`, server, tube, server, tube, pauseTimeLeft, currentTubeJobsActionsRowSample(server, tube))
}

func currentTubeJobsActionsRowSample(server string, tube string) string {
	var sample string
	for _, v := range sampleJobs.Tubes {
		if v.Name != tube {
			continue
		}
		if len(v.Keys) == 0 {
			continue
		}
		for _, k := range v.Keys {
			for _, j := range sampleJobs.Jobs {
				if j.Key != k {
					continue
				}
				sample += fmt.Sprintf(`<li><a href="?server=%s&tube=%s&action=loadSample&key=%s">%s</a></li>`, server, tube, j.Key, html.EscapeString(j.Name))
			}
		}
	}
	if sample == "" {
		return `<li><a href="#">There are no sample jobs</a></li>`
	}
	return sample + `<li class="divider"></li><li><a href="./sample?action=manageSamples">Manage samples</a></li>`
}
