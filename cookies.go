package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// readCookies read config property storage in cookie.
func readCookies(r *http.Request) {
	var servers, filters, tubeFilters []string
	var tubeSelectorValue string
	// Read servers in cookies
	beansServers, err := r.Cookie("beansServers")
	if err == nil {
		beansServersValue, _ := url.QueryUnescape(beansServers.Value)
		servers = strings.Split(beansServersValue, `;`)
	}
	// Read Filter in cookies
	filter, err := r.Cookie("filter")
	if err == nil {
		filterValue, _ := url.QueryUnescape(filter.Value)
		filters = strings.Split(filterValue, `,`)
		filters = removeArrayDuplicates(removeArrayEmpty(filters))
	} else {
		filters = []string{"current-connections", "current-jobs-buried", "current-jobs-delayed", "current-jobs-ready", "current-jobs-reserved", "current-jobs-urgent", "current-tubes"}
	}
	for _, v := range servers {
		selfConf.Servers = append(selfConf.Servers, v)
	}
	// Read Tube Filter in cookies
	tubeFilter, err := r.Cookie("tubefilter")
	if err == nil {
		tubeFilterValue, _ := url.QueryUnescape(tubeFilter.Value)
		tubeFilters = strings.Split(tubeFilterValue, `,`)
		tubeFilters = removeArrayDuplicates(removeArrayEmpty(tubeFilters))
	} else {
		tubeFilters = []string{"current-jobs-urgent", "current-jobs-ready", "current-jobs-reserved", "current-jobs-delayed", "current-jobs-buried", "total-jobs"}
	}
	tubeSelector, err := r.Cookie("tubeSelector")
	if err != nil {
		tubeSelectorValue = ""
	} else {
		tubeSelectorValue = tubeSelector.Value
	}

	selfConf.Servers = removeArrayDuplicates(removeArrayEmpty(selfConf.Servers))
	selfConf.Filter = filters
	selfConf.TubeFilters = tubeFilters
	selfConf.TubePauseSeconds = readIntCookie(r, `tubePauseSeconds`, 3600)
	selfConf.IsDisabledJSONDecode = readIntCookie(r, `isDisabledJsonDecode`, 0)
	selfConf.IsDisabledUnserialization = readIntCookie(r, `isDisabledUnserialization`, 0)
	selfConf.IsDisabledJobDataHighlight = readIntCookie(r, `isDisabledJobDataHighlight`, 0)
	selfConf.IsEnabledBase64Decode = readIntCookie(r, `isEnabledBase64Decode`, 0)
	selfConf.TubePauseSeconds = readIntCookie(r, `tubePauseSeconds`, -1)
	selfConf.AutoRefreshTimeoutMs = readIntCookie(r, `autoRefreshTimeoutMs`, 500)
	selfConf.SearchResultLimit = readIntCookie(r, `searchResultLimit`, 25)
	selfConf.TubeSelector = tubeSelectorValue
}

// readIntCookie return int value by the given string.
func readIntCookie(r *http.Request, name string, defaultValue int) int {
	cookie, err := r.Cookie(name)
	if err == nil {
		value, err := strconv.Atoi(cookie.Value)
		if err == nil {
			return value
		}
	}
	return defaultValue
}

// removeServerInCookie remove field in cookie by the given string.
func removeServerInCookie(server string, w http.ResponseWriter, r *http.Request) {
	for k, v := range selfConf.Servers {
		if v == server {
			selfConf.Servers = selfConf.Servers[:k+copy(selfConf.Servers[k:], selfConf.Servers[k+1:])]
		}
	}
	var serverInCookie string
	for _, v := range selfConf.Servers {
		serverInCookie += v + `;`
	}
	cookie := http.Cookie{
		Name:  `beansServers`,
		Value: url.QueryEscape(serverInCookie),
	}
	http.SetCookie(w, &cookie)
}
