package parser

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type JvsOption interface {
	flag.Value
	Clone() JvsOption
}

func RegisterJvsOption(v JvsOption, name string, usage string) {
	jvsOptions.Var(v, name, usage)
}

func argToOption(s string) (string, error) {
	if len(s) < 2 || s[0] != '-' {
		return "", errors.New(fmt.Sprintf("bad flag syntax: %s", s))
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			return "", nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return name, errors.New(fmt.Sprintf("bad flag syntax: %s", s))
	}
	return name, nil
}

func GetOption(arg string) (JvsOption, error) {
	args := strings.Split(arg, " ")
	if err := jvsOptions.Parse(args); err != nil {
		return nil, err
	}
	optName, err := argToOption(args[0])
	if err != nil {
		return nil, err
	}
	v, ok := jvsOptions.Lookup(optName).Value.(JvsOption)
	if !ok {
		panic(fmt.Sprintf("expect type JvsOption but get %T", jvsOptions.Lookup(optName).Value))
	}
	return v, nil
}

type jvsNonBoolOption struct {
}

func (t *RepeatOption) IsBoolFlag() bool {
	return false
}

//buildin options
type RepeatOption struct {
	jvsNonBoolOption
	n int
}

func newRepeatOption() *RepeatOption {
	inst := new(RepeatOption)
	inst.n = 1
	return inst
}

func (t *RepeatOption) Clone() JvsOption {
	inst := newRepeatOption()
	inst.n = t.n
	return inst
}

func (t *RepeatOption) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	t.n = n
	return nil
}

func (t *RepeatOption) String() string {
	return string(t.n)
}

func init() {
	RegisterJvsOption(newRepeatOption(), "repeat", "run each testcase repeatly n times")
}

//global
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)
//------------------------
