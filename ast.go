package jarivsSim

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type ASTParser interface {
	//pass1:top-down parse
	Parse(map[interface{}]interface{}) error
	//print:for debug
	GetHierString(space int) string;
}

type ASTLinker interface {
	ASTParser
	//pass2:top-down link
	Link(map[interface{}]interface{}) error
}

func astHierFmt(title string, space int, handler func() string) string {
	return fmt.Sprintln(strings.Repeat(" ", space)+strings.Repeat("-", 20-space)) +
		fmt.Sprintln(strings.Repeat(" ", space)+title+":") +
		handler() +
		"\n" +
		fmt.Sprintln(strings.Repeat(" ", space)+strings.Repeat("-", 20-space))
}

func CfgToASTItemRequired(cfg map[interface{}]interface{}, key string, handler func(interface{}) error) error {
	if item, ok := cfg[key]; ok {
		flag.Args()
		return handler(item)
	}
	return errors.New("not define " + key + "!")
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

func CfgToASTItemOptional(cfg map[interface{}]interface{}, key string, handler func(interface{}) error) error {
	if item, ok := cfg[key]; ok {
		return handler(item)
	}
	return nil
}

type ASTOptionItem struct {
	content string
}

func NewASTOptionItem(content interface{}) *ASTOptionItem {
	inst := new(ASTOptionItem)
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

type ASTSimOnlyItem struct {
	PreSimOption  *ASTOptionItem
	SimOption     *ASTOptionItem
	PostSimOption *ASTOptionItem
}

func (t *ASTSimOnlyItem) Parse(cfg map[interface{}]interface{}) error {
	if err := CfgToASTItemOptional(cfg, "pre_sim_option", func(item interface{}) error {
		t.PreSimOption = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "sim_option", func(item interface{}) error {
		t.SimOption = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "post_sim_option", func(item interface{}) error {
		t.PostSimOption = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *ASTSimOnlyItem) GetHierString(space int) string {
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

type ASTBuildItem struct {
	ASTSimOnlyItem
	PreCompileOption  *ASTOptionItem
	CompileOption     *ASTOptionItem
	PostCompileOption *ASTOptionItem
}

func (t *ASTBuildItem) Parse(cfg map[interface{}]interface{}) error {
	if err := t.ASTSimOnlyItem.Parse(cfg); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "pre_compile_option", func(item interface{}) error {
		t.PreCompileOption = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "compile_option", func(item interface{}) error {
		t.CompileOption = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "post_compile_option", func(item interface{}) error {
		t.PostCompileOption = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *ASTBuildItem) GetHierString(space int) string {
	nextSpace := space + 1
	return t.ASTSimOnlyItem.GetHierString(space) +
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
type ASTOptionAction struct {
	ASTBuildItem
}

func (t *ASTOptionAction) Parse(cfg map[interface{}]interface{}) error {
	return t.ASTBuildItem.Parse(cfg)
}

type ASTOption struct {
	On        *ASTOptionAction
	Off       *ASTOptionAction
	WithValue *ASTOptionAction
	Value     string
	Name      string
}

func NewASTOption(name string) *ASTOption {
	inst := new(ASTOption)
	inst.Name = name
	inst.Value = "false"
	return inst
}

func (t *ASTOption) Set(s string) error {
	if t.WithValue != nil {
		t.Value = s
		return nil
	}
	b, err := strconv.ParseBool(s)
	t.Value = strconv.FormatBool(b)
	return err
}

func (t *ASTOption) String() string {
	return t.Value
}

func (t *ASTOption) IsBoolFlag() bool {
	return t.WithValue == nil
}

func (t *ASTOption) Parse(cfg map[interface{}]interface{}) error {
	if err := CfgToASTItemOptional(cfg, "on_action", func(item interface{}) error {
		t.On = new(ASTOptionAction)
		fmt.Println("get an on action!")
		return t.On.Parse(item.(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "off_action", func(item interface{}) error {
		t.Off = new(ASTOptionAction)
		fmt.Println("get an off action!")
		return t.Off.Parse(item.(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	if err := CfgToASTItemOptional(cfg, "with_value_action", func(item interface{}) error {
		t.WithValue = new(ASTOptionAction)
		fmt.Println("get an with_value action!")
		return t.WithValue.Parse(item.(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	for key, _ := range(cfg) {
		if key.(string) != "on_action" && key.(string) != "off_action" && key.(string) != "with_value_action" {
			return errors.New("action of option " + t.Name + "must be \"on_action\", \"off_action\" or \"with_value_action\" but get \"" + key.(string) + "\"!" )
		}
	}
	//add to flagSet
	jvsOptions.Var(t, t.Name, "user-defined flag")
	return nil
}

func (t *ASTOption) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return astHierFmt("Value:", nextSpace, func() string {
			return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.Value)
		}) +
			astHierFmt("On:", nextSpace, func() string {
				if t.On != nil {
					return t.On.GetHierString(nextSpace + 1)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace + 1) + "null")
			}) +
			astHierFmt("Off:", nextSpace, func() string {
				if t.Off != nil {
					return t.Off.GetHierString(nextSpace + 1)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace + 1) + "null")
			}) +
			astHierFmt("WithValue:", nextSpace, func() string {
				if t.WithValue != nil {
					return t.WithValue.GetHierString(nextSpace + 1)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace + 1) + "null")
			})
	})
}

//------------------------

//env
//------------------------
type ASTEnv struct {
	CompileCmd *ASTOptionItem
	SimCmd     *ASTOptionItem
}

func (t *ASTEnv) Parse(cfg map[interface{}]interface{}) error {
	if err := CfgToASTItemRequired(cfg, "compile_cmd", func(item interface{}) error {
		t.CompileCmd = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToASTItemRequired(cfg, "sim_cmd", func(item interface{}) error {
		t.SimCmd = NewASTOptionItem(item)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *ASTEnv) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("CompileCmd:", nextSpace, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(t.CompileCmd)
	}) +
		astHierFmt("SimCmd:", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(t.SimCmd)
		})
}

//------------------------

//Build
//------------------------
type ASTBuild struct {
	ASTBuildItem
	Name string
}

func NewASTBuild(name string) *ASTBuild {
	inst := new(ASTBuild)
	inst.Name = name
	return inst
}

func (t *ASTBuild) Parse(cfg map[interface{}]interface{}) error {
	return t.ASTBuildItem.Parse(cfg)
}

func (t *ASTBuild) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.ASTBuildItem.GetHierString(nextSpace)
	})
}

//------------------------

//Test and Group, linkable
//------------------------
//bottom-up search
type aSTTestOpts interface {
	SetParent(parent aSTTestOpts)
	GetName() string
	GetOptionArgs() []*ASTOption
	GetBuild() *ASTBuild
}

type aSTTest struct {
	Name       string
	OptionArgs []*ASTOption
	parent     aSTTestOpts
}

func (t *aSTTest) init(name string) {
	t.Name = name
}

func (t *aSTTest) SetParent(parent aSTTestOpts) {
	t.parent = parent
}

func (t *aSTTest) GetName() string {
	return t.Name
}

func (t *aSTTest) GetOptionArgs() []*ASTOption {
	if t.parent != nil {
		return append(t.parent.GetOptionArgs(), t.OptionArgs...)
	}
	return t.OptionArgs
}

func (t *aSTTest) Parse(cfg map[interface{}]interface{}) error {
	return nil
}

func (t *aSTTest) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return astHierFmt("parent:", nextSpace, func() string {
			if t.parent != nil {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.parent.GetName())
			}
			return fmt.Sprintln(strings.Repeat(" ", nextSpace) + "null")
		}) +
			astHierFmt("OptionArgs:", nextSpace, func() string {
				s := ""
				for _, option := range t.GetOptionArgs() {
					s += fmt.Sprintln(strings.Repeat(" ", nextSpace) + option.Name)
				}
				return s
			})
	})
}

//because Link is top-down, the last repeated args take effect
func (t *aSTTest) Link(cfg map[interface{}]interface{}) error {
	if err := CfgToASTItemOptional(cfg, "args", func(item interface{}) error {
		t.OptionArgs = make([]*ASTOption, 0)
		for _, arg := range (item.([]interface{})) {
			//Options have been all parsed
			args := strings.Split(arg.(string), " ")
			if err := jvsOptions.Parse(args); err != nil {
				return err
			}
			optName, err := argToOption(args[0])
			if err != nil {
				return err
			}
			t.OptionArgs = append(t.OptionArgs, jvsOptions.Lookup(optName).Value.(*ASTOption))
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type ASTTestCase struct {
	aSTTest
}

func NewASTTestCase(name string) *ASTTestCase {
	inst := new(ASTTestCase)
	inst.init(name)
	return inst
}

func (t *ASTTestCase) GetBuild() *ASTBuild {
	return t.parent.GetBuild()
}

func (t *ASTTestCase) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.aSTTest.GetHierString(nextSpace) +
			astHierFmt("Builds:", nextSpace, func() string {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.GetBuild().Name)
			})
	})
}

type ASTGroup struct {
	aSTTest
	Build  *ASTBuild
	Tests  map[string]*ASTTestCase
	Groups []*ASTGroup
}

func NewASTGroup(name string) *ASTGroup {
	inst := new(ASTGroup)
	inst.init(name)
	return inst
}

func (t *ASTGroup) GetBuild() *ASTBuild {
	if t.Build != nil {
		return t.Build
	}
	if t.parent != nil {
		return t.parent.GetBuild()
	}
	return nil
}

func (t *ASTGroup) Parse(cfg map[interface{}]interface{}) error {
	if err := t.aSTTest.Parse(cfg); err != nil {
		return err
	}
	//parse tests
	if err := CfgToASTItemOptional(cfg, "tests", func(item interface{}) error {
		t.Tests = make(map[string]*ASTTestCase)
		for name, test := range item.(map[interface{}]interface{}) {
			t.Tests[name.(string)] = NewASTTestCase(name.(string))
			if err := t.Tests[name.(string)].Parse(test.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *ASTGroup) Link(cfg map[interface{}]interface{}) error {
	//link build
	//builds have been all parsed
	if err := CfgToASTItemOptional(cfg, "build", func(item interface{}) error {
		build := jvsASTRoot.GetBuild(item.(string))
		if build == nil {
			return errors.New(item.(string) + " is undef!")
		}
		t.Build = build
		return nil
	}); err != nil {
		return err
	}
	//link args
	if err := t.aSTTest.Link(cfg); err != nil {
		return err
	}
	//link tests
	for name, test := range t.Tests {
		if err := test.Link(cfg["tests"].(map[interface{}]interface{})[name].(map[interface{}]interface{})); err != nil {
			return err
		}
		test.SetParent(t)
	}
	//link groups
	if err := CfgToASTItemOptional(cfg, "groups", func(item interface{}) error {
		t.Groups = make([]*ASTGroup, 0)
		for _, name := range item.([]interface{}) {
			group := jvsASTRoot.GetGroup(name.(string))
			if group == nil {
				return errors.New("group " + name.(string) + " is undef!")
			}
			t.Groups = append(t.Groups, group)
			group.SetParent(t)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *ASTGroup) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.aSTTest.GetHierString(nextSpace) +
			astHierFmt("Builds:", nextSpace, func() string {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.GetBuild().Name)
			}) +
			astHierFmt("Tests:", nextSpace, func() string {
				s := ""
				for _, test := range t.Tests {
					s += test.GetHierString(nextSpace + 1)
				}
				return s
			}) +
			astHierFmt("Groups:", nextSpace, func() string {
				s := ""
				for _, group := range t.Groups {
					s += fmt.Sprintln(strings.Repeat(" ", nextSpace) + group.Name)
				}
				return s
			})
	})
}

//------------------------

//Root
//------------------------
type ASTRoot struct {
	Env     *ASTEnv
	Options map[string]*ASTOption
	Builds  map[string]*ASTBuild
	Groups  map[string]*ASTGroup
}

func (t *ASTRoot) GetBuild(name string) *ASTBuild {
	if build, ok := t.Builds[name]; ok {
		return build
	}
	return nil
}

func (t *ASTRoot) GetGroup(name string) *ASTGroup {
	if group, ok := t.Groups[name]; ok {
		return group
	}
	return nil
}

func (t *ASTRoot) Parse(cfg map[interface{}]interface{}) error {
	//parsing Env
	if err := CfgToASTItemRequired(cfg, "env", func(item interface{}) error {
		t.Env = new(ASTEnv)
		return t.Env.Parse(item.(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	//parsing builds
	if err := CfgToASTItemRequired(cfg, "builds", func(item interface{}) error {
		t.Builds = make(map[string]*ASTBuild)
		for name, build := range item.(map[interface{}]interface{}) {
			t.Builds[name.(string)] = NewASTBuild(name.(string))
			if err := t.Builds[name.(string)].Parse(build.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	//parsing options
	if err := CfgToASTItemOptional(cfg, "options", func(item interface{}) error {
		t.Options = make(map[string]*ASTOption)
		for name, option := range item.(map[interface{}]interface{}) {
			t.Options[name.(string)] = NewASTOption(name.(string))
			if err := t.Options[name.(string)].Parse(option.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	//parsing groups
	if err := CfgToASTItemOptional(cfg, "groups", func(item interface{}) error {
		t.Groups = make(map[string]*ASTGroup)
		for name, group := range item.(map[interface{}]interface{}) {
			t.Groups[name.(string)] = NewASTGroup(name.(string))
			if err := t.Groups[name.(string)].Parse(group.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (t *ASTRoot) Link(cfg map[interface{}]interface{}) error {
	//link groups
	for name, group := range t.Groups {
		if err := group.Link(cfg["groups"].(map[interface{}]interface{})[name].(map[interface{}]interface{})); err != nil {
			return err
		}
	}
	return nil
}

func (t *ASTRoot) GetHierString(space int) string {
	nextSpace := space + 1
	return fmt.Sprintln(strings.Repeat(" ", space)+"ASTRoot") +
		astHierFmt("Env:", nextSpace, func() string {
			return t.Env.GetHierString(nextSpace + 1)
		}) +
		astHierFmt("Options:", nextSpace, func() string {
			s := ""
			for _, option := range t.Options {
				s += option.GetHierString(nextSpace + 1)
			}
			return s
		}) +
		astHierFmt("Builds:", nextSpace, func() string {
			s := ""
			for _, build := range t.Builds {
				s += build.GetHierString(nextSpace + 1)
			}
			return s
		}) +
		astHierFmt("Groups:", nextSpace, func() string {
			s := ""
			for _, group := range t.Groups {
				s += group.GetHierString(nextSpace + 1)
			}
			return s
		})

}

//global
var jvsASTRoot = ASTRoot{}
var jvsOptions = flag.NewFlagSet("jvsOptions", flag.ExitOnError)
//------------------------
