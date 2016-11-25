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
	ConfigFileTemplate      = "servers = []\r\nlisten = \"127.0.0.1:3000\"\r\nversion = 1.6\r\n\r\n[auth]\r\nenabled = false\r\npassword = \"password\"\r\nusername = \"admin\"\r\n\r\n[sample]\r\nstorage = \"{}\""
	DefaultDelay            = 0
	DefaultPriority         = 1024 // most urgent: 0, least urgent: 4294967295.
	DefaultTTR              = 60   // 1 minute
	DefaultTubePauseSeconds = 3600
	TplHead                 = `<head><meta charset="UTF-8"/><meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"><!--[if IE]><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><![endif]--><meta name="description" content="Beanstalk Console"><meta name="keywords" content="Beanstalk Console, beanstalk, console"><meta content="always" name="referrer"><meta name="language" content="en-US"><meta name="category" content="Tools"><meta name="summary" content="Beanstalk Console"><meta name="apple-mobile-web-app-capable" content="yes"/><link rel="copyright" href="http://www.opensource.org/licenses/mit-license.php"/><link rel="icon" sizes="32x32" href="./images/aurora-32x32.ico"><link rel="apple-touch-icon-precomposed" sizes="180x180" href="./images/apple-touch-icon-180x180-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="152x152" href="./images/apple-touch-icon-152x152-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="144x144" href="./images/apple-touch-icon-144x144-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="120x120" href="./images/apple-touch-icon-120x120-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="114x114" href="./images/apple-touch-icon-114x114-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="76x76" href="./images/apple-touch-icon-76x76-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="72x72" href="./images/apple-touch-icon-72x72-precomposed.png"><link rel="apple-touch-icon-precomposed" href="./images/apple-touch-icon-precomposed-57x57.png"><title>Beanstalk Console</title><!-- Bootstrap core CSS --><link href="./assets/vendor/bootstrap/css/bootstrap.min.css" rel="stylesheet"><link href="./css/customer.css" rel="stylesheet"><link href="./highlight/styles/magula.css" rel="stylesheet"><!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries --><!--[if lt IE 9]><script src="./js/libs/html5shiv/3.7.0/html5shiv.js"></script><script src="./js/libs/respond.js/1.4.2/respond.min.js"></script><![endif]--></head>`
	TplLinks                = `<li class="dropdown"><a href="#" class="dropdown-toggle" data-toggle="dropdown"> Links <span class="caret"></span></a><ul class="dropdown-menu"><li><a href="https://github.com/kr/beanstalkd">Beanstalk (GitHub)</a></li><li><a href="https://github.com/Luxurioust/aurora">Aurora (GitHub)</a></li></ul></li>`
	TplNoScript             = `<noscript><div class="container"><div class="alert alert-danger" role="alert">Aurora beanstalk console requires JavaScript supports, please refresh after enable browser JavaScript support.</div></div></noscript>`
	UpdateURL               = `https://api.github.com/repos/Luxurioust/aurora/tags`
	Version                 = 1.6
)

