package main

import (
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	server                   = "http://127.0.0.1:3000"
	bstk                     = "127.0.0.1:11300"
	configFileWithSampleJobs = `servers = []
listen = "127.0.0.1:3000"
version = 2.0

[openpage]
enabled = true

[auth]
  enabled = false
  password = "password"
  username = "admin"

[sample]
  storage = "{\"jobs\":[{\"key\":\"97ec882fd75855dfa1b4bd00d4a367d4\",\"name\":\"sample_1\",\"tubes\":[\"default\"],\"data\":\"aurora_test_sample_job\",\"ttr\":60},{\"key\":\"d44782912092260cec11275b73f78434\",\"name\":\"sample_2\",\"tubes\":[\"aurora_test\"],\"data\":\"aurora_test_sample_job\",\"ttr\":60},{\"key\":\"5def7a2daf5e0292bed42db9f0017c94\",\"name\":\"sample_3\",\"tubes\":[\"default\",\"aurora_test\"],\"data\":\"aurora_test_sample_job\",\"ttr\":60}],\"tubes\":[{\"name\":\"default\",\"keys\":[\"97ec882fd75855dfa1b4bd00d4a367d4\",\"5def7a2daf5e0292bed42db9f0017c94\"]},{\"name\":\"aurora_test\",\"keys\":[\"d44782912092260cec11275b73f78434\",\"5def7a2daf5e0292bed42db9f0017c94\"]}]}"`
)

var (
	once sync.Once
	urls = []string{
		"/",                      // Static files server
		"/public",                // Server list
		"/server?server=" + bstk, // Server status
		"/index?server=&action=reloader&tplMain=ajax&tplBlock=serversList",                                                                 // Reload server status
		"/serversRemove?action=serversRemove&removeServer=" + bstk,                                                                         // Remove server
		"/server?server=" + bstk + "&action=reloader&tplMain=ajax&tplBlock=allTubes",                                                       // Reload tube status
		"/tube?server=" + bstk + "&tube=default",                                                                                           // Tube status
		"/tube?server=" + bstk + "&tube=default1",                                                                                          // Tube status with no exits tube
		"/tube?server=not_exist_server_addr&tube=default",                                                                                  // Tube status with no exits server
		"/tube?server=" + bstk + "&tube=default&action=pause&count=-1",                                                                     // Pause tube
		"/tube?server=" + bstk + "&tube=default&action=pause&count=0",                                                                      // Pause tube
		"/tube?server=not_exist_server_addr&tube=default&action=pause&count=0",                                                             // Pause tube with no exits server
		"/tube?server=" + bstk + "&tube=default&action=kick&count=1",                                                                       // Kick 1 job
		"/tube?server=" + bstk + "&tube=default&action=kick&count=10",                                                                      // Kick 10 job
		"/tube?server=" + bstk + "&tube=default&state=ready&action=kickJob&jobid=1",                                                        // Kick job by given ID
		"/tube?server=" + bstk + "&tube=default&state=ready&action=kickJob&jobid=badID",                                                    // Kick job by given ID with no exits ID
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=kickJob&jobid=1",                                               // Kick job by given ID with no exits server
		"/tube?server=" + bstk + "&tube=default&action=loadSample&key=97ec882fd75855dfa1b4bd00d4a367d4&redirect=tube&action=manageSamples", // Load sample job by given key
		"/tube?server=" + bstk + "&tube=aurora_test&state=&action=search&limit=25&searchStr=t",                                             // Search job
		"/tube?server=" + bstk + "&tube=aurora_test&state=&action=search&limit=25&searchStr=match",                                         // Search job with not match string
		"/tube?server=" + bstk + "&tube=aurora_test&action=moveJobsTo&destState=buried&state=ready",                                        // Move job from ready to buried state
		"/tube?server=not_exist_server_addr&tube=aurora_test&action=moveJobsTo&destState=buried&state=ready",                               // Move job from ready to buried state with no exits server
		"/tube?server=" + bstk + "&tube=aurora_test&action=moveJobsTo&destState=&state=ready",                                              // Move job from ready to buried state without destState
		"/tube?server=" + bstk + "&tube=aurora_test&action=moveJobsTo&destTube=aurora_test&state=buried",                                   // Move job from buried to ready state
		"/tube?server=not_exist_server_addr&tube=aurora_test&action=moveJobsTo&destTube=aurora_test&state=buried",                          // Move job from buried to ready state with no exits server
		"/sample?action=manageSamples",                                                                                                     // Manage sample jobs
		"/tube?server=" + bstk + "&tube=auto&action=loadSample&key=xxx&redirect=tube?action=manageSamples",                                 // Kick job to tubes
		"/sample?action=newSample",          // New sample job
		"/sample?action=editSample&key=xxx", // Edit sample job
		"/tube?server=" + bstk + "&tube=default&state=ready&action=deleteJob&jobid=1",          // Delete a job
		"/tube?server=" + bstk + "&tube=default&state=ready&action=deleteJob&jobid=badID",      // Delete a no exists job
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=deleteJob&jobid=1", // Delete a job with no exits server
		"/tube?server=" + bstk + "&tube=default&state=ready&action=deleteAll&count=1",          // Delete all jobs in empty tube
		"/tube?server=" + bstk + "&tube=aurora_test&state=ready&action=deleteAll&count=1",      // Delete all jobs
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=deleteAll&count=1", // Delete all jobs with no exits server
	}
)

