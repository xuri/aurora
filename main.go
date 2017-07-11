//go:generate statik -src=./public

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/rakyll/statik/fs"
	_ "github.com/xuri/aurora/statik"
)

// main function defines the entry point for the program if read config file or
// init failed, the application will be exit.
func main() {
	parseFlags()
	var err error
	err = readConf()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	statikFS, err := fs.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	http.FileServer(statikFS)

	// handle static files include HTML, CSS and JavaScripts.
	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
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
	fmt.Println("To view beanstalk console open", url, "in browser")
	if !pubConf.OpenPage.Enabled {
		return
	}
	var err error
	switch runtime.GOOS {
	case "linux":
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
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
