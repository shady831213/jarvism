package errors

import "github.com/shady831213/jarvism/utils"

type JVSRuntimeStatus int

const (
	_ JVSRuntimeStatus = iota
	JVSRuntimePass
	JVSRuntimeWarning
	JVSRuntimeFail
	JVSRuntimeUnknown
)

func StatusColor(status JVSRuntimeStatus) func(str string, modifier ...interface{}) string {
	switch status {
	case JVSRuntimePass:
		return utils.Green
	case JVSRuntimeWarning:
		return utils.Yellow
	case JVSRuntimeFail:
		return utils.Red
	case JVSRuntimeUnknown:
		return utils.LightRed
	}
	return utils.LightRed
}

func StatusString(status JVSRuntimeStatus) string {
	switch status {
	case JVSRuntimePass:
		return "PASS"
	case JVSRuntimeWarning:
		return "WARNING"
	case JVSRuntimeFail:
		return "FAIL"
	case JVSRuntimeUnknown:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}

func StatusShortString(status JVSRuntimeStatus) string {
	switch status {
	case JVSRuntimePass:
		return "P"
	case JVSRuntimeWarning:
		return "W"
	case JVSRuntimeFail:
		return "F"
	case JVSRuntimeUnknown:
		return "U"
	}
	return "U"
}

type JVSRuntimeResult struct {
	Status JVSRuntimeStatus
	Msg    string
}

func (e *JVSRuntimeResult) Error() string {
	return StatusColor(e.Status)(StatusString(e.Status) + e.Msg)
}

func JVSRuntimeResultPass(msg string) *JVSRuntimeResult {
	inst := new(JVSRuntimeResult)
	inst.Status = JVSRuntimePass
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + msg
	}
	return inst
}

func JVSRuntimeResultFail(msg string) *JVSRuntimeResult {
	inst := new(JVSRuntimeResult)
	inst.Status = JVSRuntimeFail
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + "Error:" + msg
	}
	return inst
}

func JVSRuntimeResultWarning(msg string) *JVSRuntimeResult {
	inst := new(JVSRuntimeResult)
	inst.Status = JVSRuntimeWarning
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + "Warning:" + msg
	}
	return inst
}

func JVSRuntimeResultUnknown(msg string) *JVSRuntimeResult {
	inst := new(JVSRuntimeResult)
	inst.Status = JVSRuntimeUnknown
	inst.Msg = "!"
	if msg != "" {
		inst.Msg += "\n" + "UnKnown:" + msg
	}
	return inst
}

type JVSAstError struct {
	Msg   string
	Item  string
	phase string
}

func (e *JVSAstError) Error() string {
	return utils.Red(e.phase + " Error in " + e.Item + ": " + e.Msg)
}

func JVSAstParseError(item, msg string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Parse"
	return inst
}

func JVSAstLinkError(item, msg string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Link"
	return inst
}

func JVSAstLexError(item, msg string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Lex"
	return inst
}