func testSetup() {
	time.Sleep(1 * time.Second) // Wait Beanstalkd server ready.
	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	readConf()
	saveSample()
	selfConf := selfDir + string(os.PathSeparator) + `aurora.toml`
	createFile(selfConf)
	writeFile(selfConf, ConfigFileTemplate)
	parseFlags()
	readConf()
	public, _ := fs.Sub(staticFiles, "public")
	// handle static files include HTML, CSS and JavaScripts.
	http.Handle("/", http.FileServer(http.FS(public)))
	http.HandleFunc("/public", basicAuth(handlerMain))
	http.HandleFunc("/index", basicAuth(handlerServerList))
	http.HandleFunc("/serversRemove", basicAuth(serversRemove))
	http.HandleFunc("/server", basicAuth(handlerServer))
	http.HandleFunc("/tube", basicAuth(handlerTube))
	http.HandleFunc("/sample", basicAuth(handlerSample))
	http.HandleFunc("/statistics", basicAuth(handlerStatistics))
	go func() {
		http.ListenAndServe(pubConf.Listen, nil)
	}()
	go statistic()
}

func TestIndex(t *testing.T) {
	once.Do(testSetup)
	var resp *http.Response
	var err error
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server=not_exist_server_addr&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {"test_sample_1"}, "tubes[aurora_test]": {"1"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {"test_sample_1"}, "tubes[aurora_test]": {"1"}, "addsamplettr": {"60"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {""}, "tubes[aurora_test]": {"1"}, "addsamplettr": {"60"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"default"}, "addsamplejobid": {"1"}, "addsamplename": {"sample_1"}, "tubes[aurora_test]": {"1"}, "addsamplettr": {"60"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server=not_exist_server_addr&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {"sample_2"}, "tubes[default]": {"1"}, "addsamplettr": {"60"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {""}, "addsamplename": {"sample_2"}, "tubes[aurora_test]": {"1"}, "addsamplettr": {"60"}})
	if err != nil {
		t.Log(err)
	}
	_, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"not_int"}, "addsamplename": {"sample_2"}, "tubes[aurora_test]": {"1"}, "addsamplettr": {"60"}})
	if err != nil {
		t.Log(err)
	}
	for _, v := range urls {
		req, err := http.NewRequest("GET", server+v, nil)
		if err != nil {
			t.Log(err)
		}
		cookie := http.Cookie{Name: "beansServers", Value: `127.0.0.1%3A11300%3B127.0.0.1%3A11300%3B127.0.0.1%3A11301%3B`}
		req.AddCookie(&cookie)
		var client = &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			t.Log(err)
		}
	}
	resp, err = http.PostForm(server+"/server?server="+bstk+"&action=clearTubes",
		url.Values{"default": {"1"}})
	if err != nil {
		t.Log(err)
	}
	defer resp.Body.Close()
}

