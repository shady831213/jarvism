package parser

import (
	"errors"
	"flag"
	"fmt"
	"github.com/shady831213/jarvisSim/utils"
	"strconv"
	"strings"
)

type astParser interface {
	//pass1:top-down AstParse
	Parse(map[interface{}]interface{}) error
	KeywordsChecker(s string) (bool, []string, string)
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

func AstParse(parser astParser, cfg map[interface{}]interface{}) error {
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

type astItem struct {
	content string
}

func newAstItem(content interface{}) *astItem {
	inst := new(astItem)
	if value, ok := content.(string); ok {
		inst.content = value
		return inst
	}
	if value, ok := content.([]interface{}); ok {
		for _, i := range (value) {
			s, ok := i.(string)
			if !ok {
				panic(fmt.Sprintf("content must be string or []interface{}, but it is %T !", content))
				return nil
			}
			inst.content += " " + s
		}

		return inst
	}
	panic(fmt.Sprintf("content must be string or []interface{}, but it is %T !", content))
	return nil
}

func (item *astItem) Cat(i *astItem) {
	if i == nil {
		return
	}
	item.content += i.content
}

func (item *astItem) Replace(old, new string, cnt int) {
	item.content = strings.Replace(item.content, old, new, cnt)
}

type astItems struct {
	Items map[string]*astItem
}

func (items *astItems) init() {
	items.Items = make(map[string]*astItem)
	items.Items["pre_sim_option"] = nil
	items.Items["sim_option"] = nil
	items.Items["post_sim_option"] = nil
	items.Items["pre_compile_option"] = nil
	items.Items["compile_option"] = nil
	items.Items["post_compile_option"] = nil
}

func (items *astItems) Cat(i *astItems) {
	if i == nil {
		return
	}
	for k, _ := range i.Items {
		if items.Items[k] == nil && i.Items[k] != nil {
			items.Items[k] = newAstItem("")
		}
		items.Items[k].Cat(i.Items[k])
	}
}

func (items *astItems) Replace(old, new string, cnt int) {
	for k, v := range items.Items {
		if v != nil {
			items.Items[k].Replace(old, new, cnt)
		}
	}
}

func (items *astItems) IsSimOnly() bool {
	return items.Items["pre_compile_option"] == nil && items.Items["compile_option"] == nil && items.Items["post_compile_option"] == nil
}

func (items *astItems) GetHierString(space int) string {
	nextSpace := space + 1
	s := ""
	keys := make([]string, 0)
	for k := range items.Items {
		keys = append(keys, k)
	}
	utils.ForeachStringKeysInOrder(keys, func(i string) {
		s += astHierFmt(i+":", nextSpace, func() string {
			return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
				fmt.Sprintln(items.Items[i])
		})
	})
	return s
}

type astParseItem struct {
	astItems
}

func (items *astParseItem) KeywordsChecker(s string) (bool, []string, string) {
	keywords := make(map[string]interface{})
	for k, _ := range items.Items {
		keywords[k] = nil
	}
	if !checkKeyWord(s, keywords) {
		return false, utils.KeyOfStringMap(keywords), ""
	}
	return true, nil, ""
}

func (items *astParseItem) Parse(cfg map[interface{}]interface{}) error {
	for k, _ := range items.Items {
		if err := cfgToastItemOptional(cfg, k, func(i interface{}) error {
			items.Items[k] = newAstItem(i)
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

//Options
//------------------------
type AstOptionAction struct {
	astParseItem
}

func NewAstOptionAction() *AstOptionAction {
	inst := new(AstOptionAction)
	inst.astParseItem.init()
	return inst
}

type JvsOptionForTest interface {
	JvsOption
	IsCompileOption() bool
	TestHandler(test *AstTestCase)
}

type AstOption struct {
	On        *AstOptionAction
	WithValue *AstOptionAction
	Value     string
	Name      string
}

func newAstOption(name string) *AstOption {
	inst := new(AstOption)
	inst.Init(name)
	return inst
}

func (t *AstOption) Init(name string) {
	t.Name = name
	t.Value = "false"
}

func (t *AstOption) GetName() string {
	return t.Name
}

func (t *AstOption) Clone() JvsOption {
	inst := newAstOption(t.Name)
	inst.Value = t.Value
	inst.On = t.On
	inst.WithValue = t.WithValue
	return inst
}

func (t *AstOption) Set(s string) error {
	if !t.IsBoolFlag() {
		t.Value = s
		return nil
	}
	b, err := strconv.ParseBool(s)
	t.Value = strconv.FormatBool(b)
	return err
}

func (t *AstOption) String() string {
	return t.Value
}

func (t *AstOption) IsBoolFlag() bool {
	return t.WithValue == nil
}

func (t *AstOption) IsCompileOption() bool {
	return t.WithValue != nil && !t.WithValue.IsSimOnly() || t.On != nil && !t.On.IsSimOnly()
}

func (t *AstOption) Usage() string {
	return "user-defined flag"
}

func (t *AstOption) TestHandler(test *AstTestCase) {
	if t.IsBoolFlag() {
		test.Cat(&t.On.astItems)
		return
	}
	t.WithValue.Replace("$"+t.Name, t.Value, -1)
	test.Cat(&t.WithValue.astItems)
}

func (t *AstOption) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"on_action": nil, "with_value_action": nil}
	if !checkKeyWord(s, keywords) {
		return false, utils.KeyOfStringMap(keywords), "Error in " + t.Name + ":"
	}
	return true, nil, ""
}

func (t *AstOption) Parse(cfg map[interface{}]interface{}) error {
	if err := cfgToastItemOptional(cfg, "on_action", func(item interface{}) error {
		if t.WithValue != nil {
			return errors.New("Error in " + t.Name + ": on_action and with_value_action are both defined!")
		}
		t.On = NewAstOptionAction()
		return AstParse(t.On, item.(map[interface{}]interface{}))
	}); err != nil {
		return astError("on_action of "+t.Name, err)
	}
	if err := cfgToastItemOptional(cfg, "with_value_action", func(item interface{}) error {
		if t.On != nil {
			return errors.New("Error in " + t.Name + ": on_action and with_value_action are both defined!")
		}
		t.WithValue = NewAstOptionAction()
		return AstParse(t.WithValue, item.(map[interface{}]interface{}))
	}); err != nil {
		return astError("with_value_action of "+t.Name, err)
	}
	//add to flagSet
	RegisterJvsOption(t, t.Usage())
	return nil
}

func (t *AstOption) GetHierString(space int) string {
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
}

func (t *astEnv) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"simulator": nil}
	if !checkKeyWord(s, keywords) {
		return false, utils.KeyOfStringMap(keywords), "Error in Env:"
	}
	return true, nil, ""
}

func (t *astEnv) Parse(cfg map[interface{}]interface{}) error {
	if err := cfgToastItemRequired(cfg, "simulator", func(item interface{}) error {
		simulator, ok := validSimulators[item.(string)]
		if !ok {
			errMsg := "Error in Env: simulator " + item.(string) + " is invalid! valid simulator are [ "
			for k, _ := range validSimulators {
				errMsg += k + " "
			}
			errMsg += "!"
			return errors.New(errMsg)
		}
		if err := LoadBuildInOptions(simulator.BuildInOptionFile); err != nil {
			panic("Error in loading " + simulator.BuildInOptionFile + ":" + err.Error())
		}
		setSimulator(simulator)
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
			fmt.Sprintln(GetSimulator().Name)
	})
}

//------------------------

//Build
//------------------------
type astBuild struct {
	astParseItem
	Name string
}

func newAstBuild(name string) *astBuild {
	inst := new(astBuild)
	inst.Name = name
	inst.astParseItem.init()
	return inst
}

func (t *astBuild) KeywordsChecker(s string) (bool, []string, string) {
	if ok, buildKeywords, _ := t.astParseItem.KeywordsChecker(s); !ok {
		return false, buildKeywords, "Error in build " + t.Name + ":"
	}

	return true, nil, ""
}

func (t *astBuild) Parse(cfg map[interface{}]interface{}) error {
	return t.astParseItem.Parse(cfg)
}

func (t *astBuild) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.astParseItem.GetHierString(nextSpace)
	})
}

