package runners

import "strings"

func parseBuildName(name string) (jobId, buildName string) {
	s := strings.Split(name, "__")
	return s[0], s[1]
}

func parseTestName(name string) (jobId, buildName, testName, seed string, groupsName []string) {
	s := strings.Split(name, "__")
	jobId = s[0]
	buildName = s[1]
	groupsName = s[2 : len(s)-2]
	testName = s[len(s)-2]
	seed = s[len(s)-1]
	return
}
