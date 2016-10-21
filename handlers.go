package main

import (
	"fmt"
	"io"
	"net/http"
)

// handlerMain handle request on router: /
func handlerMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go WebServer")
	w.Header().Set("Content-Type", "text/html")
	server := r.URL.Query().Get("server")
	readCookies(r)
	io.WriteString(w, tplMain(getServerStatus(), server))
}

// handlerServerList handle request on router: /index
func handlerServerList(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	io.WriteString(w, getServerStatus())
}

// serversRemove handle request on router: /serversRemove
func serversRemove(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	server := r.URL.Query().Get("removeServer")
	removeServerInCookie(server, w, r)
	removeServerInConfig(server)
	http.Redirect(w, r, "/public", 301)
}

// handlerServer handle request on router: /server
func handlerServer(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	server := r.URL.Query().Get("server")
	action := r.URL.Query().Get("action")
	switch action {
	case "reloader":
		io.WriteString(w, getServerTubes(server))
		return
	case "clearTubes":
		r.ParseForm()
		clearTubes(server, r.Form)
		io.WriteString(w, `{"result":true}`)
		return
	}
	io.WriteString(w, tplServer(getServerTubes(server), server))
}

// handlerTube handle request on router: /tube
func handlerTube(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	server := r.URL.Query().Get("server")
	tube := r.URL.Query().Get("tube")
	action := r.URL.Query().Get("action")
	count := r.URL.Query().Get("count")
	switch action {
	case "addjob":
		addJob(server, r.PostFormValue("tubeName"), r.PostFormValue("tubeData"), r.PostFormValue("tubePriority"), r.PostFormValue("tubeDelay"), r.PostFormValue("tubeTtr"))
		io.WriteString(w, `{"result":true}`)
		return
	case "search":
		content := searchTube(server, tube, r.URL.Query().Get("limit"), r.URL.Query().Get("searchStr"))
		io.WriteString(w, tplTube(content, server, tube))
		return
	case "addSample":
		r.ParseForm()
		addSample(server, r.Form, w)
		return
	default:
		handleRedirect(w, r, server, tube, action, count)
	}
}

// handleRedirect handle request with redirect response.
func handleRedirect(w http.ResponseWriter, r *http.Request, server string, tube string, action string, count string) {
	switch action {
	case "kick":
		kick(server, tube, count)
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, tube), 302)
	case "kickJob":
		kickJob(server, tube, r.URL.Query().Get("jobid"))
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, tube), 302)
	case "pause":
		pause(server, tube, count)
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, tube), 302)
	case "moveJobsTo":
		destTube := tube
		if r.URL.Query().Get("destTube") != "" {
			destTube = r.URL.Query().Get("destTube")
		}
		moveJobsTo(server, tube, destTube, r.URL.Query().Get("state"), r.URL.Query().Get("destState"))
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, destTube), 302)
	case "deleteAll":
		deleteAll(server, tube)
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, tube), 302)
	case "deleteJob":
		deleteJob(server, tube, r.URL.Query().Get("jobid"))
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, tube), 302)
	case "loadSample":
		loadSample(server, tube, r.URL.Query().Get("key"))
		http.Redirect(w, r, fmt.Sprintf("/tube?server=%s&tube=%s", server, tube), 302)
	}
	io.WriteString(w, tplTube(currentTube(server, tube), server, tube))
}

// handlerSample handle request on router /sample
func handlerSample(w http.ResponseWriter, r *http.Request) {
	setHeader(w, r)
	readCookies(r)
	action := r.URL.Query().Get("action")
	server := r.URL.Query().Get("server")
	switch action {
	case "manageSamples":
		io.WriteString(w, tplSampleJobsManage(getSampleJobList(), server))
		return
	case "newSample":
		io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit("", ""), server))
		return
	case "editSample":
		io.WriteString(w, tplSampleJobsManage(tplSampleJobEdit(r.URL.Query().Get("key"), ""), server))
		return
	case "actionNewSample":
		r.ParseForm()
		newSample(server, r.Form, w, r)
		return
	case "actionEditSample":
		r.ParseForm()
		editSample(server, r.Form, r.URL.Query().Get("key"), w, r)
		return
	case "deleteSample":
		deleteSamples(r.URL.Query().Get("key"))
		http.Redirect(w, r, "/sample?action=manageSamples", 301)
		return
	}
}
