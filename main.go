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
//
// See https://xuri.me/aurora for more information about aurora.
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

//go:embed public
var staticFiles embed.FS

// main function defines the entry point for the program if read config file or
// init failed, the application will be exit.
func main() {
	parseFlags()
	err := readConf()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
		err = http.ListenAndServe(pubConf.Listen, nil)
		if err != nil {
			fmt.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
	go statistic()
	openPage()
	handleSignals()
}

// openPage function can be open system default browser automatic.
func openPage() {
	url := fmt.Sprintf("http://%v", pubConf.Listen)
	fmt.Println("To view beanstalkd console open", url, "in browser")
	if !pubConf.OpenPage.Enabled {
		return
	}
	var err error
	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd":
		err = runCmd("xdg-open", url)
	case "darwin":
		err = runCmd("open", url)
	case "windows":
		r := strings.NewReplacer("&", "^&")
		err = runCmd("cmd", "/c", "start", r.Replace(url))
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}

// handleSignals handle kill signal.
func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-c
}
