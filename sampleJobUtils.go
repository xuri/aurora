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
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/xuri/aurora/beanstalk"
)

// addSample provide a function to add sample job by parse form with POST
// method.
func addSample(server string, data url.Values, w http.ResponseWriter) {
	var err error
	var key = randToken()
	var sampleName, body string
	var sampleTTR int
	var tubes []string

	err = readConf()
	if err != nil {
		return
	}

	sampleName, sampleTTR, body, err = sampleValidate(server, data, w)
	if err != nil {
		return
	}

	for k := range data { // range over map
		switch k {
		case "action", "tube", "addsamplejobid", "addsamplename", "addsamplettr", "server":
			continue
		default:
			t := strings.TrimSuffix(strings.TrimPrefix(k, `tubes[`), `]`)
			tubes = append(tubes, t)
			addSampleTube(t, key)
		}
	}
	sampleJobs.Jobs = append(sampleJobs.Jobs, SampleJob{
		Key:   key,
		Name:  sampleName,
		Tubes: tubes,
		Data:  body,
		TTR:   sampleTTR,
	})

	err = saveSample()
	if err != nil {
		return
	}
	_, _ = io.WriteString(w, `{"result":true}`)
}

// sampleValidate validate sample job if exists.
func sampleValidate(server string, data url.Values, w http.ResponseWriter) (string, int, string, error) {
	var bstkConn *beanstalk.Conn
	var sampleName string
	var sampleTTR = DefaultTTR
	var body []byte
	var err error
	sampleName = data.Get("addsamplename")
	if sampleName == "" {
		_, _ = io.WriteString(w, `{"result":false,"error":"You should give a name with this sample"}`)
		return sampleName, sampleTTR, string(body), errors.New("You should give a name with this sample")
	}
	if checkSampleJobs(sampleName) {
		_, _ = io.WriteString(w, `{"result":false,"error":"You already have a job with this name"}`)
		return sampleName, sampleTTR, string(body), errors.New("You already have a job with this name")
	}
	ID := data.Get("addsamplejobid")
	if ID == "" {
		_, _ = io.WriteString(w, `{"result":false,"error":"Job ID for add sample is empty"}`)
		return sampleName, sampleTTR, string(body), errors.New("Job ID for add sample is empty")
	}
	jobID, err := strconv.Atoi(ID)
	if err != nil {
		_, _ = io.WriteString(w, `{"result":false,"error":"Retrieve beanstalkd job ID error"}`)
		return sampleName, sampleTTR, string(body), errors.New("Retrieve beanstalkd job ID error")
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		_, _ = io.WriteString(w, `{"result":false,"error":"Connect to beanstal server fail"}`)
		return sampleName, sampleTTR, string(body), errors.New("Connect to beanstal server fail")
	}
	body, err = bstkConn.Peek(uint64(jobID))
	if err != nil {
		_, _ = io.WriteString(w, `{"result":false,"error":"Read beanstalkd job content fail"}`)
		return sampleName, sampleTTR, string(body), errors.New("Read beanstalkd job content fail")
	}
	// Read beanstalkd job TTR in job ID.
	jobStats, err := bstkConn.StatsJob(uint64(jobID))
	if err != nil {
		_, _ = io.WriteString(w, `{"result":false,"error":"Read beanstalkd job stats fail"}`)
		return sampleName, sampleTTR, string(body), errors.New("Read beanstalkd job stats fail")
	}
	sampleTTR, err = strconv.Atoi(jobStats["ttr"])
	if err != nil {
		_, _ = io.WriteString(w, `{"result":false,"error":"Read beanstalkd job TTR fail"}`)
		return sampleName, sampleTTR, string(body), errors.New("Read beanstalkd job TTR fail")
	}
	bstkConn.Close()
	return sampleName, sampleTTR, string(body), nil
}

// addSampleTube provide a method add a sample job tube in global config
// variable.
func addSampleTube(tube string, key string) {
	for k, v := range sampleJobs.Tubes {
		if v.Name == tube {
			sampleJobs.Tubes[k].Keys = append(sampleJobs.Tubes[k].Keys, key)
			return
		}
	}
	sampleJobs.Tubes = append(sampleJobs.Tubes, SampleTube{
		Name: tube,
		Keys: []string{key},
	})
}

