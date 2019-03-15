package ast

import (
	"bufio"
	"fmt"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/utils"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Checker interface {
	astParser
	Name() string
	Check() *errors.JVSRuntimeResult
	Input(reader io.Reader)
}

type CheckerBase struct {
	name      string
	pats      map[errors.JVSRuntimeStatus][]*regexp.Regexp
	exclPats  map[errors.JVSRuntimeStatus][]*regexp.Regexp
	finishPat *regexp.Regexp
	result    *errors.JVSRuntimeResult
	finished  bool
	input     *bufio.Reader
}

func (c *CheckerBase) Name() string {
	return c.name
}

func (c *CheckerBase) Init(name string) {
	c.name = name
	c.result = errors.JVSRuntimeResultUnknown("Unfinished!")
	c.finished = true
	c.pats = make(map[errors.JVSRuntimeStatus][]*regexp.Regexp)
	c.exclPats = make(map[errors.JVSRuntimeStatus][]*regexp.Regexp)
	c.pats[errors.JVSRuntimeFail] = make([]*regexp.Regexp, 0)
	c.pats[errors.JVSRuntimeWarning] = make([]*regexp.Regexp, 0)
	c.pats[errors.JVSRuntimeUnknown] = make([]*regexp.Regexp, 0)
	c.exclPats[errors.JVSRuntimeFail] = make([]*regexp.Regexp, 0)
	c.exclPats[errors.JVSRuntimeWarning] = make([]*regexp.Regexp, 0)
	c.exclPats[errors.JVSRuntimeUnknown] = make([]*regexp.Regexp, 0)
}

func (c *CheckerBase) AddPats(status errors.JVSRuntimeStatus, exl bool, pat *regexp.Regexp, pats ...*regexp.Regexp) {
	if exl {
		c.exclPats[status] = append(c.exclPats[status], pat)
		if pats != nil && len(pats) > 0 {
			c.exclPats[status] = append(c.exclPats[status], pats...)
		}
	}
	c.pats[status] = append(c.pats[status], pat)
	if pats != nil && len(pats) > 0 {
		c.pats[status] = append(c.pats[status], pats...)
	}
}

func (c *CheckerBase) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	//AstParse pats and exclPats
	for k := range c.pats {
		name := strings.ToLower(errors.StatusString(k))
		exclName := "exclude_" + name
		if err := CfgToAstItemOptional(cfg, name, func(item interface{}) *errors.JVSAstError {
			pats, ok := item.([]interface{})
			if !ok {
				return errors.JVSAstParseError(name, fmt.Sprintf("must be a list of strings but get %T!", item))
			}
			for _, pat := range pats {
				re, err := regexp.Compile(pat.(string))
				if err != nil {
					return errors.JVSAstParseError(name, err.Error())
				}
				c.pats[k] = append(c.pats[k], re)
			}
			return nil
		}); err != nil {
			return errors.JVSAstParseError(name+" of "+c.Name(), err.Error())
		}

		if err := CfgToAstItemOptional(cfg, exclName, func(item interface{}) *errors.JVSAstError {
			pats, ok := item.([]interface{})
			if !ok {
				return errors.JVSAstParseError(exclName, fmt.Sprintf("must be a list of strings but get %T!", item))
			}
			for _, pat := range pats {
				re, err := regexp.Compile(pat.(string))
				if err != nil {
					return errors.JVSAstParseError(exclName, err.Error())
				}
				c.exclPats[k] = append(c.exclPats[k], re)
			}
			return nil
		}); err != nil {
			return errors.JVSAstParseError(exclName+" of "+c.Name(), err.Error())
		}
	}
	if err := CfgToAstItemOptional(cfg, "finish_flag", func(item interface{}) *errors.JVSAstError {
		pat, ok := item.(string)
		if !ok {
			return errors.JVSAstParseError("finish_flag", fmt.Sprintf("must be a string but get %T!", item))
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			return errors.JVSAstParseError("finish_flag", err.Error())
		}
		c.finishPat = re
		c.finished = false
		return nil
	}); err != nil {
		return errors.JVSAstParseError("finish_flag of "+c.Name(), err.Error())
	}
	return nil
}

func (c *CheckerBase) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	for k, _ := range c.pats {
		name := strings.ToLower(errors.StatusString(k))
		exclName := "exclude_" + name
		keywords.AddKey(name)
		keywords.AddKey(exclName)
	}
	keywords.AddKey("finish_flag")
	if !CheckKeyWord(s, keywords) {
		return false, keywords, "Error in " + c.Name() + ":"
	}
	return true, nil, ""
}

func (c *CheckerBase) Input(reader io.Reader) {
	c.input = bufio.NewReader(reader)
}

func (c *CheckerBase) Check() (res *errors.JVSRuntimeResult) {
	if c.input == nil {
		panic("nil input! Call Input(reader) first!")
	}
	line := 1
	for {
		s, e := c.input.ReadString('\n')
		if e != nil && e != io.EOF {
			return errors.JVSRuntimeResultUnknown(e.Error())
		}
		if c.result.Status < errors.JVSRuntimeFail {
			c.checkLine(s, line)
		}
		if e == io.EOF {
			return c.status()
		}
		line++
	}
}

func (c *CheckerBase) checkLine(s string, line int) {
	//check error first
	for status := errors.JVSRuntimeFail; status > errors.JVSRuntimePass; status-- {
		if c.checkStatus(status, s, line) {
			break
		}
	}
	//check finish flag
	if c.finishPat != nil && c.finishPat.MatchString(s) {
		c.finished = true
	}
}

func (c *CheckerBase) status() *errors.JVSRuntimeResult {
	if c.result.Status > errors.JVSRuntimeUnknown {
		return c.result
	}
	if c.finished {
		return errors.JVSRuntimeResultPass("")
	}
	return c.result
}

func (c *CheckerBase) checkStatus(status errors.JVSRuntimeStatus, s string, line int) bool {
	for _, p := range c.pats[status] {
		if p.MatchString(s) {
			for _, exp := range c.exclPats[status] {
				//excluded
				if exp.MatchString(s) {
					return false
				}
			}
			//matched
			if status > c.result.Status {
				c.result = errors.NewJVSRuntimeResult(status, strings.Replace(s, "\n", "", -1)+", line "+strconv.Itoa(line))
			}
			return true
		}
	}
	return false
}

var validChecker = make(map[string]func() Checker)

func GetChecker(key string) Checker {
	if v, ok := validChecker[key]; ok {
		return v()
	}
	return nil
}

func RegisterChecker(d func() Checker) {
	inst := d()
	if _, ok := validChecker[inst.Name()]; ok {
		panic("Checker " + inst.Name() + " has been registered!")
	}
	validChecker[inst.Name()] = d
}
