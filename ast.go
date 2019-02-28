package jarivsSim

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type astParser interface {
	//pass1:top-down astParse
	Parse(map[interface{}]interface{}) error
	KeywordsChecker(s string) (bool, []string, string)
	//print:for debug
	GetHierString(space int) string;
}

type astLinker interface {
	astParser
	//pass2:top-down link
	Link() error
}

func astError(item string, err error) error {
	return errors.New("Error in " + item + ": " + err.Error())
}

func checkKeyWord(s string, keyWords map[string]interface{}) bool {
	_, ok := keyWords[s]
	return ok
}

func astHierFmt(title string, space int, handler func() string) string {
	return fmt.Sprintln(strings.Repeat(" ", space)+strings.Repeat("-", 20-space)) +
		fmt.Sprintln(strings.Repeat(" ", space)+title) +
		handler() +
		"\n"
}

func astParse(parser astParser, cfg map[interface{}]interface{}) error {
	for name, _ := range cfg {
		if ok, keywords, tag := parser.KeywordsChecker(name.(string)); !ok {
			return errors.New(tag + "syntax error of " + name.(string) + "! expect " + fmt.Sprint(keywords))
		}
	}
	if err := parser.Parse(cfg); err != nil {
		return err
	}
	return nil
}

func cfgToastItemRequired(cfg map[interface{}]interface{}, key string, handler func(interface{}) error) error {
	if item, ok := cfg[key]; ok {
		flag.Args()
		return handler(item)
	}
	return errors.New("not define " + key + "!")
}

func cfgToastItemOptional(cfg map[interface{}]interface{}, key string, handler func(interface{}) error) error {
	if item, ok := cfg[key]; ok {
		return handler(item)
	}
	return nil
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

type astOptionItem struct {
	content string
}

func newAstOptionItem(content interface{}) *astOptionItem {
	inst := new(astOptionItem)
	if value, ok := content.(string); ok {
		inst.content = value
		return inst
	}
	if value, ok := content.([]string); ok {
		inst.content = strings.Join(value, "")
		return inst
	}
	if value, ok := content.([]interface{}); ok {
		for _, i := range (value) {
			s, ok := i.(string)
			if !ok {
				panic(fmt.Sprintf("content must be string or []string, but it is %T !", content))
				return nil
			}
			inst.content += s
		}

		return inst
	}
	panic(fmt.Sprintf("content must be string or []string, but it is %T !", content))
	return nil
}

type astSimOnlyItem struct {
	PreSimOption  *astOptionItem
	SimOption     *astOptionItem
	PostSimOption *astOptionItem
}

func (t *astSimOnlyItem) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"pre_sim_option": nil, "sim_option": nil, "post_sim_option": nil}
	if !checkKeyWord(s, keywords) {
		return false, KeyOfStringMap(keywords), ""
	}
	return true, nil, ""
}