//------------------------

//Test and Group, linkable
//------------------------
type astTestOpts interface {
	SetParent(parent astTestOpts)
	//bottom-up search
	GetOptionArgs() map[string]JvsOptionForTest
	//bottom-up search
	GetBuild() *astBuild
	//top-down search
	GetTestCases() []*AstTestCase
}

type astTest struct {
	Name       string
	Build      *astBuild
	OptionArgs map[string]JvsOptionForTest
	args       []string
	parent     astTestOpts
}

func (t *astTest) init(name string) {
	t.Name = name
	t.OptionArgs = make(map[string]JvsOptionForTest)
}

func (t *astTest) SetParent(parent astTestOpts) {
	t.parent = parent
}

func (t *astTest) GetOptionArgs() map[string]JvsOptionForTest {
	if t.parent != nil {
		options := make(map[string]JvsOptionForTest)
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

func (t *astTest) GetBuild() *astBuild {
	if t.Build != nil {
		return t.Build
	}
	if t.parent != nil {
		return t.parent.GetBuild()
	}
	return nil
}

func (t *astTest) KeywordsChecker(s string) (bool, []string, string) {
	keywords := map[string]interface{}{"args": nil}
	if !checkKeyWord(s, keywords) {
		return false, utils.KeyOfStringMap(keywords), "Error in " + t.Name + ":"
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
		opt, err := GetOption(arg)
		if err != nil {
			return astError("args of "+t.Name, err)
		}
		v, ok := opt.Clone().(JvsOptionForTest)
		if !ok {
			return nil
		}
		if v.IsCompileOption() {
			return errors.New("Error in args of " + t.Name + ":" + v.GetName() + " is a compile option! compile option can't be in groups and tests!")
		}
		t.OptionArgs[v.GetName()] = v

	}
	return nil
}

func (t *astTest) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("parent:", nextSpace, func() string {
		if t.parent != nil {
			return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.parent.(*astGroup).Name)
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
			utils.ForeachStringKeysInOrder(keys, func(i string) {
				if v, ok := args[i].(*AstOption); ok {
					s += v.GetHierString(nextSpace + 1)
				} else {
					s += fmt.Sprintln(strings.Repeat(" ", nextSpace) + "buildIn Option: " + args[i].GetName())
				}
			})
			return s
		})
}

type AstTestCase struct {
	astTest
	astItems
	seeds []int
}

func newAstTestCase(name string) *AstTestCase {
	inst := new(AstTestCase)
	inst.astTest.init(name)
	inst.astItems.init()
	return inst
}

func (t *AstTestCase) GetTestCases() []*AstTestCase {
	testcases := make([]*AstTestCase, len(t.seeds))
	for i := range testcases {
		testcases[i] = newAstTestCase(t.Name + "__" + strconv.Itoa(t.seeds[i]))
		//copy sim_options and set seed
		testcases[i].astItems.Cat(&t.astItems)
		testcases[i].Items["sim_option"].Cat(newAstItem(GetSimulator().SeedOption + strconv.Itoa(t.seeds[i])))
		//copy build
		testcases[i].Build = t.Build
	}
	return testcases
}

func (t *AstTestCase) Link() error {
	if err := t.astTest.Link(); err != nil {
		return err
	}
	//get build sim_options
	t.Build = t.GetBuild()
	t.Cat(&t.Build.astItems)

	//get options sim_options in order
	keys := make([]string, 0)
	args := t.GetOptionArgs()
	for k := range args {
		keys = append(keys, k)
	}
	utils.ForeachStringKeysInOrder(keys, func(i string) {
		args[i].TestHandler(t)
	})
	return nil
}

func (t *AstTestCase) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return t.astTest.GetHierString(nextSpace) +
			t.astItems.GetHierString(nextSpace) +
			astHierFmt("seeds:", nextSpace, func() string {
				return strings.Repeat(" ", nextSpace) + fmt.Sprintln(t.seeds)
			}) +
			astHierFmt("Builds:", nextSpace, func() string {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.GetBuild().Name)
			}) +
			astHierFmt("Flatten Tests:", nextSpace, func() string {
				s := ""
				for _, test := range t.GetTestCases() {
					s += test.GetHierString(nextSpace + 1)
				}
				return s
			})
	})
}

