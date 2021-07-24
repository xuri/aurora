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
	"container/list"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/aurora/beanstalk"
)

// statisticPreferenceSave provide method to save statistics preference
// settings.
func statisticPreferenceSave(f url.Values, w http.ResponseWriter, r *http.Request) {
	var err error
	var collection, frequency string
	var tubes []string
	alert := `<div class="alert alert-danger" id="sfsa"><button type="button" class="close" onclick="$('#sfsa').fadeOut('fast');">×</button><span> Required fields are not set correct</span></div>`
	err = readConf()
	if err != nil {
		_, _ = io.WriteString(w, tplStatisticSetting(tplStatisticEdit(`<div class="alert alert-danger"><button type="button" class="close" data-dismiss="alert">×</button><span> Read config error</span></div>`)))
		return
	}
	for k, v := range f {
		switch k {
		case "frequency":
			frequency = v[0]
		case "collection":
			collection = v[0]
		case "action":
			continue
		default:
			t := strings.TrimSuffix(strings.TrimPrefix(k, `tubes[`), `]`)
			tubes = append(tubes, t)
		}
	}
	if len(tubes) == 0 || collection == "" || frequency == "" {
		_, _ = io.WriteString(w, tplStatisticSetting(tplStatisticEdit(alert)))
		return
	}
	err = statisticCashier(collection, frequency, tubes)
	if err != nil {
		_, _ = io.WriteString(w, tplStatisticSetting(tplStatisticEdit(`<div class="alert alert-danger" id="sfsa"><button type="button" class="close" onclick="$('#sfsa').fadeOut('fast');">×</button><span> Save statistics preference error</span></div>`)))
		return
	}
	_, _ = io.WriteString(w, tplStatisticSetting(tplStatisticEdit(`<div class="alert alert-success" id="sfsa"><button type="button" class="close" onclick="$('#sfsa').fadeOut('fast');">×</button><span> Statistics preference saved</span></div>`)))
}

// statisticCashier validate collection and frequency parameter and send notify
// to statistic Goroutine that the configuration of statistics preference
// settings has changed.
func statisticCashier(collection string, frequency string, tubes []string) error {
	c, err := strconv.Atoi(collection)
	if err != nil {
		return err
	}
	f, err := strconv.Atoi(frequency)
	if err != nil {
		return err
	}
	if c < 1 {
		c = 0
	}
	if f < 1 {
		f = 1
	}
	selfConf.StatisticsCollection = c
	selfConf.StatisticsFrequency = f
	statisticsDataServer = make(map[string]map[string]map[string]*list.List)
	for _, v := range tubes {
		addr := strings.Split(v, `:`)
		if len(addr) != 3 {
			continue
		}
		tube := make(map[string]map[string]*list.List)
		tube[addr[2]] = make(map[string]*list.List)
		s, ok := statisticsDataServer[addr[0]+`:`+addr[1]]
		if !ok {
			statisticsDataServer[addr[0]+`:`+addr[1]] = tube
		} else {
			s[addr[2]] = tube[addr[2]]
		}
	}
	statisticsData.Lock()
	statisticsData.Server = statisticsDataServer
	statisticsData.Unlock()
	notify <- true
	return nil
}

// statistic provide method to control statisticAgent collect the statistics
// data in a Goroutine.
func statistic() {
	for {
		tick := time.Tick(time.Duration(selfConf.StatisticsFrequency) * time.Second)
	NOTIFY:
		for {
			select {
			case <-notify:
				break NOTIFY
			case <-tick:
				for k, v := range statisticsData.Server {
					for t := range v {
						if selfConf.StatisticsCollection == 0 {
							continue
						}
						err := statisticAgent(k, t)
						if err != nil {
							continue
						}
					}
				}
			}
		}
	}
}

// statisticAgent collect the statistics data by given server and tube.
func statisticAgent(server string, tube string) error {
	var err error
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		return err
	}
	tubeStats := &beanstalk.Tube{
		Conn: bstkConn,
		Name: tube,
	}
	statsMap, err := tubeStats.Stats()
	if err != nil {
		bstkConn.Close()
		return err
	}
	for _, field := range statisticsFields {
		for k, v := range field {
			t := time.Now()
			stats, err := strconv.Atoi(statsMap[v])
			if err != nil {
				continue
			}
			statisticsData.Lock()
			_, ok := statisticsData.Server[server][tube][k]
			if !ok {
				statisticsData.Server[server][tube][k] = list.New()
			}
			if statisticsData.Server[server][tube][k].Len() >= selfConf.StatisticsCollection {
				front := statisticsData.Server[server][tube][k].Back()
				statisticsData.Server[server][tube][k].Remove(front)
			}
			statisticsData.Server[server][tube][k].PushFront([]int{t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), stats})
			statisticsData.Unlock()
		}
	}
	bstkConn.Close()
	return nil
}

// statisticWaitress return real-time statistics data by given server and tube.
func statisticWaitress(server string, tube string) string {
	var buf, b, s, l strings.Builder
	b.WriteString(`{`)
	statisticsData.RLock()
	for _, field := range statisticsFields {
		for k := range field {
			b.WriteString(`"`)
			b.WriteString(k)
			b.WriteString(`":[`)
			_, ok := statisticsData.Server[server][tube][k]
			if !ok {
				b.WriteString(`],`)
				continue
			}
			s.Reset()
			for e := statisticsData.Server[server][tube][k].Front(); e != nil; e = e.Next() {
				s.WriteString(`[`)
				l.Reset()
				for _, v := range e.Value.([]int) {
					l.WriteString(strconv.Itoa(v))
					l.WriteString(`,`)
				}
				s.WriteString(strings.TrimSuffix(l.String(), `,`))
				s.WriteString(`],`)
			}
			b.WriteString(strings.TrimSuffix(s.String(), `,`))
			b.WriteString(`],`)
		}
	}
	statisticsData.RUnlock()
	buf.WriteString(strings.TrimSuffix(b.String(), `,`))
	buf.WriteString(`}`)
	return buf.String()
}