func (t *astSimOnlyItem) Parse(cfg map[interface{}]interface{}) error {
	if err := cfgToastItemOptional(cfg, "pre_sim_option", func(item interface{}) error {
		t.PreSimOption = newAstOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := cfgToastItemOptional(cfg, "sim_option", func(item interface{}) error {
		t.SimOption = newAstOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := cfgToastItemOptional(cfg, "post_sim_option", func(item interface{}) error {
		t.PostSimOption = newAstOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *astSimOnlyItem) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("PreSimOption:", nextSpace, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(t.PreSimOption)
	}) +
		astHierFmt("SimOption:", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(t.SimOption)
		}) +
		astHierFmt("PostSimOption:", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(t.PostSimOption)
		})
}

type astBuildItem struct {
	astSimOnlyItem
	PreCompileOption  *astOptionItem
	CompileOption     *astOptionItem
	PostCompileOption *astOptionItem
}

func (t *astBuildItem) KeywordsChecker(s string) (bool, []string, string) {
	if ok, simKeywords, _ := t.astSimOnlyItem.KeywordsChecker(s); !ok {
		compileKeywords := map[string]interface{}{"pre_compile_option": nil, "compile_option": nil, "post_compile_option": nil}
		if !checkKeyWord(s, compileKeywords) {
			return false, append(simKeywords, KeyOfStringMap(compileKeywords)...), ""
		}
	}

	return true, nil, ""
}

func (t *astBuildItem) Parse(cfg map[interface{}]interface{}) error {
	if err := t.astSimOnlyItem.Parse(cfg); err != nil {
		return err
	}
	if err := cfgToastItemOptional(cfg, "pre_compile_option", func(item interface{}) error {
		t.PreCompileOption = newAstOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := cfgToastItemOptional(cfg, "compile_option", func(item interface{}) error {
		t.CompileOption = newAstOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := cfgToastItemOptional(cfg, "post_compile_option", func(item interface{}) error {
		t.PostCompileOption = newAstOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *astBuildItem) GetHierString(space int) string {
	nextSpace := space + 1
	return t.astSimOnlyItem.GetHierString(space) +
		astHierFmt("PreCompileOption:", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(t.PreCompileOption)
		}) +
		astHierFmt("CompileOption:", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(t.CompileOption)
		}) +
		astHierFmt("PostCompileOption:", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(t.PostCompileOption)
		})
}

//Options
//------------------------
type astOptionAction struct {
	astBuildItem
}

type astOption struct {
	On        *astOptionAction
	Off       *astOptionAction
	WithValue *astOptionAction
	Value     string
	Name      string
}

func newAstOption(name string) *astOption {
	inst := new(astOption)
	inst.Name = name
	inst.Value = "false"
	return inst
}

func (t *astOption) Clone() *astOption {
	inst := newAstOption(t.Name)
	inst.Value = t.Value
	inst.On = t.On
	inst.Off = t.Off
	inst.WithValue = t.WithValue
	return inst
}

func (t *astOption) Set(s string) error {
	if t.WithValue != nil {
		t.Value = s
		return nil
	}
	b, err := strconv.ParseBool(s)
	t.Value = strconv.FormatBool(b)
	return err
}

func (t *astOption) String() string {
	return t.Value
}

func (t *astOption) IsBoolFlag() bool {
	return t.WithValue == nil
}

func (t *astOption) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"on_action": nil, "off_action": nil, "with_value_action": nil}
	if !checkKeyWord(s, keywords) {
		return false, KeyOfStringMap(keywords), "Error in " + t.Name + ":"
	}
	return true, nil, ""
}

func (t *astOption) Parse(cfg map[interface{}]interface{}) error {
	if err := cfgToastItemOptional(cfg, "on_action", func(item interface{}) error {
		t.On = new(astOptionAction)
		return astParse(t.On, item.(map[interface{}]interface{}))
	}); err != nil {
		return astError("on_action of "+t.Name, err)
	}
	if err := cfgToastItemOptional(cfg, "off_action", func(item interface{}) error {
		t.Off = new(astOptionAction)
		return astParse(t.Off, item.(map[interface{}]interface{}))
	}); err != nil {
		return astError("off_action of "+t.Name, err)
	}
	if err := cfgToastItemOptional(cfg, "with_value_action", func(item interface{}) error {
		t.WithValue = new(astOptionAction)
		return astParse(t.WithValue, item.(map[interface{}]interface{}))
	}); err != nil {
		return astError("with_value_action of "+t.Name, err)
	}
	//add to flagSet
	jvsOptions.Var(t, t.Name, "user-defined flag")
	return nil
}

func (t *astOption) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return astHierFmt("Value:", nextSpace, func() string {
			return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.Value)
		}) +
			astHierFmt("On:", nextSpace, func() string {
				if t.On != nil {
					return t.On.GetHierString(nextSpace + 1)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace+1) + "null")
			}) +
			astHierFmt("Off:", nextSpace, func() string {
				if t.Off != nil {
					return t.Off.GetHierString(nextSpace + 1)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace+1) + "null")
			}) +
			astHierFmt("WithValue:", nextSpace, func() string {
				if t.WithValue != nil {
					return t.WithValue.GetHierString(nextSpace + 1)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace+1) + "null")
			})
	})
}

