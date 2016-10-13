package main

// currentTubeJobs call currentTubeJobsSummaryTable, currentTubeJobsActionsRow and currentTubeJobsShowcase functions by given server and tube config, and merge these functions return value.
func currentTubeJobs(server string, tube string) string {
	var table = currentTubeJobsSummaryTable(server, tube)
	if table == `` {
		return `Tube "` + tube + `" not found or it is empty <br><br><a href="./server?server=` + server + `"> << back </a>`
	}
	return table + currentTubeJobsActionsRow(server, tube) + currentTubeJobsShowcase(server, tube)
}
