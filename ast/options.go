package ast

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type JvsOption interface {
	flag.Value
	GetName() string
	Clone() JvsOption
}

func RegisterJvsOption(v JvsOption, usage string) {
	jvsOptions.Var(v, v.GetName(), usage)
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

func (t *jvsNonBoolOption) IsBoolFlag() bool {
	return false
}

//buildin options
type jvsTestOption struct {
	jvsNonBoolOption
}

func (t *jvsTestOption) IsCompileOption() bool {
	return false
}

type RepeatOption struct {
	jvsTestOption
	n int
}

func newRepeatOption() *RepeatOption {
	inst := new(RepeatOption)
	inst.n = 1
	return inst
}

func (t *RepeatOption) GetName() string {
	return "repeat"
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

func (t *RepeatOption) TestHandler(test *AstTestCase) {
	test.seeds = make([]int, 0)
	seedsMap := make(map[int]interface{})
	for len(seedsMap) < t.n {
		seed := jvsRand.Intn(math.MaxInt32)
		if _, ok := seedsMap[seed]; !ok {
			seedsMap[seed] = nil
			test.seeds = append(test.seeds, seed)
		}
	}
	if len(test.seeds) != t.n {
		panic(fmt.Sprintf("len of seeds %d != t.n %d !", len(test.seeds), t.n))
	}
}

type SeedOption struct {
	jvsTestOption
	n int
}

func newSeedOption() *SeedOption {
	inst := new(SeedOption)
	inst.n = 0
	return inst
}

func (t *SeedOption) GetName() string {
	return "seed"
}

func (t *SeedOption) Clone() JvsOption {
	inst := newSeedOption()
	inst.n = t.n
	return inst
}

func (t *SeedOption) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	t.n = n
	return nil
}

func (t *SeedOption) String() string {
	return string(t.n)
}

func (t *SeedOption) TestHandler(test *AstTestCase) {
	//no one touch it
	if test.seeds == nil {
		test.seeds = make([]int, 1)
		if t.n == 0 {
			test.seeds[0] = jvsRand.Intn(math.MaxInt32)
			return
		}
		test.seeds[0] = t.n
	}
}

func SetRand(rand *rand.Rand) {
	jvsRand = rand
}

func init() {
	jvsRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	RegisterJvsOption(newRepeatOption(), "run each testcase repeatly n times")
	RegisterJvsOption(newSeedOption(), "run testcase with specific seed")
}

//global
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)
var jvsRand *rand.Rand
//------------------------