//------------------------

//env
//------------------------
type astEnv struct {
	Simulator string
}

func (t *astEnv) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"simulator": nil}
	if !checkKeyWord(s, keywords) {
		return false, KeyOfStringMap(keywords), "Error in Env:"
	}
	return true, nil, ""
}

func (t *astEnv) Parse(cfg map[interface{}]interface{}) error {
	if err := cfgToastItemRequired(cfg, "simulator", func(item interface{}) error {
		t.Simulator = item.(string)
		return nil
	}); err != nil {
		return astError("Env", err)
	}
	return nil
}

func (t *astEnv) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("Simulator:", nextSpace, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(t.Simulator)
	})
}

//------------------------

//Build
//------------------------
type astBuild struct {
	astBuildItem
	Name string
}

func newAstBuild(name string) *astBuild {
	inst := new(astBuild)
	inst.Name = name
	return inst
}

func (t *astBuild) KeywordsChecker(s string) (bool, []string, string) {
	if ok, buildKeywords, _ := t.astBuildItem.KeywordsChecker(s); !ok {
		return false, buildKeywords, "Error in build " + t.Name + ":"
	}

	return true, nil, ""
}

func (t *astBuild) Parse(cfg map[interface{}]interface{}) error {
	return t.astBuildItem.Parse(cfg)
}

func (t *astBuild) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.astBuildItem.GetHierString(nextSpace)
	})
}

//------------------------

//Test and Group, linkable
//------------------------
type astTestOpts interface {
	SetParent(parent astTestOpts)
	GetName() string
	//bottom-up search
	GetOptionArgs() map[string]*astOption
	//bottom-up search
	GetBuild() *astBuild
}

type astTest struct {
	Name       string
	OptionArgs map[string]*astOption
	args       []string
	parent     astTestOpts
}

func (t *astTest) init(name string) {
	t.Name = name
	t.OptionArgs = make(map[string]*astOption)
}

func (t *astTest) SetParent(parent astTestOpts) {
	t.parent = parent
}

func (t *astTest) GetOptionArgs() map[string]*astOption {
	if t.parent != nil {
		options := make(map[string]*astOption)
		for k, v := range t.parent.GetOptionArgs() {
			options[k] = v
		}
		for k, v := range t.OptionArgs {
			options[k] = v
		}
		return options
	}
	return t.OptionArgs
}

func (t *astTest) GetName() string {
	return t.Name
}

func (t *astTest) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"args": nil}
	if !checkKeyWord(s, keywords) {
		return false, KeyOfStringMap(keywords), "Error in " + t.Name + ":"
	}
	return true, nil, ""
}

func (t *astTest) Parse(cfg map[interface{}]interface{}) error {
	if err := cfgToastItemOptional(cfg, "args", func(item interface{}) error {
		t.args = make([]string, 0)
		for _, arg := range (item.([]interface{})) {
			t.args = append(t.args, arg.(string))
		}
		return nil
	}); err != nil {
		return astError(t.Name, err)
	}
	return nil
}

//because Link is top-down, the last repeated args take effect
func (t *astTest) Link() error {
	for _, arg := range t.args {
		//Options have been all parsed
		args := strings.Split(arg, " ")
		if err := jvsOptions.Parse(args); err != nil {
			return astError("args of "+t.Name, err)
		}
		optName, err := argToOption(args[0])
		if err != nil {
			return astError("args of "+t.Name, err)
		}
		t.OptionArgs[optName] = jvsOptions.Lookup(optName).Value.(*astOption).Clone()
	}
	return nil
}

func (t *astTest) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("parent:", nextSpace, func() string {
		if t.parent != nil {
			return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.parent.GetName())
		}
		return fmt.Sprintln(strings.Repeat(" ", nextSpace) + "null")
	}) +
		astHierFmt("OptionArgs:", nextSpace, func() string {
			s := ""
			keys := make([]string, 0)
			args := t.GetOptionArgs()
			for k := range args {
				keys = append(keys, k)
			}
			ForeachStringKeysInOrder(keys, func(i string) {
				s += args[i].GetHierString(nextSpace + 1)
			})
			return s
		})
}