type astGroup struct {
	astTest
	buildName string
	Tests     map[string]*AstTestCase
	Groups    map[string]*astGroup
}

func newAstGroup(name string) *astGroup {
	inst := new(astGroup)
	inst.init(name)
	inst.buildName = ""
	return inst
}

func (t *astGroup) GetTestCases() []*AstTestCase {
	testcases := make([]*AstTestCase, 0)
	for _, test := range t.Tests {
		testcases = append(testcases, test.GetTestCases()...)
	}
	for _, group := range t.Groups {
		testcases = append(testcases, group.GetTestCases()...)
	}
	return testcases
}

func (t *astGroup) KeywordsChecker(s string) (bool, []string, string) {
	if ok, testKeywords, _ := t.astTest.KeywordsChecker(s); !ok {
		groupKeywords := map[string]interface{}{"build": nil, "tests": nil, "groups": nil}
		if !checkKeyWord(s, groupKeywords) {
			return false, append(testKeywords, utils.KeyOfStringMap(groupKeywords)...), "Error in group " + t.Name + ":"
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
	//AstParse tests
	if err := cfgToastItemOptional(cfg, "tests", func(item interface{}) error {
		t.Tests = make(map[string]*AstTestCase)
		for name, test := range item.(map[interface{}]interface{}) {
			t.Tests[name.(string)] = newAstTestCase(name.(string))
			if err := AstParse(t.Tests[name.(string)], test.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return astError("group "+t.Name, err)
	}

	//AstParse groups
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
	if t.buildName != "" {
		build := jvsAstRoot.GetBuild(t.buildName)
		if build == nil {
			return errors.New("build " + t.buildName + "of group " + t.Name + "is undef!")
		}
		t.Build = build
	}
	//link args
	if err := t.astTest.Link(); err != nil {
		return err
	}
	//link groups
	for name, _ := range t.Groups {
		group := jvsAstRoot.GetGroup(name)
		group.SetParent(t)
		if group == nil {
			return errors.New("sub group " + name + "of group " + t.Name + " is undef!")
		}
		t.Groups[name] = group
	}
	//check loop
	if t.parent != nil {
		for g := t.parent.(*astGroup); g.parent != nil; g = g.parent.(*astGroup) {
			if g.Name == t.Name {
				return errors.New("Loop include: group " + t.Name + " and group " + g.Name)
			}
		}
	}

	//link tests
	for _, test := range t.Tests {
		test.SetParent(t)
		if err := test.Link(); err != nil {
			return err
		}
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
				utils.ForeachStringKeysInOrder(keys, func(i string) {
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
				utils.ForeachStringKeysInOrder(keys, func(i string) {
					s += fmt.Sprintln(strings.Repeat(" ", nextSpace+1) + t.Groups[i].Name)
				})
				return s
			}) +
			astHierFmt("Flatten Tests:", nextSpace, func() string {
				s := ""
				for _, test := range t.GetTestCases() {
					s += test.GetHierString(nextSpace + 1)
				}
				return s
			})
	})
}

//------------------------

//Root
//------------------------
type astRoot struct {
	Env     *astEnv
	Options map[string]*AstOption
	Builds  map[string]*astBuild
	Groups  map[string]*astGroup
}

func newAstRoot() *astRoot {
	inst := new(astRoot)
	inst.Env = new(astEnv)
	inst.Builds = make(map[string]*astBuild)
	inst.Groups = make(map[string]*astGroup)
	inst.Options = make(map[string]*AstOption)
	return inst
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
		return AstParse(t.Env, item.(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	//parsing builds
	if err := cfgToastItemRequired(cfg, "builds", func(item interface{}) error {
		for name, build := range item.(map[interface{}]interface{}) {
			t.Builds[name.(string)] = newAstBuild(name.(string))
			if err := AstParse(t.Builds[name.(string)], (build.(map[interface{}]interface{}))); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	//parsing options
	if err := cfgToastItemOptional(cfg, "options", func(item interface{}) error {
		for name, option := range item.(map[interface{}]interface{}) {
			t.Options[name.(string)] = newAstOption(name.(string))
			if err := AstParse(t.Options[name.(string)], option.(map[interface{}]interface{})); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	//parsing groups
	if err := cfgToastItemOptional(cfg, "groups", func(item interface{}) error {
		for name, group := range item.(map[interface{}]interface{}) {
			t.Groups[name.(string)] = newAstGroup(name.(string))
			if err := AstParse(t.Groups[name.(string)], group.(map[interface{}]interface{})); err != nil {
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
			utils.ForeachStringKeysInOrder(keys, func(i string) {
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
			utils.ForeachStringKeysInOrder(keys, func(i string) {
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
			utils.ForeachStringKeysInOrder(keys, func(i string) {
				s += t.Groups[i].GetHierString(nextSpace + 1)
			})
			return s
		})

}

//global
var jvsAstRoot = newAstRoot()

func GetJvsAstRoot() *astRoot {
	return jvsAstRoot
}

//------------------------
