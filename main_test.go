package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rakyll/statik/fs"
)

const server = "http://127.0.0.1:3000"

var (
	once sync.Once
	urls = []string{
		"/",                                                                                                   // Static files server
		"/public",                                                                                             // Server list
		"/server?server=127.0.0.1:11300",                                                                      // Server status
		"/index?server=&action=reloader&tplMain=ajax&tplBlock=serversList",                                    // Reload server status
		"/serversRemove?action=serversRemove&removeServer=127.0.0.1:11300",                                    // Remove server
		"/server?server=127.0.0.1:11300&action=reloader&tplMain=ajax&tplBlock=allTubes",                       // Reload tube status
		"/tube?server=127.0.0.1:11300&tube=default",                                                           // Tube status
		"/tube?server=127.0.0.1:11300&tube=default&action=pause&count=-1",                                     // Pause tube
		"/tube?server=127.0.0.1:11300&tube=default&action=pause&count=0",                                      // Pause tube
		"/tube?server=127.0.0.1:11300&tube=default&action=kick&count=1",                                       // Kick 1 job
		"/tube?server=127.0.0.1:11300&tube=default&action=kick&count=10",                                      // Kick 10 job
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=deleteJob&jobid=1",                      // Delete a job
		"/tube?server=127.0.0.1:11300&tube=default&state=ready&action=deleteAll&count=1",                      // Delete all jobs
		"/tube?server=127.0.0.1:11300&tube=aurora_test&state=&action=search&limit=25&searchStr=t",             // Search job
		"/tube?server=127.0.0.1:11300&tube=aurora_test&action=moveJobsTo&destState=buried&state=ready",        // Move job from ready to buried state
		"/tube?server=127.0.0.1:11300&tube=aurora_test&action=moveJobsTo&destTube=aurora_test&state=buried",   // Move job from buried to ready state
		"/sample?action=manageSamples",                                                                        // Manage sample jobs
		"/tube?server=127.0.0.1:11300&tube=auto&action=loadSample&key=xxx&redirect=tube?action=manageSamples", // Kick job to tubes
		"/sample?action=newSample",                                                                            // New sample job
		"/sample?action=editSample&key=xxx",                                                                   // Edit sample job
	}
)

func testSetup() {
	parseFlags()
	if _, err := toml.Decode(ConfigFileTemplate, &pubConf); err != nil {
		return
	}
	if err := json.Unmarshal([]byte(pubConf.Sample.Storage), &sampleJobs); err != nil {
		return
	}
}

func TestIndex(t *testing.T) {
	once.Do(testSetup)

	var resp *http.Response
	var err error

	statikFS, err := fs.New()
	if err != nil {
		t.Log(err)
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

	time.Sleep(2 * time.Second)

	resp, err = http.PostForm(server+"/tube?server=127.0.0.1:11300&action=addjob",
		url.Values{"tubeName": {"aurora_test"}, "tubeData": {"test"}})
	if err != nil {
		t.Log(err)
	}

	resp, err = http.PostForm(server+"/tube?server=127.0.0.1:11300&action=addSample",
		url.Values{"tube": {"aurora_test"}, "addsamplejobid": {"1"}, "addsamplename": {"test_sample_1"}, "tubes[aurora_test]": {"1"}})
	if err != nil {
		t.Log(err)
	}
	defer resp.Body.Close()

	for _, v := range urls {
		req, err := http.NewRequest("GET", server+v, nil)
		if err != nil {
			t.Log(err)
		}
		cookie := http.Cookie{Name: "beansServers", Value: `127.0.0.1%3A11300`}
		req.AddCookie(&cookie)
		var client = &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			t.Log(err)
		}
	}
	return
}
