package main

import (
	"bytes"
	"html"
	"net/url"
	"strconv"

	"github.com/Luxurioust/aurora/beanstalk"
)

// currentTubeJobsActionsRow render a section include kick, pause and unpause
// job button by given server and tube.
func currentTubeJobsActionsRow(server string, tube string) string {
	var err error
	var bstkConn *beanstalk.Conn
	var buf, pauseTimeLeft bytes.Buffer
	var pause = strconv.Itoa(selfConf.TubePauseSeconds)
	if pause == "-1" {
		pause = "3600"
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return ``
	}
	tubeStats := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	statsMap, _ := tubeStats.Stats()
	if statsMap["pause-time-left"] == "0" {
		pauseTimeLeft.WriteString(`<a class="btn btn-default btn-sm" href="?server=`)
		pauseTimeLeft.WriteString(server)
		pauseTimeLeft.WriteString(`&tube=`)
		pauseTimeLeft.WriteString(url.QueryEscape(tube))
		pauseTimeLeft.WriteString(`&action=pause&count=-1" title="Temporarily prevent jobs being reserved from the given tube. Pause for: `)
		pauseTimeLeft.WriteString(pause)
		pauseTimeLeft.WriteString(` seconds"><i class="glyphicon glyphicon-pause"></i> Pause tube</a>`)
	} else {
		pauseTimeLeft.WriteString(`<a class="btn btn-default btn-sm" href="?server=`)
		pauseTimeLeft.WriteString(server)
		pauseTimeLeft.WriteString(`&tube=`)
		pauseTimeLeft.WriteString(url.QueryEscape(tube))
		pauseTimeLeft.WriteString(`&action=pause&count=0" title="Pause seconds left: `)
		pauseTimeLeft.WriteString(statsMap["pause-time-left"])
		pauseTimeLeft.WriteString(`"><i class="glyphicon glyphicon-play"></i> Unpause tube</a>`)
	}
	bstkConn.Close()
	buf.WriteString(`<section id="actionsRow"><b>Actions:</b> &nbsp;<a class="btn btn-default btn-sm" href="?server=`)
	buf.WriteString(server)
	buf.WriteString(`&tube=`)
	buf.WriteString(url.QueryEscape(tube))
	buf.WriteString(`&action=kick&count=1"><i class="glyphicon glyphicon-forward"></i> Kick 1 job</a> <a class="btn btn-default btn-sm" href="?server=`)
	buf.WriteString(server)
	buf.WriteString(`&tube=`)
	buf.WriteString(url.QueryEscape(tube))
	buf.WriteString(`&action=kick&count=10" title='To kick more jobs, edit the "count" parameter'><i class="glyphicon glyphicon-fast-forward"></i> Kick 10 job</a> `)
	buf.WriteString(pauseTimeLeft.String())
	buf.WriteString(` &nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <div class="btn-group"><a data-toggle="modal" class="btn btn-success btn-sm" href="#" id="addJob"><i class="glyphicon glyphicon-plus-sign glyphicon-white"></i> Add job</a><button class="btn btn-success btn-sm dropdown-toggle" data-toggle="dropdown"><span class="caret"></span></button><ul class="dropdown-menu">`)
	buf.WriteString(currentTubeJobsActionsRowSample(server, tube))
	buf.WriteString(`</ul></div></section>`)
	return buf.String()
}

// currentTubeJobsActionsRowSample render a dropdown sample list by given server
// and tube.
func currentTubeJobsActionsRowSample(server string, tube string) string {
	sample := bytes.Buffer{}
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
				sample.WriteString(`<li><a href="?server=`)
				sample.WriteString(server)
				sample.WriteString(`&tube=`)
				sample.WriteString(url.QueryEscape(tube))
				sample.WriteString(`&action=loadSample&key=`)
				sample.WriteString(j.Key)
				sample.WriteString(`">`)
				sample.WriteString(html.EscapeString(j.Name))
				sample.WriteString(`</a></li>`)
			}
		}
	}
	if sample.String() == "" {
		return `<li><a href="javascript:void(0);">There are no sample jobs</a></li>`
	}
	sample.WriteString(`<li class="divider"></li><li><a href="./sample?action=manageSamples">Manage samples</a></li>`)
	return sample.String()
}
