package main

import (
	"bytes"

	"github.com/kr/beanstalk"
)

// currentTubeJobsSummaryTable constructs a tube job table based on the given server and tube conf.
func currentTubeJobsSummaryTable(server string, tube string) string {
	var err error
	var th, tr, td, template bytes.Buffer
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		for _, v := range selfConf.TubeFilters {
			th.WriteString(`<th>`)
			th.WriteString(v)
			th.WriteString(`</th>`)
		}
		buf := bytes.Buffer{}
		buf.WriteString(`<section id="summaryTable"><div class="row"><div class="col-sm-12"><table class="table table-striped table-hover"><thead><tr><th>name</th>`)
		buf.WriteString(th.String())
		buf.WriteString(`</tr></thead><tbody></tbody></table></div></div></section>`)
		return buf.String()
	}
	tubes, _ := bstkConn.ListTubes()
	for _, v := range selfConf.TubeFilters {
		th.WriteString(`<th>`)
		th.WriteString(v)
		th.WriteString(`</th>`)
	}
	for _, v := range tubes {
		if v != tube {
			continue
		}
		tubeStats := &beanstalk.Tube{Conn: bstkConn, Name: v}
		statsMap, err := tubeStats.Stats()
		if err != nil {
			continue
		}
		for _, stats := range selfConf.TubeFilters {
			td.WriteString(`<td>`)
			td.WriteString(statsMap[stats])
			td.WriteString(`</td>`)
		}
		tr.WriteString(`<tr><td>`)
		tr.WriteString(v)
		tr.WriteString(`</td>`)
		tr.WriteString(td.String())
		tr.WriteString(`</tr>`)
		td.Reset()
	}
	bstkConn.Close()
	template.WriteString(`<section id="summaryTable"><div class="row"><div class="col-sm-12"><table class="table table-striped table-hover"><thead><tr><th>name</th>`)
	template.WriteString(th.String())
	template.WriteString(`</tr></thead><tbody> `)
	template.WriteString(tr.String())
	template.WriteString(`</tbody></table></div></div></section>`)
	if tr.String() == `` {
		return ``
	}
	return template.String()
}
