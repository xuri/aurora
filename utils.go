package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// readConf read external config file when program startup.
func readConf() error {
	buf := new(bytes.Buffer)
	selfDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		err := ioutil.WriteFile(selfDir+string(os.PathSeparator)+ConfigFile, []byte(ConfigFileTemplate), 0644)
		if err != nil {
			return err
		}
	}
	buf.Reset()
	tomlData, err := os.Open(filepath.Join(selfDir, ConfigFile))
	if err != nil {
		return err
	}
	io.Copy(buf, tomlData)
	tomlData.Close()
	if _, err := toml.Decode(string(buf.Bytes()), &pubConf); err != nil {
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
	for _, v := range pubConf.Servers {
		selfConf.Servers = append(selfConf.Servers, v)
	}
}

// removeArrayDuplicates provide a function remove duplicates value elements in a slice.
func removeArrayDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
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

// removeServerInConfig provide a method to remove property in config by given field.
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

// checkInSlice return bool type value to check if exits in slice by given string.
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

// parseFlags parse flags of program.
func parseFlags() {
	configPtr := flag.String("c", "", "Use config file.")
	verPtr := flag.Bool("v", false, "Output version and exit.")
	helpPtr := flag.Bool("h", false, "Output this help and exit.")
	flag.Parse()
	if *configPtr == "" {
		ConfigFile = `.` + string(os.PathSeparator) + `aurora.toml`
	} else {
		ConfigFile = *configPtr
	}
	if *verPtr == true {
		fmt.Println("aurora version: 0.1")
		os.Exit(1)
	}
	if *helpPtr == true {
		fmt.Println("aurora version: 0.1\r\nCopyright (c) 2016 Ri Xu https://xuri.me \r\n\r\nUsage: aurora [OPTIONS] [cmd [arg ...]]\n  -c <filename>   Use config file. (default: aurora.toml)\r\n  -h \t\t  Output this help and exit.\r\n  -v \t\t  Output version and exit.")
		os.Exit(1)
	}
}

// basicAuth provide a simple method to HTTP authenticate.
func basicAuth(f ViewFunc) ViewFunc {
	if pubConf.Auth.Enabled == false {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
			return
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
	w.Header().Set("Server", "Go WebServer")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
}
