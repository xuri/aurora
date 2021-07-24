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
	"container/list"
	"io"
	"net/http"
	"os"
	"sync"
)

// Define the default configuration and HTML header template.
const (
	ConfigFileTemplate      = "servers = []\r\nlisten = \"127.0.0.1:3000\"\r\nversion = 2.2\r\n\r\n[openpage]\r\nenabled = true\r\n\r\n[auth]\r\nenabled = false\r\npassword = \"password\"\r\nusername = \"admin\"\r\n\r\n[sample]\r\nstorage = \"{}\""
	DefaultDelay            = 0
	DefaultPriority         = 1024 // most urgent: 0, least urgent: 4294967295.
	DefaultTTR              = 60   // 1 minute
	DefaultTubePauseSeconds = 3600
	TplHeaderBegin          = `<!DOCTYPE html><html lang="en-US"><head><meta charset="UTF-8"/><meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"><!--[if IE]><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><![endif]--><meta name="description" content="Beanstalkd Console"><meta name="keywords" content="Beanstalkd Console, beanstalkd, console"><meta content="always" name="referrer"><meta name="language" content="en-US"><meta name="category" content="Tools"><meta name="summary" content="Beanstalkd Console"><meta name="apple-mobile-web-app-capable" content="yes"/><link rel="copyright" href="http://www.opensource.org/licenses/mit-license.php"/><link rel="icon" sizes="32x32" href="./images/aurora-32x32.ico"><link rel="apple-touch-icon" sizes="180x180" href="./images/apple-touch-icon-180x180-precomposed.png"><link rel="apple-touch-icon" sizes="152x152" href="./images/apple-touch-icon-152x152-precomposed.png"><link rel="apple-touch-icon" sizes="144x144" href="./images/apple-touch-icon-144x144-precomposed.png"><link rel="apple-touch-icon" sizes="120x120" href="./images/apple-touch-icon-120x120-precomposed.png"><link rel="apple-touch-icon" sizes="114x114" href="./images/apple-touch-icon-114x114-precomposed.png"><link rel="apple-touch-icon" sizes="76x76" href="./images/apple-touch-icon-76x76-precomposed.png"><link rel="apple-touch-icon" sizes="72x72" href="./images/apple-touch-icon-72x72-precomposed.png"><link rel="apple-touch-icon" href="./images/apple-touch-icon-precomposed-57x57.png"><title>`
	TplHeaderEnd            = ` Beanstalkd Console</title><!-- Bootstrap core CSS --><link href="./assets/vendor/bootstrap/css/bootstrap.min.css" rel="stylesheet"><link href="./css/customer.css" rel="stylesheet"><link href="./highlight/styles/magula.css" rel="stylesheet"><!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries --><!--[if lt IE 9]><script src="./js/libs/html5shiv/3.7.0/html5shiv.js"></script><script src="./js/libs/respond.js/1.4.2/respond.min.js"></script><![endif]--></head><body>`
	TplLinks                = `<li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown"> Links <span class="caret"></span></a><ul class="dropdown-menu"><li><a href="https://github.com/beanstalkd" target="_blank">Beanstalkd (GitHub)</a></li><li><a href="https://github.com/xuri/aurora" target="_blank">Aurora (GitHub)</a></li></ul></li>`
	TplNoScript             = `<noscript><div class="container"><div class="alert alert-danger" role="alert">Aurora beanstalkd console requires JavaScript supports, please refresh after enable browser JavaScript support.</div></div></noscript>`
	UpdateURL               = `https://api.github.com/repos/xuri/aurora/tags`
	Version                 = 2.2
)

