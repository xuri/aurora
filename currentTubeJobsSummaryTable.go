package main

import (
	"fmt"

	"github.com/kr/beanstalk"
)

// currentTubeJobsSummaryTable constructs a tube job table based on the given server and tube conf.
func currentTubeJobsSummaryTable(server string, tube string) string {
	var err error
	var th, tr string
	var bstkConn *beanstalk.Conn
	if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
		for _, v := range selfConf.TubeFilters {
			th += fmt.Sprintf(`<th name="%s">%s</th>`, v, v)
		}
		return fmt.Sprintf(`<section id="summaryTable"> <div class="row"> <div class="col-sm-12"> <table class="table table-striped table-hover"> <thead> <tr> <th>name</th>%s</tr> </thead> <tbody> </tbody> </table> </div> </div> </section>`, th)
	}
	tubes, _ := bstkConn.ListTubes()
	for _, v := range selfConf.TubeFilters {
		th += fmt.Sprintf(`<th name="%s">%s</th>`, v, v)
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
		var td string
		for _, stats := range selfConf.TubeFilters {
			td += fmt.Sprintf(`<td>%s</td>`, statsMap[stats])
		}
		tr += fmt.Sprintf(`<tr><td name="pause-time-left">%s</td>%s</tr>`, v, td)
	}
	bstkConn.Close()
	template := fmt.Sprintf(`<section id="summaryTable"> <div class="row"> <div class="col-sm-12"> <table class="table table-striped table-hover"> <thead> <tr> <th>name</th>%s</tr> </thead> <tbody> %s </tbody> </table> </div> </div> </section>`, th, tr)
	return template
}
