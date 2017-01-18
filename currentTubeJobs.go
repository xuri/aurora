package main

import "bytes"

// currentTubeJobs call currentTubeJobsSummaryTable, currentTubeJobsActionsRow
// and currentTubeJobsShowcase functions by given server and tube config, and
// merge these functions return value.
func currentTubeJobs(server string, tube string) string {
	var table = currentTubeJobsSummaryTable(server, tube)
	buf := bytes.Buffer{}
	if table == `` {
		buf.WriteString(`Tube "`)
		buf.WriteString(tube)
		buf.WriteString(`" not found or it is empty <br><br><a href="./server?server=`)
		buf.WriteString(server)
		buf.WriteString(`"> &lt;&lt; back </a>`)
		return buf.String()
	}
	buf.WriteString(table)
	buf.WriteString(currentTubeJobsActionsRow(server, tube))
	buf.WriteString(currentTubeJobsShowcase(server, tube))
	return buf.String()
}