// Define server and tube stats fields.
var (
	pubConf              PubConfig
	sampleJobs           SampleJobs
	selfConf             SelfConf
	Stderr               io.Writer = os.Stderr // Stderr is the io.Writer to which executed commands write standard error.
	Stdout               io.Writer = os.Stdout // Stdout is the io.Writer to which executed commands write standard output.
	ConfigFile                     = `aurora.toml`
	statisticsData                 = StatisticsData{new(sync.RWMutex), statisticsDataServer}
	statisticsDataServer           = make(map[string]map[string]map[string]*list.List)
	notify                         = make(chan bool, 1)
	updateInfo                     = "uncheck"
	// Server filter columns.
	binlogStatsGroups = []map[string]string{
		{"binlog-current-index": "the index of the current binlog file being written to. If binlog is not active this value will be 0"},
		{"binlog-max-size": "the maximum size in bytes a binlog file is allowed to get before a new binlog file is opened"},
		{"binlog-oldest-index": "the index of the oldest binlog file needed to store the current jobs"},
		{"binlog-records-migrated": "the cumulative number of records written as part of compaction"},
		{"binlog-records-written": "the cumulative number of records written to the binlog"},
	}
	cmdStatsGroups = []map[string]string{
		{"cmd-bury": "the cumulative number of bury commands"},
		{"cmd-delete": "the cumulative number of delete commands"},
		{"cmd-ignore": "the cumulative number of ignore commands"},
		{"cmd-kick": "the cumulative number of kick commands"},
		{"cmd-list-tube-used": "the cumulative number of list-tube-used commands"},
		{"cmd-list-tubes": "the cumulative number of list-tubes commands"},
		{"cmd-list-tubes-watched": "the cumulative number of list-tubes-watched commands"},
		{"cmd-pause-tube": "the cumulative number of pause-tube commands"},
		{"cmd-peek": "the cumulative number of peek commands"},
		{"cmd-peek-buried": "the cumulative number of peek-buried commands"},
		{"cmd-peek-delayed": "the cumulative number of peek-delayed commands"},
		{"cmd-peek-ready": "the cumulative number of peek-ready commands"},
		{"cmd-put": "the cumulative number of put commands"},
		{"cmd-release": "the cumulative number of release commands"},
		{"cmd-reserve": "the cumulative number of reserve commands"},
		{"cmd-stats": "the cumulative number of stats commands"},
		{"cmd-stats-job": "the cumulative number of stats-job commands"},
		{"cmd-stats-tube": "the cumulative number of stats-tube commands"},
		{"cmd-use": "the cumulative number of use commands"},
		{"cmd-watch": "the cumulative number of watch commands"},
	}
	currentStatsGroups = []map[string]string{
		{"current-connections": "the number of currently open connections"},
		{"current-jobs-buried": "the number of buried jobs"},
		{"current-jobs-delayed": "the number of delayed jobs"},
		{"current-jobs-ready": "the number of jobs in the ready queue"},
		{"current-jobs-reserved": "the number of jobs reserved by all clients"},
		{"current-jobs-urgent": "the number of ready jobs with priority &lt; 1024"},
		{"current-producers": "the number of open connections that have each issued at least one put command"},
		{"current-tubes": "the number of currently-existing tubes"},
		{"current-waiting": "the number of open connections that have issued a reserve command but not yet received a response"},
		{"current-workers": "the number of open connections that have each issued at least one reserve command"},
	}
	otherStatsGroups = []map[string]string{
		{"hostname": "the hostname of the machine as determined by uname"},
		{"id": "a random id string for this server process}, generated when emap[string]string{ach beanstalkd process starts"},
		{"job-timeouts": "the cumulative count of times a job has timed out"},
		{"max-job-size": "the maximum number of bytes in a job"},
		{"pid": "the process id of the server"},
		{"rusage-stime": "the cumulative system CPU time of this process in seconds and microseconds"},
		{"rusage-utime": "the cumulative user CPU time of this process in seconds and microseconds"},
		{"total-connections": "the cumulative count of connections"},
		{"total-jobs": "the cumulative count of jobs created"},
		{"uptime": "the number of seconds since this server process started running"},
		{"version": "the version string of the server"},
	}
	// Tube filter columns.
	tubeStatFields = []map[string]string{
		{"current-jobs-urgent": "number of ready jobs with priority &lt; 1024 in this tube"},
		{"current-jobs-ready": "number of jobs in the ready queue in this tube"},
		{"current-jobs-reserved": "number of jobs reserved by all clients in this tube"},
		{"current-jobs-delayed": "number of delayed jobs in this tube"},
		{"current-jobs-buried": "number of buried jobs in this tube"},
		{"current-using": "number of open connections that are currently using this tube"},
		{"current-waiting": "number of open connections that have issued a reserve command while watching this tube but not yet received a response"},
		{"current-watching": "number of open connections that are currently watching this tube"},
		{"cmd-delete": "cumulative number of delete commands for this tube"},
		{"cmd-pause-tube": "cumulative number of pause-tube commands for this tube"},
		{"pause": "number of seconds the tube has been paused for"},
		{"pause-time-left": "number of seconds until the tube is un-paused"},
		{"total-jobs": "cumulative count of jobs created in this tube in the current beanstalkd process"},
	}
	statisticsFields = []map[string]string{
		{"ready": "current-jobs-ready"},
		{"delayed": "current-jobs-delayed"},
		{"reserved": "current-jobs-reserved"},
		{"buried": "current-jobs-buried"},
	}
	jobStatsOrder = []string{"id", "tube", "state", "pri", "age", "delay", "ttr", "time-left", "file", "reserves", "timeouts", "releases", "buries", "kicks"}
)

// ViewFunc define HTTP Basic Auth type of return function.
type ViewFunc func(http.ResponseWriter, *http.Request)

// PubConfig define struct for prase config file.
type PubConfig struct {
	Servers  []string `toml:"servers"`
	Listen   string   `toml:"listen"`
	Version  float64  `toml:"version"`
	OpenPage struct {
		Enabled bool `toml:"enabled"`
	} `toml:"openpage"`
	Auth struct {
		Enabled  bool   `toml:"enabled"`
		Password string `toml:"password"`
		Username string `toml:"username"`
	} `toml:"auth"`
	Sample struct {
		Storage string `toml:"storage"`
	} `toml:"sample"`
}

// SampleJobs define beanstalkd sample jobs storage struct.
type SampleJobs struct {
	Jobs  []SampleJob  `json:"jobs"`
	Tubes []SampleTube `json:"tubes"`
}

// SampleJob define beanstalkd sample job storage struct.
type SampleJob struct {
	Key   string   `json:"key"`
	Name  string   `json:"name"`
	Tubes []string `json:"tubes"`
	Data  string   `json:"data"`
	TTR   int      `json:"ttr"`
}

// SampleTube define beanstalkd sample job's tube storage struct.
type SampleTube struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// StatisticsData define the data struct for storage statistics data.
type StatisticsData struct {
	*sync.RWMutex
	Server map[string]map[string]map[string]*list.List
}

// SelfConf define fields storage in cookies and statistics parameter storage in
// RAM.
type SelfConf struct {
	Filter                     []string
	Servers                    []string
	TubeFilters                []string
	TubeSelector               string
	TubePauseSeconds           int
	IsDisabledJSONDecode       int
	IsDisabledUnserialization  int
	IsDisabledJobDataHighlight int
	IsEnabledBase64Decode      int
	AutoRefreshTimeoutMs       int
	SearchResultLimit          int
	StatisticsCollection       int
	StatisticsFrequency        int
}

// SearchResult define the search result of jobs in tube.
type SearchResult struct {
	ID    uint64
	State string
	Data  string
}

// UpdateTags define the tag of versions control.
type UpdateTags []struct {
	Name string `json:"name"`
}
