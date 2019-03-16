package loader

import (
	"errors"
	"flag"
	"fmt"
	"github.com/shady831213/jarvism/core"
	jvsErrors "github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/options"
	"math"
	"math/rand"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type JvsAstOption interface {
	flag.Value
	GetName() string
	Clone() JvsAstOption
	Usage() string
}

type JvsAstOptionForTest interface {
	JvsAstOption
	TestHandler(test *AstTestCase)
}

type JvsAstOptionForBuild interface {
	JvsAstOption
	BuildHandler(build *AstBuild)
}

func RegisterJvsAstOption(v JvsAstOption) {
	options.GetJvsOptions().Var(v, v.GetName(), v.Usage())
}

func GetJvsAstOption(arg string) (JvsAstOption, error) {
	args := strings.Split(regexp.MustCompile(`^\s+`).ReplaceAllString(arg, ""), " ")
	optName, err := options.ArgToOption(args[0])
	if err != nil {
		return nil, err
	}
	if len(args) > 1 {
		args[0] += "="
	}
	if err := options.GetJvsOptions().Parse([]string{strings.Join(args, "")}); err != nil {
		return nil, err
	}
	v, ok := options.GetJvsOptions().Lookup(optName).Value.(JvsAstOption)
	if !ok {
		return nil, errors.New("Not JvsAstOption")
	}
	return v, nil
}

func LoadBuildInOptions(configFile string) error {
	if configFile == "" {
		return nil
	}
	cfg, err := Lex(configFile)
	if err != nil {
		panic(err)
	}
	if err := CfgToAstItemRequired(cfg, "options", WithCheckMap(func(item map[interface{}]interface{}) *jvsErrors.JVSAstError {
		for name, optionCfg := range item {
			if option, ok := jvsAstRoot.Options[name.(string)]; ok {
				return jvsErrors.JVSAstParseError("options", "option conflict["+option.Name+","+name.(string)+"]! option name must be unique!")
			}
			jvsAstRoot.Options[name.(string)] = newAstOption(name.(string))
			if err := AstParse(jvsAstRoot.Options[name.(string)], optionCfg); err != nil {
				if err.Item == "" {
					err.Item = name.(string)
				}
				return err
			}
		}
		return nil
	})); err != nil {
		return err
	}
	return nil
}

type jvsAstNonBoolOption struct {
}

func (t *jvsAstNonBoolOption) IsBoolFlag() bool {
	return false
}

//buildin options

type RepeatOption struct {
	jvsAstNonBoolOption
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

func (t *SeedOption) Usage() string {
	return "run testcase with specific seed"
}

//------------------------

func (t *RepeatOption) Usage() string {
	return "run each testcase repeatly n times"
}

type SeedOption struct {
	jvsAstNonBoolOption
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

var jvsRand *rand.Rand

func init() {
	jvsRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := LoadBuildInOptions(path.Join(core.BuildInOptionPath(), "global_options.yaml")); err != nil {
		panic("Error in loading " + path.Join(core.BuildInOptionPath(), "global_options.yaml") + ":" + err.Error())
	}
	RegisterJvsAstOption(newRepeatOption())
	RegisterJvsAstOption(newSeedOption())
}
