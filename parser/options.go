package parser

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type JvsOption interface {
	flag.Value
	GetName() string
	Usage() string
	AfterParse()
}

func RegisterOption(option JvsOption) {
	jvsOptions.Var(option, option.GetName(), option.Usage())
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

//global
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)
//------------------------