// Define server and tube stats fields.
var (
	pubConf              PubConfig
	sampleJobs           SampleJobs
	selfConf             SelfConf
	Stderr               io.Writer = os.Stderr // Stderr is the io.Writer to which executed commands write standard error.
	Stdout               io.Writer = os.Stdout // Stdout is the io.Writer to which executed commands write standard output.
	ConfigFile                     = `.` + string(os.PathSeparator) + `aurora.toml`
	statisticsData                 = StatisticsData{new(sync.RWMutex), statisticsDataServer}
	statisticsDataServer           = make(map[string]map[string]map[string]*list.List)
	notify                         = make(chan bool, 1)
	updateInfo                     = "uncheck"
	// Server filter columns.
	binlogStatsGroups = []map[string]string{
		map[string]string{"binlog-current-index": "the index of the current binlog file being written to. If binlog is not active this value will be 0"},
		map[string]string{"binlog-max-size": "the maximum size in bytes a binlog file is allowed to get before a new binlog file is opened"},
		map[string]string{"binlog-oldest-index": "the index of the oldest binlog file needed to store the current jobs"},
		map[string]string{"binlog-records-migrated": "the cumulative number of records written as part of compaction"},
		map[string]string{"binlog-records-written": "the cumulative number of records written to the binlog"},
	}
	cmdStatsGroups = []map[string]string{
		map[string]string{"cmd-bury": "the cumulative number of bury commands"},
		map[string]string{"cmd-delete": "the cumulative number of delete commands"},
		map[string]string{"cmd-ignore": "the cumulative number of ignore commands"},
		map[string]string{"cmd-kick": "the cumulative number of kick commands"},
		map[string]string{"cmd-list-tube-used": "the cumulative number of list-tube-used commands"},
		map[string]string{"cmd-list-tubes": "the cumulative number of list-tubes commands"},
		map[string]string{"cmd-list-tubes-watched": "the cumulative number of list-tubes-watched commands"},
		map[string]string{"cmd-pause-tube": "the cumulative number of pause-tube commands"},
		map[string]string{"cmd-peek": "the cumulative number of peek commands"},
		map[string]string{"cmd-peek-buried": "the cumulative number of peek-buried commands"},
		map[string]string{"cmd-peek-delayed": "the cumulative number of peek-delayed commands"},
		map[string]string{"cmd-peek-ready": "the cumulative number of peek-ready commands"},
		map[string]string{"cmd-put": "the cumulative number of put commands"},
		map[string]string{"cmd-release": "the cumulative number of release commands"},
		map[string]string{"cmd-reserve": "the cumulative number of reserve commands"},
		map[string]string{"cmd-stats": "the cumulative number of stats commands"},
		map[string]string{"cmd-stats-job": "the cumulative number of stats-job commands"},
		map[string]string{"cmd-stats-tube": "the cumulative number of stats-tube commands"},
		map[string]string{"cmd-use": "the cumulative number of use commands"},
		map[string]string{"cmd-watch": "the cumulative number of watch commands"},
	}
	currentStatsGroups = []map[string]string{
		map[string]string{"current-connections": "the number of currently open connections"},
		map[string]string{"current-jobs-buried": "the number of buried jobs"},
		map[string]string{"current-jobs-delayed": "the number of delayed jobs"},
		map[string]string{"current-jobs-ready": "the number of jobs in the ready queue"},
		map[string]string{"current-jobs-reserved": "the number of jobs reserved by all clients"},
		map[string]string{"current-jobs-urgent": "the number of ready jobs with priority &lt; 1024"},
		map[string]string{"current-producers": "the number of open connections that have each issued at least one put command"},
		map[string]string{"current-tubes": "the number of currently-existing tubes"},
		map[string]string{"current-waiting": "the number of open connections that have issued a reserve command but not yet received a response"},
		map[string]string{"current-workers": "the number of open connections that have each issued at least one reserve command"},
	}
	otherStatsGroups = []map[string]string{
		map[string]string{"hostname": "the hostname of the machine as determined by uname"},
		map[string]string{"id": "a random id string for this server process}, generated when emap[string]string{ach beanstalkd process starts"},
		map[string]string{"job-timeouts": "the cumulative count of times a job has timed out"},
		map[string]string{"max-job-size": "the maximum number of bytes in a job"},
		map[string]string{"pid": "the process id of the server"},
		map[string]string{"rusage-stime": "the cumulative system CPU time of this process in seconds and microseconds"},
		map[string]string{"rusage-utime": "the cumulative user CPU time of this process in seconds and microseconds"},
		map[string]string{"total-connections": "the cumulative count of connections"},
		map[string]string{"total-jobs": "the cumulative count of jobs created"},
		map[string]string{"uptime": "the number of seconds since this server process started running"},
		map[string]string{"version": "the version string of the server"},
	}
	// Tube filter columns.
	tubeStatFields = []map[string]string{
		map[string]string{"current-jobs-urgent": "number of ready jobs with priority &lt; 1024 in this tube"},
		map[string]string{"current-jobs-ready": "number of jobs in the ready queue in this tube"},
		map[string]string{"current-jobs-reserved": "number of jobs reserved by all clients in this tube"},
		map[string]string{"current-jobs-delayed": "number of delayed jobs in this tube"},
		map[string]string{"current-jobs-buried": "number of buried jobs in this tube"},
		map[string]string{"current-using": "number of open connections that are currently using this tube"},
		map[string]string{"current-waiting": "number of open connections that have issued a reserve command while watching this tube but not yet received a response"},
		map[string]string{"current-watching": "number of open connections that are currently watching this tube"},
		map[string]string{"cmd-delete": "cumulative number of delete commands for this tube"},
		map[string]string{"cmd-pause-tube": "cumulative number of pause-tube commands for this tube"},
		map[string]string{"pause": "number of seconds the tube has been paused for"},
		map[string]string{"pause-time-left": "number of seconds until the tube is un-paused"},
		map[string]string{"total-jobs": "cumulative count of jobs created in this tube in the current beanstalkd process"},
	}
	statisticsFields = []map[string]string{
		map[string]string{"ready": "current-jobs-ready"},
		map[string]string{"delayed": "current-jobs-delayed"},
		map[string]string{"reserved": "current-jobs-reserved"},
		map[string]string{"buried": "current-jobs-buried"},
	}
	jobStatsOrder = []string{"id", "tube", "state", "pri", "age", "delay", "ttr", "time-left", "file", "reserves", "timeouts", "releases", "buries", "kicks"}
)

// ViewFunc define HTTP Basic Auth type of return function.
type ViewFunc func(http.ResponseWriter, *http.Request)

// PubConfig define struct for prase config file.
type PubConfig struct {
	Servers []string `toml:"servers"`
	Listen  string   `toml:"listen"`
	Version float64  `toml:"version"`
	Auth    struct {
		Enabled  bool   `toml:"enabled"`
		Password string `toml:"password"`
		Username string `toml:"username"`
	} `toml:"auth"`
	Sample struct {
		Storage string `toml:"storage"`
	} `toml:"sample"`
}

// SampleJobs define beanstalk sample jobs storage struct.
type SampleJobs struct {
	Jobs  []SampleJob  `json:"jobs"`
	Tubes []SampleTube `json:"tubes"`
}

// SampleJob define beanstalk sample job storage struct.
type SampleJob struct {
	Key   string   `json:"key"`
	Name  string   `json:"name"`
	Tubes []string `json:"tubes"`
	Data  string   `json:"data"`
}

// SampleTube define beanstalk sample job's tube storage struct.
type SampleTube struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// StatisticsData define the data struct for storage statistics data.
type StatisticsData struct {
	*sync.RWMutex
	Server map[string]map[string]map[string]*list.List
}

// SelfConf define fields storage in cookies and statistics parameter storage in RAM.
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
