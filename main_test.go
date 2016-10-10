package main

import (
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/rakyll/statik/fs"
)

const server = "http://127.0.0.1:3000"

var (
	once sync.Once
	urls = []string{
		"/",                                                                                                       // Static files server
		"/public",                                                                                                 // Server list
		"/server?server=127.0.0.1:11300",                                                                          // Server status
		"/server?server=not_exist_server_addr",                                                                    // Server status with no exits server
		"/index?server=&action=reloader&tplMain=ajax&tplBlock=serversList",                                        // Reload server status
		"/serversRemove?action=serversRemove&removeServer=127.0.0.1:11300",                                        // Remove server
		"/server?server=127.0.0.1:11300&action=reloader&tplMain=ajax&tplBlock=allTubes",                           // Reload tube status
		"/tube?server=not_exist_server_addr&tube=default",                                                         // Tube status with no exits server
		"/tube?server=127.0.0.1:11300&tube=default&action=pause&count=-1",                                         // Pause tube
		"/tube?server=127.0.0.1:11300&tube=default&action=pause&count=0",                                          // Pause tube
		"/tube?server=not_exist_server_addr&tube=default&action=pause&count=0",                                    // Pause tube with no exits server
		"/tube?server=127.0.0.1:11300&tube=default&action=kick&count=1",                                           // Kick 1 job
		"/tube?server=127.0.0.1:11300&tube=default&action=kick&count=10",                                          // Kick 10 job
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=kickJob&jobid=1",                            // Kick job by given ID
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=kickJob&jobid=badID",                        // Kick job by given ID with no exits ID
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=kickJob&jobid=1",                      // Kick job by given ID with no exits server
		"/tube?server=127.0.0.1:11300&tube=aurora_test&state=&action=search&limit=25&searchStr=t",                 // Search job
		"/tube?server=127.0.0.1:11300&tube=aurora_test&state=&action=search&limit=25&searchStr=match",             // Search job with not match string
		"/tube?server=127.0.0.1:11300&tube=aurora_test&action=moveJobsTo&destState=buried&state=ready",            // Move job from ready to buried state
		"/tube?server=not_exist_server_addr&tube=aurora_test&action=moveJobsTo&destState=buried&state=ready",      // Move job from ready to buried state with no exits server
		"/tube?server=127.0.0.1:11300&tube=aurora_test&action=moveJobsTo&destState=&state=ready",                  // Move job from ready to buried state without destState
		"/tube?server=127.0.0.1:11300&tube=aurora_test&action=moveJobsTo&destTube=aurora_test&state=buried",       // Move job from buried to ready state
		"/tube?server=not_exist_server_addr&tube=aurora_test&action=moveJobsTo&destTube=aurora_test&state=buried", // Move job from buried to ready state with no exits server
		"/sample?action=manageSamples",                                                                            // Manage sample jobs
		"/tube?server=127.0.0.1:11300&tube=auto&action=loadSample&key=xxx&redirect=tube?action=manageSamples",     // Kick job to tubes
		"/sample?action=newSample",                                                                                // New sample job
		"/sample?action=editSample&key=xxx",                                                                       // Edit sample job
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=deleteJob&jobid=1",                          // Delete a job
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=deleteJob&jobid=badID",                      // Delete a no exists job
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=deleteJob&jobid=1",                    // Delete a job with no exits server
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=deleteAll&count=1",                          // Delete all jobs in empty tube
		"/tube?server=127.0.0.1:11300&tube=aurora_test&state=ready&action=deleteAll&count=1",                      // Delete all jobs
		"/tube?server=not_exist_server_addr&tube=default&state=ready&action=deleteAll&count=1",                    // Delete all jobs with no exits server
	}
)

// Prepare for testing
func testSetup() {
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

func TestCookie(t *testing.T) {
	once.Do(testSetup)
	cookies := []http.Cookie{
		http.Cookie{Name: "filter", Value: `binlog-current-index,binlog-max-size`},
		http.Cookie{Name: "tubefilter", Value: `current-using,current-waiting`},
		http.Cookie{Name: "tubeSelector", Value: `default`},
		http.Cookie{Name: "isDisabledJsonDecode", Value: `1`},
	}
	for _, v := range cookies {
		req, err := http.NewRequest("GET", server+urls[1], nil)
		if err != nil {
			t.Log(err)
		}
		cookie := v
		req.AddCookie(&cookie)
		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			t.Log(err)
		}
	}
}

func TestForm(t *testing.T) {
	once.Do(testSetup)
	var resp *http.Response
	var err error
	resp, err = http.PostForm(server+"/tube?server=127.0.0.1:11300&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}
	resp, err = http.PostForm(server+"/tube?server=not_exist_server_addr&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}
	resp, err = http.PostForm(server+"/tube?server=127.0.0.1:11300&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {"test_sample_1"}, "tubes[aurora_test]": {"1"}})
	if err != nil {
		t.Log(err)
	}

	resp, err = http.PostForm(server+"/server?server=127.0.0.1:11300&action=clearTubes",
		url.Values{"default": {"1"}})
	if err != nil {
		t.Log(err)
	}
	defer resp.Body.Close()
}

func TestURL(t *testing.T) {
	once.Do(testSetup)
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
}
