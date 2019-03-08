package core

import (
	"github.com/shady831213/jarvism/utils"
	"strconv"
)

func statusString(buildPass,
buildFail,
totalBuild,
testPass,
testFail,
testWarning,
testUnknown,
totalTest int) string {
	return utils.Brown("[B:(") + utils.Green("P:"+strconv.Itoa(buildPass)) + utils.Brown("/") + utils.Red("F:"+strconv.Itoa(buildFail)) + utils.Brown("/") +
		utils.Brown("D:"+strconv.Itoa(buildPass+buildFail)+"/T"+strconv.Itoa(totalBuild)+")][T:(") + utils.Green("P:"+strconv.Itoa(testPass)) + utils.Brown("/") +
		utils.Red("F:"+strconv.Itoa(testFail)) + utils.Brown("/") + utils.Yellow("W:"+strconv.Itoa(testWarning)) + utils.Brown("/") + utils.LightRed("U:"+strconv.Itoa(testUnknown)) +
		utils.Brown("/") + utils.Brown("D:"+strconv.Itoa(testPass+testFail+testWarning+testUnknown)+"/T:"+strconv.Itoa(totalTest)+")]")
}

func finishStatusString(buildPass,
buildFail,
totalBuild,
testPass,
testFail,
testWarning,
testUnknown,
totalTest int) string {
	return utils.Brown("[Builds:(") + utils.Green("PASS:"+strconv.Itoa(buildPass)) + utils.Brown("/") + utils.Red("FAIL:"+strconv.Itoa(buildFail)) + utils.Brown("/") +
		utils.Brown("DONE:"+strconv.Itoa(buildPass+buildFail)+"/TOTAL:"+strconv.Itoa(totalBuild)+")][Tests:(") + utils.Green("PASS:"+strconv.Itoa(testPass)) + utils.Brown("/") +
		utils.Red("FAIL:"+strconv.Itoa(testFail)) + utils.Brown("/") + utils.Yellow("WARNING:"+strconv.Itoa(testWarning)) + utils.Brown("/") + utils.LightRed("UNKNOWN:"+strconv.Itoa(testUnknown)) +
		utils.Brown("/") + utils.Brown("DONE:"+strconv.Itoa(testPass+testFail+testWarning+testUnknown)+"/TOTAL:"+strconv.Itoa(totalTest)+")]")
}

func statusMonitor(status *string, totalBuild, totalTest int, buildDone chan error, testDone chan *JVSTestResult, done chan bool) {
	buildPass := 0
	buildFail := 0
	testPass := 0
	testFail := 0
	testWarning := 0
	testUnknown := 0
LableFor:
	for {
		select {
		case err := <-buildDone:
			{
				if err == nil {
					buildPass++
				} else {
					buildFail++
				}
				*status = statusString(buildPass, buildFail, totalBuild, testPass, testFail, testWarning, testUnknown, totalTest)
				break
			}
		case result := <-testDone:
			{
				switch result.status {
				case JVSTestPass:
					testPass++
					break
				case JVSTestFail:
					testFail++
					break
				case JVSTestWarning:
					testWarning++
					break
				case JVSTestUnknown:
					testUnknown++
					break
				}
				*status = statusString(buildPass, buildFail, totalBuild, testPass, testFail, testWarning, testUnknown, totalTest)
				break
			}
		case <-done:
			break LableFor
		}
	}
	*status = finishStatusString(buildPass, buildFail, totalBuild, testPass, testFail, testWarning, testUnknown, totalTest)
}
