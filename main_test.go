package main

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/rakyll/statik/fs"
)

const (
	server                   = "http://127.0.0.1:3000"
	bstk                     = "127.0.0.1:11300"
	configFileWithSampleJobs = `servers = []
listen = "127.0.0.1:3000"
version = 1.3
[auth]
  enabled = false
  password = "password"
  username = "admin"

[sample]
  storage = "{\"jobs\":[{\"key\":\"97ec882fd75855dfa1b4bd00d4a367d4\",\"name\":\"sample_1\",\"tubes\":[\"default\"],\"data\":\"aurora_test_sample_job\"},{\"key\":\"d44782912092260cec11275b73f78434\",\"name\":\"sample_2\",\"tubes\":[\"aurora_test\"],\"data\":\"aurora_test_sample_job\"},{\"key\":\"5def7a2daf5e0292bed42db9f0017c94\",\"name\":\"sample_3\",\"tubes\":[\"default\",\"aurora_test\"],\"data\":\"aurora_test_sample_job\"}],\"tubes\":[{\"name\":\"default\",\"keys\":[\"97ec882fd75855dfa1b4bd00d4a367d4\",\"5def7a2daf5e0292bed42db9f0017c94\"]},{\"name\":\"aurora_test\",\"keys\":[\"d44782912092260cec11275b73f78434\",\"5def7a2daf5e0292bed42db9f0017c94\"]}]}"
`
)

var (
	once sync.Once
	urls = []string{
		"/",                                                                                                                                // Static files server
		"/public",                                                                                                                          // Server list
		"/server?server=" + bstk,                                                                                                           // Server status
		"/index?server=&action=reloader&tplMain=ajax&tplBlock=serversList",                                                                 // Reload server status
		"/serversRemove?action=serversRemove&removeServer=" + "127.0.0.1:11300",                                                            // Remove server
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
		"/sample?action=newSample",                                                                                                         // New sample job
		"/sample?action=editSample&key=xxx",                                                                                                // Edit sample job
		"/tube?server=" + bstk + "&tube=default&state=ready&action=deleteJob&jobid=1",                                                      // Delete a job
		"/tube?server=" + bstk + "&tube=default&state=ready&action=deleteJob&jobid=badID",                                                  // Delete a no exists job
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=deleteJob&jobid=1",                                             // Delete a job with no exits server
		"/tube?server=" + bstk + "&tube=default&state=ready&action=deleteAll&count=1",                                                      // Delete all jobs in empty tube
		"/tube?server=" + bstk + "&tube=aurora_test&state=ready&action=deleteAll&count=1",                                                  // Delete all jobs
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=deleteAll&count=1",                                             // Delete all jobs with no exits server
	}
)

func testSetup() {
	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	selfConf := selfDir + string(os.PathSeparator) + `aurora.toml`
	createFile(selfConf)
	writeFile(selfConf, ConfigFileTemplate)
	parseFlags()
	readConf()
	statikFS, err := fs.New()
	if err != nil {
		return
	}
	http.FileServer(statikFS)
	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	http.HandleFunc("/public", basicAuth(handlerMain))
	http.HandleFunc("/index", basicAuth(handlerServerList))
	http.HandleFunc("/serversRemove", basicAuth(serversRemove))
	http.HandleFunc("/server", basicAuth(handlerServer))
	http.HandleFunc("/tube", basicAuth(handlerTube))
	http.HandleFunc("/sample", basicAuth(handlerSample))
	go func() {
		http.ListenAndServe(pubConf.Listen, nil)
	}()
}

func TestIndex(t *testing.T) {
	once.Do(testSetup)

	var resp *http.Response
	var err error

	resp, err = http.PostForm(server+"/tube?server="+bstk+"&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}

	resp, err = http.PostForm(server+"/tube?server=not_exist_server_addr&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}

	resp, err = http.PostForm(server+"/tube?server="+bstk+"&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {"test_sample_1"}, "tubes[aurora_test]": {"1"}})
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

	return
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

	return
}

func TestCurrentTubeJobsActionsRow(t *testing.T) {
	once.Do(testSetup)

	var err error
	var req *http.Request
	var cookie http.Cookie
	var client = &http.Client{}

	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	selfConf := selfDir + string(os.PathSeparator) + `aurora.toml`
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

	return
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