type astTestCase struct {
	astTest
}

func newAstTestCase(name string) *astTestCase {
	inst := new(astTestCase)
	inst.init(name)
	return inst
}

func (t *astTestCase) GetBuild() *astBuild {
	return t.parent.GetBuild()
}

func (t *astTestCase) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.astTest.GetHierString(nextSpace) +
			astHierFmt("Builds:", nextSpace, func() string {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.GetBuild().Name)
			})
	})
}

type astGroup struct {
	astTest
	buildName string
	Build     *astBuild
	Tests     map[string]*astTestCase
	Groups    map[string]*astGroup
}

func newAstGroup(name string) *astGroup {
	inst := new(astGroup)
	inst.init(name)
	return inst
}

func (t *astGroup) GetBuild() *astBuild {
	if t.Build != nil {
		return t.Build
	}
	if t.parent != nil {
		return t.parent.GetBuild()
	}
	return nil
}

func (t *astGroup) KeywordsChecker(s string) (bool, []string, string) {
	if ok, testKeywords, _ := t.astTest.KeywordsChecker(s); !ok {
		groupKeywords := map[string]interface{}{"build": nil, "tests": nil, "groups": nil}
		if !checkKeyWord(s, groupKeywords) {
			return false, append(testKeywords, KeyOfStringMap(groupKeywords)...), "Error in group " + t.Name + ":"
		}
	}
	return true, nil, ""
}

func (t *astGroup) Parse(cfg map[interface{}]interface{}) error {
	if err := t.astTest.Parse(cfg); err != nil {
		return err
	}
	if err := cfgToastItemOptional(cfg, "build", func(item interface{}) error {
		t.buildName = item.(string)
		return nil
	}); err != nil {
		return astError("group "+t.Name, err)
	}
	//astParse tests
	if err := cfgToastItemOptional(cfg, "tests", func(item interface{}) error {
		t.Tests = make(map[string]*astTestCase)
		for name, test := range item.(map[interface{}]interface{}) {
			t.Tests[name.(string)] = newAstTestCase(name.(string))
			if err := astParse(t.Tests[name.(string)], test.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return astError("group "+t.Name, err)
	}
	//astParse groups
	if err := cfgToastItemOptional(cfg, "groups", func(item interface{}) error {
		t.Groups = make(map[string]*astGroup)
		for _, name := range item.([]interface{}) {
			if _, ok := t.Groups[name.(string)]; ok {
				return errors.New("sub group " + name.(string) + " is redefined in group " + t.Name + "!")
			}
			t.Groups[name.(string)] = nil
		}
		return nil
	}); err != nil {
		return astError("group "+t.Name, err)
	}
	return nil
}

func (t *astGroup) Link() error {
	//link build
	//builds have been all parsed
	build := jvsAstRoot.GetBuild(t.buildName)
	if build == nil {
		return errors.New("build " + t.buildName + "of group " + t.Name + "is undef!")
	}
	t.Build = build
	//link args
	if err := t.astTest.Link(); err != nil {
		return err
	}
	//link tests
	for _, test := range t.Tests {
		if err := test.Link(); err != nil {
			return err
		}
		test.SetParent(t)
	}
	//link groups
	for name, _ := range t.Groups {
		group := jvsAstRoot.GetGroup(name)
		if group == nil {
			return errors.New("sub group " + name + "of group " + t.Name + " is undef!")
		}
		t.Groups[name] = group
		group.SetParent(t)
	}
	return nil
}

func (t *astGroup) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.astTest.GetHierString(nextSpace) +
			astHierFmt("Builds:", nextSpace, func() string {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.GetBuild().Name)
			}) +
			astHierFmt("Tests:", nextSpace, func() string {
				s := ""
				keys := make([]string, 0)
				for k := range t.Tests {
					keys = append(keys, k)
				}
				ForeachStringKeysInOrder(keys, func(i string) {
					s += t.Tests[i].GetHierString(nextSpace + 1)
				})
				return s
			}) +
			astHierFmt("Groups:", nextSpace, func() string {
				s := ""
				keys := make([]string, 0)
				for k := range t.Groups {
					keys = append(keys, k)
				}
				ForeachStringKeysInOrder(keys, func(i string) {
					s += fmt.Sprintln(strings.Repeat(" ", nextSpace+1) + t.Groups[i].Name)
				})
				return s
			})
	})
}

