// Copyright 2016 - 2021 The aurora Authors. All rights reserved. Use of this
// source code is governed by a MIT license that can be found in the LICENSE
// file.
//
// The aurora is a web-based beanstalkd queue server console written in Go
// and works on macOS, Linux and Windows machines. Main idea behind using Go
// for backend development is to utilize ability of the compiler to produce
// zero-dependency binaries for multiple platforms. aurora was created as an
// attempt to build very simple and portable application to work with local or
// remote beanstalkd server.

package main

import (
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/aurora/beanstalk"
)

// addJob puts a job into tube by given config.
func addJob(server string, tube string, data string, priority string, delay string, TTR string) {
	var (
		err                error
		tubeDelay, tubeTTR int
		tubePriority       uint64
		bstkConn           *beanstalk.Conn
	)
	tubePriority, err = strconv.ParseUint(priority, 10, 32)
	if err != nil || tubePriority > math.MaxUint32 {
		tubePriority = DefaultPriority
	}
	tubeDelay, err = strconv.Atoi(delay)
	if err != nil {
		tubeDelay = DefaultDelay
	}
	tubeTTR, err = strconv.Atoi(TTR)
	if err != nil {
		tubeTTR = DefaultTTR
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	_, _ = bstkTube.Put([]byte(data), uint32(tubePriority), time.Duration(tubeDelay)*time.Second, time.Duration(tubeTTR)*time.Second)
	bstkConn.Close()
}

// deleteJob delete a job in tube by given config.
func deleteJob(server string, tube string, jobID string) {
	var err error
	var id int
	var bstkConn *beanstalk.Conn
	id, err = strconv.Atoi(jobID)
	if err != nil {
		return
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	_ = bstkConn.Delete(uint64(id))
	bstkConn.Close()
}

// deleteAll delete all jobs in tube by given server and tube.
func deleteAll(server string, tube string) {
	var err error
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	for {
		readyJob, _, err := bstkTube.PeekReady()
		if err != nil {
			break
		}
		_ = bstkConn.Delete(readyJob)
	}
	for {
		buriedJob, _, err := bstkTube.PeekBuried()
		if err != nil {
			break
		}
		_ = bstkConn.Delete(buriedJob)
	}
	for {
		delayedJob, _, err := bstkTube.PeekDelayed()
		if err != nil {
			break
		}
		_ = bstkConn.Delete(delayedJob)
	}
	bstkConn.Close()
}

// kick takes up to bound jobs from the holding area and moves them into the
// ready queue, then returns the number of jobs moved. Jobs will be taken in the
// order in which they were last buried.
func kick(server string, tube string, count string) {
	var err error
	var bound int
	var bstkConn *beanstalk.Conn
	bound, err = strconv.Atoi(count)
	if err != nil {
		bound = 0
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	_, _ = bstkTube.Kick(bound)
	bstkConn.Close()
}

// kickJob kick single job in tube by given server, tube name and job ID.
func kickJob(server string, tube string, id string) {
	var err error
	var bound int
	var bstkConn *beanstalk.Conn
	bound, err = strconv.Atoi(id)
	if err != nil {
		return
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	_ = bstkConn.KickJob(uint64(bound))
	bstkConn.Close()
}

// pause pauses new reservations in tube for time duration.
func pause(server string, tube string, count string) {
	var err error
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	switch count {
	case "-1": // Pause tube
		if selfConf.TubePauseSeconds == -1 {
			_ = bstkTube.Pause(DefaultTubePauseSeconds * time.Second)
		} else {
			_ = bstkTube.Pause(time.Duration(selfConf.TubePauseSeconds) * time.Second)
		}
	case "0":
		_ = bstkTube.Pause(0 * time.Second) // Unpause tube
	}
	bstkConn.Close()
}

// moveJobsTo switch two case when move a job.
func moveJobsTo(server string, tube string, destTube string, state string, destState string) {
	switch state {
	case "ready": // ready to buried or ready
		moveReadyJobsTo(server, tube, destTube, destState)
	case "buried": // move job across the tube
		moveBuriedJobsTo(server, tube, destTube, destState)
	}
}

// moveReadyJobsTo process job moved origin stats in ready.
func moveReadyJobsTo(server string, tube string, destTube string, destState string) {
	var err error
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	switch destState {
	case "buried":
		tubeSet := beanstalk.NewTubeSet(bstkConn, tube)
		for {
			id, _, err := tubeSet.Reserve(time.Second)
			if err != nil {
				break
			}
			err = bstkConn.Bury(id, DefaultPriority)
			if err != nil {
				break
			}
		}
	default:
		if tube == destTube {
			bstkConn.Close()
			return
		}
		bstkDestTube := &beanstalk.Tube{
			Conn: bstkConn,
			Name: destTube,
		}
		for {
			readyJob, readyBody, err := bstkTube.PeekReady()
			if err != nil {
				break
			}
			_, err = bstkDestTube.Put(readyBody, DefaultPriority, DefaultDelay, DefaultTTR)
			if err != nil {
				break
			}
			err = bstkConn.Delete(readyJob)
			if err != nil {
				break
			}
		}
	}
	bstkConn.Close()
}

// moveBuriedJobsTo process job moved origin stats in buried.
func moveBuriedJobsTo(server string, tube string, destTube string, destState string) {
	var err error
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	bstkDestTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: destTube,
	}
	for {
		buriedJob, buriedBody, err := bstkTube.PeekBuried()
		if err != nil {
			break
		}
		_, err = bstkDestTube.Put(buriedBody, DefaultPriority, DefaultDelay, DefaultTTR)
		if err != nil {
			break
		}
		err = bstkConn.Delete(buriedJob)
		if err != nil {
			break
		}
	}
	bstkConn.Close()
}

// clearTubes delete all jobs in all tubes by given server.
func clearTubes(server string, data url.Values) {
	for tube := range data { // range over map
		deleteAll(server, tube)
	}
}

// searchTube search job by given search string in ready, delayed and buried
// stats.
func searchTube(server string, tube string, limit string, searchStr string) string {
	var (
		bstkConn      *beanstalk.Conn
		bstkConnStats map[string]string
		err           error
		result        = []SearchResult{}
		searchLimit   int
		statsFilter   = []string{"ready", "delayed", "buried"}
		table         = currentTubeJobsSummaryTable(server, tube)
		totalJobs, id uint64
	)
	if table == `` {
		return `Tube "` + tube + `" not found or it is empty <br><br><a href="./server?server=` + server + `"> &lt;&lt; back </a>`
	}
	searchLimit, err = strconv.Atoi(limit)
	if err != nil {
		return table
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return table
	}
	if bstkConnStats, err = bstkConn.Stats(); err != nil {
		return table
	}
	if totalJobs, err = strconv.ParseUint(bstkConnStats["total-jobs"], 10, 64); err != nil {
		return table
	}
	// Get ready stat job total
	for _, state := range statsFilter {
		var cnt int
		for id = totalJobs; id > 0; id-- {
			if cnt >= searchLimit {
				continue
			}
			ret := searchTubeInStats(tube, searchStr, state, bstkConn, id)
			if ret != nil {
				result = append(result, *ret)
				cnt++
			}
		}
	}
	bstkConn.Close()
	return table + currentTubeSearchResults(server, tube, limit, searchStr, result)
}

// searchTubeInStats search job in tube by given stats.
func searchTubeInStats(tube, searchStr, stat string, bstkConn *beanstalk.Conn, id uint64) *SearchResult {
	jobStats, err := bstkConn.StatsJob(id)
	if err != nil {
		return nil
	}
	if jobStats["tube"] != tube || jobStats["state"] != stat {
		return nil
	}
	readyBody, err := bstkConn.Peek(id)
	if err != nil {
		return nil
	}
	body := string(readyBody)
	if !strings.Contains(body, searchStr) {
		return nil
	}
	return &SearchResult{
		ID:    id,
		State: stat,
		Data:  string(readyBody),
	}
}
