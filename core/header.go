package core

import "github.com/shady831213/jarvism/utils"

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
		return utils.Green("Pass")
	case JVSTestWarning:
		return utils.Yellow("Warning")
	case JVSTestFail:
		return utils.Red("Fail")
	case JVSTestUnknown:
		return utils.LightRed("Unknown")
	}
	return utils.LightRed("Unknown")
}

func (e *JVSTestResult) Error() string {
	return e.StatusString() + e.msg
}

func JVSTestResultPass(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestPass
	inst.msg = "!\n" + msg
	return inst
}

func JVSTestResultFail(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestFail
	inst.msg = "!\n" + "Error:" + msg
	return inst
}

func JVSTestResultWarning(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestWarning
	inst.msg = "!\n" + "Warning:" + msg
	return inst
}

func JVSTestResultUnknown(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.status = JVSTestUnknown
	inst.msg = "!\n" + "UnKnown:" + msg
	return inst
}
