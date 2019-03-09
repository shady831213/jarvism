package core

import (
	"github.com/shady831213/jarvism/utils"
	"strconv"
)

type JVSTestStatus int

const (
	_ JVSTestStatus = iota
	JVSTestPass
	JVSTestWarning
	JVSTestFail
	JVSTestUnknown
)

type JVSTestResult struct {
	status JVSTestStatus
	msg    string
}

func (e *JVSTestResult) StatusString() string {
	switch e.status {
	case JVSTestPass:
		return utils.Green("PASS")
	case JVSTestWarning:
		return utils.Yellow("WARNING")
	case JVSTestFail:
		return utils.Red("FAIL")
	case JVSTestUnknown:
		return utils.LightRed("UNKNOWN")
	}
	return utils.LightRed("UNKNOWN")
}

func (e *JVSTestResult) Error() string {
	return e.StatusString() + e.msg
}

func JVSTestResultPass(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestPass
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + msg
	}
	return inst
}

func JVSTestResultFail(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestFail
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + "Error:" + msg
	}
	return inst
}

func JVSTestResultWarning(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestWarning
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + "Warning:" + msg
	}
	return inst
}

func JVSTestResultUnknown(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestUnknown
	inst.msg = "!"
	if msg != "" {
		inst.msg += "\n" + "UnKnown:" + msg
	}
	return inst
}

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
