package core

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

type JvsAstOption interface {
	flag.Value
	GetName() string
	Clone() JvsAstOption
}

func RegisterJvsAstOption(v JvsAstOption, usage string) {
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

func getJvsAstOption(arg string) (JvsAstOption, error) {
	args := strings.Split(arg, " ")
	if err := jvsOptions.Parse(args); err != nil {
		return nil, err
	}
	optName, err := argToOption(args[0])
	if err != nil {
		return nil, err
	}
	v, ok := jvsOptions.Lookup(optName).Value.(JvsAstOption)
	if !ok {
		return nil, errors.New("Not JvsAstOption")
	}
	return v, nil
}

type jvsAstNonBoolOption struct {
}

func (t *jvsAstNonBoolOption) IsBoolFlag() bool {
	return false
}

//buildin options
var runTimeMaxJob int

type jvsAstTestOption struct {
	jvsAstNonBoolOption
}

func (t *jvsAstTestOption) IsCompileOption() bool {
	return false
}

type RepeatOption struct {
	jvsAstTestOption
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

func (t *RepeatOption) Clone() JvsAstOption {
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
	//no one touch it
	if test.seeds == nil {
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
}

type SeedOption struct {
	jvsAstTestOption
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

func (t *SeedOption) Clone() JvsAstOption {
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
	test.seeds = make([]int, 1)
	if t.n == 0 {
		test.seeds[0] = jvsRand.Intn(math.MaxInt32)
		return
	}
	test.seeds[0] = t.n
}

func init() {
	jvsRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	RegisterJvsAstOption(newRepeatOption(), "run each testcase repeatly n times")
	RegisterJvsAstOption(newSeedOption(), "run testcase with specific seed")
	jvsOptions.IntVar(&runTimeMaxJob, "max_job", -1, "limit of runtime coroutines, default is unlimited.")
}

//global
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)
var jvsRand *rand.Rand
//------------------------
