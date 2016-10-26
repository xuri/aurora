package main

import (
	"fmt"
	"html"

	"github.com/kr/beanstalk"
)

// currentTubeJobsShowcase return a section include three stats of job,
// call currentTubeJobsShowcaseReadySection, currentTubeJobsShowcaseDelayedSection
// and currentTubeJobsShowcaseBuriedSection function and combine these return value
//  based on the given server and tube config.
func currentTubeJobsShowcase(server string, tube string) string {
	return fmt.Sprintf(`<section class="jobsShowcase">%s%s%s</section>`, currentTubeJobsShowcaseReadySection(server, tube), currentTubeJobsShowcaseDelayedSection(server, tube), currentTubeJobsShowcaseBuriedSection(server, tube))
}

// currentTubeJobsShowcaseReadySection constructs a tube job in ready
// stats table based on the given server and tube config.
func currentTubeJobsShowcaseReadySection(server string, tube string) string {
	var err error
	var statsJobStr, jobBodyStr, btnGroup, moveAllReadyTo, readyStats string
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "ready" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	tubes, _ := bstkConn.ListTubes()
	tubeStats := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	jobID, jobBody, err := tubeStats.PeekReady()
	if err != nil {
		bstkConn.Close()
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "ready" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}

	statsJob, err := bstkConn.StatsJob(jobID)
	if err != nil {
		bstkConn.Close()
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "ready" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	bstkConn.Close()

	for _, v := range jobStatsOrder {
		statsJobStr += fmt.Sprintf(`<tr><td>%s</td><td>%s</td></tr>`, v, statsJob[v])
	}

	for _, v := range tubes {
		moveAllReadyTo += fmt.Sprintf(`<li><a href="?server=%s&tube=%s&action=moveJobsTo&destTube=%s&state=ready">%s</a></li>`, server, tube, v, html.EscapeString(v))
	}
	if jobBody != nil {
		btnGroup = fmt.Sprintf(`<div class="pull-right">
                                    <div style="margin-bottom: 3px;">
                                        <a class="btn btn-sm btn-info addSample" data-jobid="%d"
                                           href="?server=%s&tube=%s&action=addSample"><i class="glyphicon glyphicon-plus glyphicon-white"></i> Add to
                                            samples</a>

                                        <div class="btn-group">
                                            <button class="btn btn-info btn-sm dropdown-toggle" data-toggle="dropdown">
                                                <i class="glyphicon glyphicon-arrow-right glyphicon-white"></i> Move all ready to
                                            </button>
                                            <ul class="dropdown-menu">
                                                <li><input class="moveJobsNewTubeName input-medium" type="text" data-href="?server=%s&tube=%s&action=moveJobsTo&state=ready&destTube=" placeholder="New tube name"/></li>
                                                    %s
                                                <li class="divider"></li>
                                                <li>
                                                    <a href="?server=%s&tube=%s&action=moveJobsTo&destState=buried&state=ready">Buried</a>
                                                </li>
                                            </ul>
                                        </div>
                                        <a class="btn btn-sm btn-danger"
                                           href="?server=%s&tube=%s&state=ready&action=deleteAll&count=1"
                                           onclick="return confirm('This process might hang a while on tubes with lots of jobs. Are you sure you want to continue?');"><i
                                                class="glyphicon glyphicon-trash glyphicon-white"></i> Delete all ready jobs</a>
                                        <a class="btn btn-sm btn-danger"
                                           href="?server=%s&tube=%s&state=ready&action=deleteJob&jobid=%d"><i
                                                class="glyphicon glyphicon-remove glyphicon-white"></i> Delete</a>
                                    </div>
                                </div>`, int(jobID), server, tube, server, tube, moveAllReadyTo, server, tube, server, tube, server, tube, int(jobID))
	}
	if jobBody != nil {
		jobBodyStr = preformat(jobBody)
	}

	if jobBody != nil {
		readyStats = fmt.Sprintf(`<hr>
            <div class="pull-left">
                <h3>Next job in "ready" state</h3>
            </div>
            <div class="clearfix"></div>
                <div class="row show-grid">
                    <div class="col-sm-3">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>Stats:</th>
                                    <th>&nbsp;</th>
                                </tr>
                            </thead>
                            <tbody>
                                %s
                            </tbody>
                        </table>
                    </div>
                    <div class="col-sm-9">
                        <div class="clearfix">
                            <div class="pull-left">
                                <b>Job data:</b>
                            </div>
                            %s
                        </div>
                        <pre><code>%s</code></pre>
                    </div>
                </div>`, statsJobStr, btnGroup, jobBodyStr)
	} else {
		readyStats = `<hr>
            <div class="pull-left">
                <h3>Next job in "ready" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	return readyStats
}

// currentTubeJobsShowcaseDelayedSection constructs a tube job in delayed
// stats table based on the given server and tube conf.
func currentTubeJobsShowcaseDelayedSection(server string, tube string) string {
	var err error
	var statsJobStr, jobBodyStr, btnGroup, moveAllDelayedTo, delayedStats string
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "delayed" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	tubes, _ := bstkConn.ListTubes()
	tubeStats := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	jobID, jobBody, err := tubeStats.PeekDelayed()
	if err != nil {
		bstkConn.Close()
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "delayed" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}

	statsJob, err := bstkConn.StatsJob(jobID)
	if err != nil {
		bstkConn.Close()
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "delayed" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	bstkConn.Close()

	for _, v := range jobStatsOrder {
		statsJobStr += fmt.Sprintf(`<tr><td>%s</td><td>%s</td></tr>`, v, statsJob[v])
	}

	for _, v := range tubes {
		moveAllDelayedTo += fmt.Sprintf(`<li><a href="?server=%s&tube=%s&action=moveJobsTo&destTube=%s&state=delayed">%s</a></li>`, server, tube, v, html.EscapeString(v))
	}
	if jobBody != nil {
		btnGroup = fmt.Sprintf(`<div class="pull-right">
                                    <div style="margin-bottom: 3px;">
                                        <a class="btn btn-sm btn-info addSample" data-jobid="%d"
                                           href="?server=%s&tube=%s&action=addSample"><i class="glyphicon glyphicon-plus glyphicon-white"></i> Add to
                                            samples</a>

                                        <div class="btn-group">
                                            <button class="btn btn-info btn-sm dropdown-toggle" data-toggle="dropdown">
                                                <i class="glyphicon glyphicon-arrow-right glyphicon-white"></i> Move all delayed to
                                            </button>
                                            <ul class="dropdown-menu">
                                                <li><input class="moveJobsNewTubeName input-medium" type="text" data-href="?server=%s&tube=%s&action=moveJobsTo&state=delayed&destTube=" placeholder="New tube name"/></li>
                                                    %s
                                                <li class="divider"></li>
                                                <li>
                                                    <a href="?server=%s&tube=%s&action=moveJobsTo&destState=buried&state=delayed">Buried</a>
                                                </li>
                                            </ul>
                                        </div>
                                        <a class="btn btn-sm btn-danger"
                                           href="?server=%s&tube=%s&state=delayed&action=deleteAll&count=1"
                                           onclick="return confirm('This process might hang a while on tubes with lots of jobs. Are you sure you want to continue?');"><i
                                                class="glyphicon glyphicon-trash glyphicon-white"></i> Delete all delayed jobs</a>
                                        <a class="btn btn-sm btn-danger"
                                           href="?server=%s&tube=%s&state=delayed&action=deleteJob&jobid=%d"><i
                                                class="glyphicon glyphicon-remove glyphicon-white"></i> Delete</a>
                                    </div>
                                </div>`, int(jobID), server, tube, server, tube, moveAllDelayedTo, server, tube, server, tube, server, tube, int(jobID))
	}
	if jobBody != nil {
		jobBodyStr = preformat(jobBody)
	}

	if jobBody != nil {
		delayedStats = fmt.Sprintf(`<hr>
            <div class="pull-left">
                <h3>Next job in "delayed" state</h3>
            </div>
            <div class="clearfix"></div>
                <div class="row show-grid">
                    <div class="col-sm-3">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>Stats:</th>
                                    <th>&nbsp;</th>
                                </tr>
                            </thead>
                            <tbody>
                                %s
                            </tbody>
                        </table>
                    </div>
                    <div class="col-sm-9">
                        <div class="clearfix">
                            <div class="pull-left">
                                <b>Job data:</b>
                            </div>
                            %s
                        </div>
                        <pre><code>%s</code></pre>
                    </div>
                </div>`, statsJobStr, btnGroup, jobBodyStr)
	} else {
		delayedStats = `<hr>
            <div class="pull-left">
                <h3>Next job in "delayed" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	return delayedStats
}

// currentTubeJobsShowcaseBuriedSection constructs a tube job in buried
// stats table based on the given server and tube config.
func currentTubeJobsShowcaseBuriedSection(server string, tube string) string {
	var err error
	var statsJobStr, jobBodyStr, btnGroup, moveAllBuriedTo, buriedStats string
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "buried" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	tubes, _ := bstkConn.ListTubes()
	tubeStats := &beanstalk.Tube{Conn: bstkConn, Name: tube}
	jobID, jobBody, err := tubeStats.PeekBuried()
	if err != nil {
		bstkConn.Close()
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "buried" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}

	statsJob, err := bstkConn.StatsJob(jobID)
	if err != nil {
		bstkConn.Close()
		return `<hr>
            <div class="pull-left">
                <h3>Next job in "buried" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	bstkConn.Close()

	for _, v := range jobStatsOrder {
		statsJobStr += fmt.Sprintf(`<tr><td>%s</td><td>%s</td></tr>`, v, statsJob[v])
	}

	for _, v := range tubes {
		moveAllBuriedTo += fmt.Sprintf(`<li><a href="?server=%s&tube=%s&action=moveJobsTo&destTube=%s&state=buried">%s</a></li>`, server, tube, v, html.EscapeString(v))
	}
	if jobBody != nil {
		btnGroup = fmt.Sprintf(`<div class="pull-right">
                                    <div style="margin-bottom: 3px;">
                                        <a class="btn btn-sm btn-info addSample" data-jobid="%d"
                                           href="?server=%s&tube=%s&action=addSample"><i class="glyphicon glyphicon-plus glyphicon-white"></i> Add to
                                            samples</a>

                                        <div class="btn-group">
                                            <button class="btn btn-info btn-sm dropdown-toggle" data-toggle="dropdown">
                                                <i class="glyphicon glyphicon-arrow-right glyphicon-white"></i> Move all buried to
                                            </button>
                                            <ul class="dropdown-menu">
                                                <li><input class="moveJobsNewTubeName input-medium" type="text" data-href="?server=%s&tube=%s&action=moveJobsTo&state=buried&destTube=" placeholder="New tube name"/></li>
                                                    %s
                                            </ul>
                                        </div>
                                        <a class="btn btn-sm btn-danger"
                                           href="?server=%s&tube=%s&state=buried&action=deleteAll&count=1"
                                           onclick="return confirm('This process might hang a while on tubes with lots of jobs. Are you sure you want to continue?');"><i
                                                class="glyphicon glyphicon-trash glyphicon-white"></i> Delete all buried jobs</a>
                                        <a class="btn btn-sm btn-danger"
                                           href="?server=%s&tube=%s&state=buried&action=deleteJob&jobid=%d"><i
                                                class="glyphicon glyphicon-remove glyphicon-white"></i> Delete</a>
                                    </div>
                                </div>`, int(jobID), server, tube, server, tube, moveAllBuriedTo, server, tube, server, tube, int(jobID))
	}
	if jobBody != nil {
		jobBodyStr = preformat(jobBody)
	}

	if jobBody != nil {
		buriedStats = fmt.Sprintf(`<hr>
            <div class="pull-left">
                <h3>Next job in "buried" state</h3>
            </div>
            <div class="clearfix"></div>
                <div class="row show-grid">
                    <div class="col-sm-3">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th>Stats:</th>
                                    <th>&nbsp;</th>
                                </tr>
                            </thead>
                            <tbody>
                                %s
                            </tbody>
                        </table>
                    </div>
                    <div class="col-sm-9">
                        <div class="clearfix">
                            <div class="pull-left">
                                <b>Job data:</b>
                            </div>
                            %s
                        </div>
                        <pre><code>%s</code></pre>
                    </div>
                </div>`, statsJobStr, btnGroup, jobBodyStr)
	} else {
		buriedStats = `<hr>
            <div class="pull-left">
                <h3>Next job in "ready" state</h3>
            </div>
            <div class="clearfix"></div><i>empty</i>`
	}
	return buriedStats
}
