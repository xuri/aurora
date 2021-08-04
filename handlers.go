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
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// handlerMain handle request on router: /
func handlerMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go WebServer")
	w.Header().Set("Content-Type", "text/html")
	server := r.URL.Query().Get("server")
	readCookies(r)
	_, _ = io.WriteString(w, tplMain(getServerStatus(), server))
}

// handlerServerList handle request on router: /index
func handlerServerList(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	_, _ = io.WriteString(w, getServerStatus())
}

// serversRemove handle request on router: /serversRemove
func serversRemove(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	server := r.URL.Query().Get("removeServer")
	removeServerInCookie(server, w, r)
	removeServerInConfig(server)
	w.Header().Set("Location", "./public")
	w.WriteHeader(307)
}

// handlerServer handle request on router: /server
func handlerServer(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	server := html.EscapeString(r.URL.Query().Get("server"))
	action := r.URL.Query().Get("action")
	switch action {
	case "reloader":
		_, _ = io.WriteString(w, getServerTubes(server))
		return
	case "clearTubes":
		_ = r.ParseForm()
		clearTubes(server, r.Form)
		_, _ = io.WriteString(w, `{"result":true}`)
		return
	}
	_, _ = io.WriteString(w, tplServer(getServerTubes(server), server))
}

// handlerTube handle request on router: /tube
func handlerTube(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	server := html.EscapeString(r.URL.Query().Get("server"))
	tube := html.EscapeString(r.URL.Query().Get("tube"))
	action := r.URL.Query().Get("action")
	count := html.EscapeString(r.URL.Query().Get("count"))
	switch action {
	case "addjob":
		addJob(server, r.PostFormValue("tubeName"), r.PostFormValue("tubeData"), r.PostFormValue("tubePriority"), r.PostFormValue("tubeDelay"), r.PostFormValue("tubeTtr"))
		_, _ = io.WriteString(w, `{"result":true}`)
		return
	case "search":
		content := searchTube(server, tube, html.EscapeString(r.URL.Query().Get("limit")), html.EscapeString(r.URL.Query().Get("searchStr")))
		_, _ = io.WriteString(w, tplTube(content, server, tube))
		return
	case "addSample":
		_ = r.ParseForm()
		addSample(server, r.Form, w)
		return
	default:
		handleRedirect(w, r, server, tube, action, count)
	}
}

// handleRedirect handle request with redirect response.
func handleRedirect(w http.ResponseWriter, r *http.Request, server string, tube string, action string, count string) {
	var link strings.Builder
	link.WriteString(`./tube?server=`)
	link.WriteString(server)
	link.WriteString(`&tube=`)
	switch action {
	case "kick":
		kick(server, tube, count)
		link.WriteString(url.QueryEscape(tube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	case "kickJob":
		kickJob(server, tube, r.URL.Query().Get("jobid"))
		link.WriteString(url.QueryEscape(tube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	case "pause":
		pause(server, tube, count)
		link.WriteString(url.QueryEscape(tube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	case "moveJobsTo":
		destTube := tube
		if r.URL.Query().Get("destTube") != "" {
			destTube = r.URL.Query().Get("destTube")
		}
		moveJobsTo(server, tube, destTube, r.URL.Query().Get("state"), r.URL.Query().Get("destState"))
		link.WriteString(url.QueryEscape(destTube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	case "deleteAll":
		deleteAll(server, tube)
		link.WriteString(url.QueryEscape(tube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	case "deleteJob":
		deleteJob(server, tube, r.URL.Query().Get("jobid"))
		link.WriteString(url.QueryEscape(tube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	case "loadSample":
		loadSample(server, tube, r.URL.Query().Get("key"))
		link.WriteString(url.QueryEscape(tube))
		w.Header().Set("Location", link.String())
		w.WriteHeader(307)
	}
	_, _ = io.WriteString(w, tplTube(currentTube(server, tube), server, tube))
}

// handlerSample handle request on router: /sample
func handlerSample(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	action := r.URL.Query().Get("action")
	server := html.EscapeString(r.URL.Query().Get("server"))
	switch action {
	case "manageSamples":
		_, _ = io.WriteString(w, tplSampleJobsManage(getSampleJobList(), server))
		return
	case "newSample":
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", ""), server))
		return
	case "editSample":
		_, _ = io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit(html.EscapeString(r.URL.Query().Get("key")), ""), server))
		return
	case "actionNewSample":
		_ = r.ParseForm()
		newSample(server, r.Form, w, r)
		return
	case "actionEditSample":
		_ = r.ParseForm()
		editSample(server, r.Form, r.URL.Query().Get("key"), w, r)
		return
	case "deleteSample":
		deleteSamples(r.URL.Query().Get("key"))
		w.Header().Set("Location", "./sample?action=manageSamples")
		w.WriteHeader(307)
		return
	}
}

// handlerStatistics handle request on router: /statistics
func handlerStatistics(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	action := r.URL.Query().Get("action")
	server := html.EscapeString(r.URL.Query().Get("server"))
	tube := html.EscapeString(r.URL.Query().Get("tube"))
	switch action {
	case "preference":
		_, _ = io.WriteString(w, tplStatisticSetting(tplStatisticEdit("")))
		return
	case "save":
		_ = r.ParseForm()
		statisticPreferenceSave(r.Form, w, r)
		return
	case "reloader":
		_, _ = io.WriteString(w, statisticWaitress(server, tube))
		return
	}
	_, _ = io.WriteString(w, tplStatistic(server, tube))
}
