/*
JVSRuntime* for runtime results

JVSAst* for Lex, ast parser and Plugin loader
*/
package errors

import (
	"github.com/shady831213/jarvism/core/utils"
	"strings"
)

type JVSRuntimeStatus int

const (
	_ JVSRuntimeStatus = iota
	JVSRuntimePass
	JVSRuntimeUnknown
	JVSRuntimeWarning
	JVSRuntimeFail
)

//render status:
//
//pass green
//
//warning yellow
//
//fail red
//
//unknown light red
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

//convert status to string
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

//convert status to short string
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

//runtime result, for build and test
//
//Status: pass, fail, unknown, warning
//
//title:  "", Error, Unknown, Warning
//
//msg: messages
//
//Name: build name or test name
type JVSRuntimeResult struct {
	Status JVSRuntimeStatus
	title  string
	msg    []string
	Name   string
}

func (e *JVSRuntimeResult) Error() string {
	msg := e.GetMsg()
	if msg != "" {
		msg = "\n" + e.title + msg
	}
	return StatusColor(e.Status)(StatusString(e.Status) + msg)
}

func (e *JVSRuntimeResult) GetMsg() string {
	return strings.Join(e.msg, "\n")
}

func (e *JVSRuntimeResult) addMsgs(msgs ...string) {
	for _, msg := range msgs {
		if strings.Replace(strings.Replace(msg, " ", "", -1), "\n", "", -1) != "" {
			e.msg = append(e.msg, strings.TrimRight(msg, "\n"))
		}
	}
}
//create runtime result
func NewJVSRuntimeResult(status JVSRuntimeStatus, msgs ...string) *JVSRuntimeResult {
	switch status {
	case JVSRuntimePass:
		return JVSRuntimeResultPass(msgs...)
	case JVSRuntimeWarning:
		return JVSRuntimeResultWarning(msgs...)
	case JVSRuntimeFail:
		return JVSRuntimeResultFail(msgs...)
	case JVSRuntimeUnknown:
		return JVSRuntimeResultUnknown(msgs...)
	}
	return JVSRuntimeResultUnknown(msgs...)
}

//create pass runtime result
func JVSRuntimeResultPass(msgs ...string) *JVSRuntimeResult {
	inst := &JVSRuntimeResult{
		JVSRuntimePass,
		"",
		make([]string, 0),
		"",
	}
	inst.addMsgs(msgs...)
	return inst
}

//create fail runtime result
func JVSRuntimeResultFail(msgs ...string) *JVSRuntimeResult {
	inst := &JVSRuntimeResult{
		JVSRuntimeFail,
		"Error:",
		make([]string, 0),
		"",
	}
	inst.addMsgs(msgs...)
	return inst
}

//create waring runtime result
func JVSRuntimeResultWarning(msgs ...string) *JVSRuntimeResult {
	inst := &JVSRuntimeResult{
		JVSRuntimeWarning,
		"Warning:",
		make([]string, 0),
		"",
	}
	inst.addMsgs(msgs...)
	return inst
}

//create unknown runtime result
func JVSRuntimeResultUnknown(msgs ...string) *JVSRuntimeResult {
	inst := &JVSRuntimeResult{
		JVSRuntimeUnknown,
		"Unknown:",
		make([]string, 0),
		"",
	}
	inst.addMsgs(msgs...)
	return inst
}

//for lexer, parser and plugin loader
//Msg: messages
//Item: file, plugin or ast item
//phase: lex, parse, link or pluginLoad
type JVSAstError struct {
	Msg   string
	Item  string
	phase string
}

func (e *JVSAstError) Error() string {
	return utils.Red(e.phase + " Error in " + e.Item + ": " + e.Msg)
}

//create parse error
func JVSAstParseError(item, msg string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Parse"
	return inst
}

//create link error
func JVSAstLinkError(item, msg string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Link"
	return inst
}

//create lex error
func JVSAstLexError(item, msg string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = item
	inst.phase = "Lex"
	return inst
}

//create plugin load error
func JVSPluginLoadError(pluginName, msg, filePath string) *JVSAstError {
	inst := new(JVSAstError)
	inst.Msg = msg
	inst.Item = pluginName + "(" + filePath + ")"
	inst.phase = "Load Plugin"
	return inst
}

//stderr writer, collect stderr output
type JVSStderr struct {
	Msg string
}

func (e *JVSStderr) Write(p []byte) (n int, err error) {
	e.Msg += string(p)
	return len(p), nil
}
