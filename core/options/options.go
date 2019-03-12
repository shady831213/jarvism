package options

import (
	"errors"
	"flag"
	"fmt"
)

func ArgToOption(s string) (string, error) {
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


//global
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)

//------------------------

func GetJvsOptions() *flag.FlagSet {
	return jvsOptions
}
