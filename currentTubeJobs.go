package main

// currentTubeJobs call currentTubeJobsSummaryTable, currentTubeJobsActionsRow and currentTubeJobsShowcase functions by given server and tube config, and merge these functions return value.
func currentTubeJobs(server string, tube string) string {
	return currentTubeJobsSummaryTable(server, tube) + currentTubeJobsActionsRow(server, tube) + currentTubeJobsShowcase(server, tube)
}