// checkSampleJobs check if exists of sample job by given name.
func checkSampleJobs(name string) bool {
	for _, v := range sampleJobs.Jobs {
		if v.Name == name {
			return true
		}
	}
	return false
}

// saveSample provide a method to storage sample job in config file.
func saveSample() error {
	sampleJobsTOML, err := json.Marshal(sampleJobs)
	if err != nil {
		return err
	}
	pubConf.Sample.Storage = string(sampleJobsTOML)
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(pubConf); err != nil {
		return err
	}

	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		err := ioutil.WriteFile(ConfigFile, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(ConfigFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_, _ = buf.WriteTo(file)
	buf.Reset()
	return nil
}

// getSampleJobDataByKey return sample job body by given key.
func getSampleJobDataByKey(key string) string {
	var data string
	for _, j := range sampleJobs.Jobs {
		if j.Key != key {
			continue
		}
		data = j.Data
	}
	return data
}

// getSampleJobNameByKey return sample job name by given key.
func getSampleJobNameByKey(key string) string {
	var data string
	for _, j := range sampleJobs.Jobs {
		if j.Key != key {
			continue
		}
		data = j.Name
	}
	return data
}

// getSampleJobTtrByKey return sample job TTR by given key.
func getSampleJobTtrByKey(key string) int {
	var data = DefaultTTR
	for _, j := range sampleJobs.Jobs {
		if j.Key != key {
			continue
		}
		data = j.TTR
	}
	return data
}

// deleteSamples drop sample job by given key.
func deleteSamples(key string) {
	if key == "" {
		return
	}

	for k, j := range sampleJobs.Jobs {
		if j.Key == key {
			sampleJobs.Jobs = sampleJobs.Jobs[:k+copy(sampleJobs.Jobs[k:], sampleJobs.Jobs[k+1:])]
		}
	}
	for k, v := range sampleJobs.Tubes {
		for i, t := range v.Keys {
			if t == key {
				sampleJobs.Tubes[k].Keys = sampleJobs.Tubes[k].Keys[:i+copy(sampleJobs.Tubes[k].Keys[i:], sampleJobs.Tubes[k].Keys[i+1:])]
			}
		}
	}
	_ = saveSample()
}

// loadSample puts a job into tube by given sample job key.
func loadSample(server string, tube string, key string) {
	var err error
	var bstkConn *beanstalk.Conn
	data := getSampleJobDataByKey(key)
	if data == "" {
		return
	}
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return
	}
	bstkTube := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	_, _ = bstkTube.Put([]byte(data), uint32(DefaultPriority), time.Duration(DefaultDelay)*time.Second, time.Duration(getSampleJobTtrByKey(key))*time.Second)
	bstkConn.Close()
}

// newSample provide method to add a sample job.
func newSample(server string, f url.Values, w http.ResponseWriter, r *http.Request) {
	var err error
	var key = randToken()
	var name, body, ttr string
	var sampleTTR int
	var tubes []string
	alert := `<div class="alert alert-danger" id="sjsa"><button type="button" class="close" onclick="$('#sjsa').fadeOut('fast');">×</button><span> Required fields are not set</span></div>`
	err = readConf()
	if err != nil {
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", `<div class="alert alert-danger" id="sjsa"><button type="button" class="close" onclick="$('#sjsa').fadeOut('fast');">×</button><span> Read config error</span></div>`), server))
		return
	}
	for k, v := range f {
		switch k {
		case "jobdata":
			body = v[0]
		case "name":
			name = v[0]
		case "ttr":
			ttr = v[0]
		case "action", "key":
			continue
		default:
			t := strings.TrimSuffix(strings.TrimPrefix(k, `tubes[`), `]`)
			tubes = append(tubes, t)
			addSampleTube(t, key)
		}
	}
	if len(tubes) == 0 || name == "" || body == "" || ttr == "" {
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", alert), server))
		return
	}
	if checkSampleJobs(name) {
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", `<div class="alert alert-danger" id="sjsa"><button type="button" class="close" onclick="$('#sjsa').fadeOut('fast');">×</button><span> You already have a job with this name</span></div>`), server))
		return
	}
	sampleTTR, err = strconv.Atoi(ttr)
	if err != nil {
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", `<div class="alert alert-danger" id="sjsa"><button type="button" class="close" onclick="$('#sjsa').fadeOut('fast');">×</button><span> You should give a correct TTR with this sample</span></div>`), server))
		return
	}
	sampleJobs.Jobs = append(sampleJobs.Jobs, SampleJob{
		Key:   key,
		Name:  name,
		Tubes: tubes,
		Data:  body,
		TTR:   sampleTTR,
	})
	err = saveSample()
	if err != nil {
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", `<div class="alert alert-danger" id="sjsa"><button type="button" class="close" onclick="$('#sjsa').fadeOut('fast');">×</button><span> Save sample job error</span></div>`), server))
		return
	}
	w.Header().Set("Location", "./sample?action=manageSamples")
	w.WriteHeader(307)
}

// editSample provide method to update a sample job.
func editSample(server string, f url.Values, key string, w http.ResponseWriter, r *http.Request) {
	deleteSamples(key)
	newSample(server, f, w, r)
}

// getSampleJobList render a table of sample job.
func getSampleJobList() string {
	if len(sampleJobs.Jobs) == 0 {
		return `<div class="clearfix"><div class="pull-left">There are no saved jobs.</div><div class="pull-right"><a href="?action=newSample" class="btn btn-default btn-sm"><i class="glyphicon glyphicon-plus"></i> Add job to samples</a></div></div>`
	}
	var tr, td, serverList, buf strings.Builder
	for _, j := range sampleJobs.Jobs {
		for _, v := range j.Tubes {
			for _, s := range selfConf.Servers {
				serverList.Reset()
				serverList.WriteString(`<li><a href="./tube?server=`)
				serverList.WriteString(s)
				serverList.WriteString(`&tube=`)
				serverList.WriteString(v)
				serverList.WriteString(`&action=loadSample&key=`)
				serverList.WriteString(j.Key)
				serverList.WriteString(`&redirect=`)
				serverList.WriteString(url.QueryEscape(`tube?action=manageSamples`))
				serverList.WriteString(`">`)
				serverList.WriteString(s)
				serverList.WriteString(`</a></li>`)
			}
			td.WriteString(` <div class="btn-group"><a class="btn btn-default btn-sm" href="#" data-toggle="dropdown"><i class="glyphicon glyphicon-forward"></i> `)
			td.WriteString(html.EscapeString(v))
			td.WriteString(`</a><button class="btn btn-default btn-sm dropdown-toggle" data-toggle="dropdown"><span class="caret"></span></button><ul class="dropdown-menu">`)
			td.WriteString(serverList.String())
			td.WriteString(`</ul></div>`)
		}
		tr.WriteString(`<tr><td style="line-height: 25px !important;"><a href="?action=editSample&key=`)
		tr.WriteString(j.Key)
		tr.WriteString(`">`)
		tr.WriteString(html.EscapeString(j.Name))
		tr.WriteString(`</a></td><td>`)
		tr.WriteString(td.String())
		tr.WriteString(`</td><td><div class="pull-right"><a class="btn btn-default btn-sm" href="?action=editSample&key=`)
		tr.WriteString(j.Key)
		tr.WriteString(`"><i class="glyphicon glyphicon-pencil"></i> Edit</a> <a class="btn btn-default btn-sm" href="?action=deleteSample&key=`)
		tr.WriteString(j.Key)
		tr.WriteString(`"><i class="glyphicon glyphicon-trash"></i> Delete</a></div></td></tr>`)
		td.Reset()
	}
	buf.WriteString(`<div class="clearfix"><div class="pull-right"><a href="?action=newSample" class="btn btn-default btn-sm"><i class="glyphicon glyphicon-plus"></i> Add job to samples</a></div></div><section id="summaryTable"><table class="table table-striped table-hover"><thead><tr><th>Name</th><th>Kick job to tubes</th><th></th></tr></thead><tbody>`)
	buf.WriteString(tr.String())
	buf.WriteString(`</tbody></table></section>`)
	return buf.String()
}
