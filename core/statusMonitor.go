package core

import (
	"github.com/shady831213/jarvisSim/utils"
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
	return utils.Brown("[b:(") + utils.Green("p:"+strconv.Itoa(buildPass)) + utils.Brown("/") + utils.Red("f:"+strconv.Itoa(buildFail)) + utils.Brown("/") +
		utils.Brown("d:"+strconv.Itoa(buildPass+buildFail)+"/t"+strconv.Itoa(totalBuild)+")][t:(") + utils.Green("p:"+strconv.Itoa(testPass)) + utils.Brown("/") +
		utils.Red("f:"+strconv.Itoa(testFail)) + utils.Brown("/") + utils.Yellow("w:"+strconv.Itoa(testWarning)) + utils.Brown("/") + utils.LightRed("u:"+strconv.Itoa(testUnknown)) +
		utils.Brown("/") + utils.Brown("d:"+strconv.Itoa(testPass+testFail+testWarning+testUnknown)+"/t:"+strconv.Itoa(totalTest)+")]")
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
			*status = statusString(buildPass, buildFail, totalBuild, testPass, testFail, testWarning, testUnknown, totalTest)
			break LableFor
		default:
			*status = statusString(buildPass, buildFail, totalBuild, testPass, testFail, testWarning, testUnknown, totalTest)
		}
	}
}
