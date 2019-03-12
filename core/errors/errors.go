package errors

import "github.com/shady831213/jarvism/utils"

type JVSTestStatus int

const (
	_ JVSTestStatus = iota
	JVSTestPass
	JVSTestWarning
	JVSTestFail
	JVSTestUnknown
)

func StatusColor(status JVSTestStatus) func(str string, modifier ...interface{}) string {
	switch status {
	case JVSTestPass:
		return utils.Green
	case JVSTestWarning:
		return utils.Yellow
	case JVSTestFail:
		return utils.Red
	case JVSTestUnknown:
		return utils.LightRed
	}
	return utils.LightRed
}

func StatusString(status JVSTestStatus) string {
	switch status {
	case JVSTestPass:
		return "PASS"
	case JVSTestWarning:
		return "WARNING"
	case JVSTestFail:
		return "FAIL"
	case JVSTestUnknown:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}

func StatusShortString(status JVSTestStatus) string {
	switch status {
	case JVSTestPass:
		return "P"
	case JVSTestWarning:
		return "W"
	case JVSTestFail:
		return "F"
	case JVSTestUnknown:
		return "U"
	}
	return "U"
}

type JVSTestResult struct {
	Status JVSTestStatus
	Msg    string
}

func (e *JVSTestResult) Error() string {
	return StatusColor(e.Status)(StatusString(e.Status) + e.Msg)
}

func JVSTestResultPass(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.Status = JVSTestPass
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + msg
	}
	return inst
}

func JVSTestResultFail(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.Status = JVSTestFail
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + "Error:" + msg
	}
	return inst
}

func JVSTestResultWarning(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.Status = JVSTestWarning
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + "Warning:" + msg
	}
	return inst
}

func JVSTestResultUnknown(msg string) *JVSTestResult {
	inst := new(JVSTestResult)
	inst.Status = JVSTestUnknown
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + "UnKnown:" + msg
	}
	return inst
}

type AstError struct {
	Msg   string
	Item  string
	phase string
}

func (e *AstError) Error() string {
	return utils.Red(e.phase + " Error in " + e.Item + ": " + e.Msg)
}

func NewAstParseError(item, msg string) *AstError {
	inst := new(AstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Parse"
	return inst
}

func NewAstLinkError(item, msg string) *AstError {
	inst := new(AstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Link"
	return inst
}