//------------------------

//Root
//------------------------
type astRoot struct {
	Env     *astEnv
	Options map[string]*astOption
	Builds  map[string]*astBuild
	Groups  map[string]*astGroup
}

func (t *astRoot) GetBuild(name string) *astBuild {
	if build, ok := t.Builds[name]; ok {
		return build
	}
	return nil
}

func (t *astRoot) GetGroup(name string) *astGroup {
	if group, ok := t.Groups[name]; ok {
		return group
	}
	return nil
}

func (t *astRoot) KeywordsChecker(s string) (bool, []string, string) {
	return true, nil, ""
}

func (t *astRoot) Parse(cfg map[interface{}]interface{}) error {
	//parsing Env
	if err := cfgToastItemRequired(cfg, "env", func(item interface{}) error {
		t.Env = new(astEnv)
		return astParse(t.Env, item.(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	//parsing builds
	if err := cfgToastItemRequired(cfg, "builds", func(item interface{}) error {
		t.Builds = make(map[string]*astBuild)
		for name, build := range item.(map[interface{}]interface{}) {
			t.Builds[name.(string)] = newAstBuild(name.(string))
			if err := astParse(t.Builds[name.(string)], (build.(map[interface{}]interface{}))); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	//parsing options
	if err := cfgToastItemOptional(cfg, "options", func(item interface{}) error {
		t.Options = make(map[string]*astOption)
		for name, option := range item.(map[interface{}]interface{}) {
			t.Options[name.(string)] = newAstOption(name.(string))
			if err := astParse(t.Options[name.(string)], option.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	//parsing groups
	if err := cfgToastItemOptional(cfg, "groups", func(item interface{}) error {
		t.Groups = make(map[string]*astGroup)
		for name, group := range item.(map[interface{}]interface{}) {
			t.Groups[name.(string)] = newAstGroup(name.(string))
			if err := astParse(t.Groups[name.(string)], group.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *astRoot) Link() error {
	//link groups
	for _, group := range t.Groups {
		if err := group.Link(); err != nil {
			return err
		}
	}
	return nil
}

func (t *astRoot) GetHierString(space int) string {
	nextSpace := space + 1
	return fmt.Sprintln(strings.Repeat(" ", space)+"astRoot") +
		astHierFmt("Env:", nextSpace, func() string {
			return t.Env.GetHierString(nextSpace + 1)
		}) +
		astHierFmt("Options:", nextSpace, func() string {
			s := ""
			keys := make([]string, 0)
			for k := range t.Options {
				keys = append(keys, k)
			}
			ForeachStringKeysInOrder(keys, func(i string) {
				s += t.Options[i].GetHierString(nextSpace + 1)
			})
			return s
		}) +
		astHierFmt("Builds:", nextSpace, func() string {
			s := ""
			keys := make([]string, 0)
			for k := range t.Builds {
				keys = append(keys, k)
			}
			ForeachStringKeysInOrder(keys, func(i string) {
				s += t.Builds[i].GetHierString(nextSpace + 1)
			})
			return s
		}) +
		astHierFmt("Groups:", nextSpace, func() string {
			s := ""
			keys := make([]string, 0)
			for k := range t.Groups {
				keys = append(keys, k)
			}
			ForeachStringKeysInOrder(keys, func(i string) {
				s += t.Groups[i].GetHierString(nextSpace + 1)
			})
			return s
		})

}

//global
var jvsAstRoot = astRoot{}
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)
//------------------------