func TestCookie(t *testing.T) {
	once.Do(testSetup)
	var err error
	var req *http.Request
	var cookie http.Cookie
	var client = &http.Client{}
	req, err = http.NewRequest("GET", server+urls[1], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "filter", Value: `binlog-current-index,binlog-max-size`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	req, err = http.NewRequest("GET", server+urls[1], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "tubefilter", Value: `current-using,current-waiting`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	req, err = http.NewRequest("GET", server+urls[1], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "tubeSelector", Value: `default`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	// Test readIntCookie
	req, err = http.NewRequest("GET", server+urls[1], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "isDisabledJsonDecode", Value: `1`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	// Test removeServerInCookie
	req, err = http.NewRequest("GET", server+urls[4], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "beansServers", Value: `not_exist_server_addr`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	req, err = http.NewRequest("GET", server+urls[14], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "filter", Value: `binlog-current-index,binlog-max-size`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
}

func TestCurrentTubeJobsActionsRow(t *testing.T) {
	once.Do(testSetup)
	var err error
	var req *http.Request
	var cookie http.Cookie
	var client = &http.Client{}
	selfConf := `.` + string(os.PathSeparator) + `aurora.toml`
	os.Remove(selfConf)
	createFile(selfConf)
	writeFile(selfConf, configFileWithSampleJobs)
	readConf()
	req, err = http.NewRequest("GET", server+urls[6], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "filter", Value: `binlog-current-index,binlog-max-size`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	req, err = http.NewRequest("GET", server+urls[6], nil)
	if err != nil {
		t.Log(err)
	}
	cookie = http.Cookie{Name: "filter", Value: `binlog-current-index,binlog-max-size`}
	req.AddCookie(&cookie)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
}

func TestMoveReadyJobsTo(t *testing.T) {
	once.Do(testSetup)
	moveReadyJobsTo(bstk, `default`, `aurora_test_1`, `ready`)
}

func TestSearchTube(t *testing.T) {
	once.Do(testSetup)
	searchTube(bstk, `default`, `not_int`, `ready`)
	searchTube(bstk, `aurora_test_2`, `1`, `ready`)
}

func TestAddSample(t *testing.T) {
	once.Do(testSetup)
	var resp *http.Response
	var err error
	resp, err = http.PostForm(server+"/tube?server="+bstk+"&action=addjob",
		url.Values{"tubeName": {"default"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}
	resp.Body.Close()
	time.Sleep(time.Second)
	form := url.Values{}
	form.Add(`name`, `sample_job_1`)
	form.Add(`jobdata`, `sample_job_body`)
	form.Add(`action`, `actionNewSample`)
	form.Add(`tubes[default]`, `1`)
	form.Add(`addsamplejobid`, `6`)
	form.Add(`ttr`, `60`)
	req, err := http.NewRequest("POST", server+`/sample?action=actionNewSample`, strings.NewReader(form.Encode()))
	if err != nil {
		t.Log(err)
	}
	var client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
}

func TestEditSample(t *testing.T) {
	once.Do(testSetup)
	req, err := http.NewRequest("GET", server+`/sample?action=editSample&key=d44782912092260cec11275b73f78434`, nil)
	if err != nil {
		t.Log(err)
	}
	cookie := http.Cookie{Name: "beansServers", Value: `127.0.0.1%3A11300%3B127.0.0.1%3A11300%3B127.0.0.1%3A11301%3B`}
	req.AddCookie(&cookie)
	var client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
}

func TestGetServerTubes(t *testing.T) {
	once.Do(testSetup)
	getServerTubes("")
}

func TestPrettyJSON(t *testing.T) {
	once.Do(testSetup)
	prettyJSON([]byte(`{}`))
}

func TestBase64Decode(t *testing.T) {
	once.Do(testSetup)
	base64Decode(`dGVzdA==`)
	base64Decode(`test-%?s`)
}

func TestDropEditSettings(t *testing.T) {
	once.Do(testSetup)
	selfConf.IsEnabledBase64Decode = 1
	dropEditSettings()
}

func TestRemoveServerInConfig(t *testing.T) {
	once.Do(testSetup)
	removeServerInConfig(bstk)
}

func TestCheckUpdate(t *testing.T) {
	once.Do(testSetup)
	checkUpdate()
}

func TestSaveSample(t *testing.T) {
	once.Do(testSetup)
	saveSample()
}

func TestAddSampleTube(t *testing.T) {
	once.Do(testSetup)
	addSampleTube(`aurora_test_2`, `test`)
	getSampleJobList()
	getSampleJobNameByKey(`97ec882fd75855dfa1b4bd00d4a367d4`)
	loadSample(``, `default`, `97ec882fd75855dfa1b4bd00d4a367d4`)
	loadSample(bstk, `default`, `97ec882fd75855dfa1b4bd00d4a367d4`)
	deleteSamples(`97ec882fd75855dfa1b4bd00d4a367d4`)
}

func TestBasicAuth(t *testing.T) {
	once.Do(testSetup)
	var err error
	var req *http.Request
	var client = &http.Client{}
	pubConf.Auth.Enabled = true
	http.HandleFunc("/test", basicAuth(handlerMain))
	req, err = http.NewRequest("GET", server+"/test", nil)
	if err != nil {
		t.Log(err)
	}
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
	req, err = http.NewRequest("GET", server+"/test", nil)
	if err != nil {
		t.Log(err)
	}
	req.SetBasicAuth(`admin`, `password`)
	client = &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Log(err)
	}
}

func TestDeleteSamples(t *testing.T) {
	once.Do(testSetup)
	deleteSamples(``)
	deleteSamples(`test`)
}

func TestClearTubes(t *testing.T) {
	once.Do(testSetup)
	modalClearTubes("")
}

func TestRunCmd(t *testing.T) {
	once.Do(testSetup)
	err := runCmd(`date`)
	if err != nil {
		t.Log(err)
	}
}

func TestStatistic(t *testing.T) {
	once.Do(testSetup)
	var resp *http.Response
	var err error
	var testURLs = []string{"/statistics?action=preference", "/statistics", "/statistics?action=reloader"}
	_, err = http.PostForm(server+"/statistics?action=save",
		url.Values{"frequency": {"1"}, "collection": {"10"}, "tubes[" + bstk + ":default]": {"1"}})
	if err != nil {
		t.Log(err)
	}
	time.Sleep(10 * time.Second)
	tplStatistic(bstk, "default")
	statisticWaitress(bstk, "default")
	tplStatisticEdit("")
	tplStatisticSetting("")
	for _, v := range testURLs {
		req, err := http.NewRequest("GET", server+v, nil)
		if err != nil {
			t.Log(err)
		}
		cookie := http.Cookie{Name: "beansServers", Value: `127.0.0.1%3A11300%3B127.0.0.1%3A11300%3B127.0.0.1%3A11301%3B`}
		req.AddCookie(&cookie)
		var client = &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			t.Log(err)
		}
	}
	currentTubeJobsSummaryTable(bstk, "default")
	resp, err = http.PostForm(server+"/statistics?action=save",
		url.Values{"frequency": {"-1"}, "collection": {"-1"}, "tubes[127.0.0.1:default]": {"1"}})
	if err != nil {
		t.Log(err)
	}
	time.Sleep(10 * time.Second)
	for _, v := range testURLs {
		req, err := http.NewRequest("GET", server+v, nil)
		if err != nil {
			t.Log(err)
		}
		cookie := http.Cookie{Name: "beansServers", Value: `127.0.0.1%3A11300%3B127.0.0.1%3A11300%3B127.0.0.1%3A11301%3B`}
		req.AddCookie(&cookie)
		var client = &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			t.Log(err)
		}
	}
	defer resp.Body.Close()
	statisticCashier("not_int", "", []string{})
	statisticCashier("1", "not_int", []string{})
	selfConf.StatisticsFrequency = -1
	tplStatisticEdit("")
	t.SkipNow()
}

func TestReadConf(t *testing.T) {
	once.Do(testSetup)
	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	selfConf := selfDir + string(os.PathSeparator) + `aurora.toml`
	os.Remove(selfConf)
	readConf()
}

func TestMain(t *testing.T) {
	once.Do(testSetup)
	go func() {
		openPage()
		handleSignals()
	}()
	time.Sleep(5 * time.Second)
	pubConf.OpenPage.Enabled = false
	go func() {
		openPage()
	}()
	time.Sleep(5 * time.Second)
	t.SkipNow()
}

func createFile(path string) {
	// Detect if file exists
	var _, err = os.Stat(path)
	// Create file if not exists
	if os.IsNotExist(err) {
		var file, _ = os.Create(path)
		defer file.Close()
	}
}

func writeFile(path string, content string) {
	// Open file using READ & WRITE permission
	var file, _ = os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()
	file.WriteString(content) // Write some text to file
	file.Sync()               // Save changes
}
