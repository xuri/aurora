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
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// readConf read external config file when program startup.
func readConf() error {
	buf := new(strings.Builder)
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		err := ioutil.WriteFile(ConfigFile, []byte(ConfigFileTemplate), 0644)
		if err != nil {
			return err
		}
	}
	buf.Reset()
	tomlData, err := os.Open(ConfigFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(buf, tomlData)
	if err != nil {
		return err
	}
	tomlData.Close()
	if _, err := toml.Decode(buf.String(), &pubConf); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(pubConf.Sample.Storage), &sampleJobs); err != nil {
		return err
	}
	parseConf()
	return nil
}

// parseConf parse server config in external config file.
func parseConf() {
	selfConf.Servers = append(selfConf.Servers, pubConf.Servers...)
}

// removeArrayDuplicates provide a function remove duplicates value elements in
// a slice.
func removeArrayDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// removeArrayEmpty provide a function remove empty value elements in a slice.
func removeArrayEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// removeServerInConfig provide a method to remove property in config by given
// field.
func removeServerInConfig(server string) {
	for k, v := range selfConf.Servers {
		if v == server {
			selfConf.Servers = selfConf.Servers[:k+copy(selfConf.Servers[k:], selfConf.Servers[k+1:])]
		}
	}
}

// runCmd run command opens a new browser window pointing to url.
func runCmd(prog string, args ...string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	return cmd.Run()
}

// checkInSlice return bool type value to check if exits in slice by given
// string.
func checkInSlice(list []string, value string) bool {
	set := make(map[string]bool)
	for _, v := range list {
		set[v] = true
	}
	return set[value]
}

// prettyJSON provide method get JSON string with indent.
func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		return b
	}
	return out.Bytes()
}

// base64Decode provide method get Base64 decode string.
func base64Decode(b string) string {
	data, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return string(b)
	}
	return string(data)
}

// preformat provide method get job body after format with config.
func preformat(jobBody []byte) string {
	var job = string(jobBody)
	if selfConf.IsDisabledJSONDecode != 1 {
		job = string(prettyJSON(jobBody))
	}
	if selfConf.IsEnabledBase64Decode != 0 {
		job = base64Decode(job)
	}
	job = html.EscapeString(job)
	return job
}

// parseFlags parse flags of program.
func parseFlags() {
	configPtr := flag.String("c", "", "Use config file.")
	verPtr := flag.Bool("v", false, "Output version and exit.")
	helpPtr := flag.Bool("h", false, "Output this help and exit.")
	flag.Parse()
	if *configPtr == "" {
		selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			os.Exit(0)
		}
		ConfigFile = selfDir + string(os.PathSeparator) + `aurora.toml`
	} else {
		ConfigFile = *configPtr
	}
	if *verPtr {
		fmt.Printf("aurora version: %.1f\r\n", Version)
		os.Exit(0)
	}
	if *helpPtr {
		fmt.Printf("aurora version: %.1f\r\nCopyright (c) 2016 - 2020 Ri Xu https://xuri.me All rights reserved.\r\n\r\nUsage: aurora [OPTIONS] [cmd [arg ...]]\n  -c <filename>   Use config file. (default: aurora.toml)\r\n  -h \t\t  Output this help and exit.\r\n  -v \t\t  Output version and exit.\r\n", Version)
		os.Exit(0)
	}
}

// basicAuth provide a simple method to HTTP authenticate.
func basicAuth(f ViewFunc) ViewFunc {
	if !pubConf.Auth.Enabled {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		basicAuthPrefix := "Basic "
		// Parse request header
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Decoding authentication information.
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], []byte(pubConf.Auth.Username)) &&
					bytes.Equal(pair[1], []byte(pubConf.Auth.Password)) {
					f(w, r)
					return
				}
			}
		}
		// Authorization fail, return 401 Unauthorized.
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// randToken generate a random token with MD5.
func randToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// setHeader provide common method set HTTP header response.
func setHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "WebServer")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
}

// checkUpdate render update notice alert.
func checkUpdate() string {
	if updateInfo != "uncheck" {
		return updateInfo
	}
	updateInfo = ""
	r, err := http.Get(UpdateURL)
	if err != nil {
		return updateInfo
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		r.Body.Close()
		return updateInfo
	}
	r.Body.Close()
	u := UpdateTags{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		return updateInfo
	}
	if len(u) < 1 {
		return updateInfo
	}
	v, err := strconv.ParseFloat(u[0].Name, 64)
	if err != nil {
		return updateInfo
	}
	if Version < v {
		updateInfo = fmt.Sprintf(`<br/><div class="alert alert-info" style="position: relative;top:50px;"><span>You are currently running version %.1f of aurora. A new version is available: <b>%.1f</b> Get it from <b><a href="https://github.com/xuri/aurora" target="_blank">GitHub</a></b></span></div>`, Version, v)
	}
	return updateInfo
}
