package ast

import (
	"flag"
	"fmt"
	"github.com/fatih/set"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/errors"
	"github.com/shady831213/jarvism/core/utils"
	"math"
	"path"
	"sort"
	"strconv"
	"strings"
)

type astParser interface {
	//pass1:top-down AstParse
	Parse(map[interface{}]interface{}) *errors.JVSAstError
	KeywordsChecker(string) (bool, *utils.StringMapSet, string)
}

type astLinker interface {
	astParser
	//pass2:top-down link
	Link() *errors.JVSAstError
}

func CheckKeyWord(s string, keyWords *utils.StringMapSet) bool {
	_, ok := keyWords.Get(s)
	return ok
}

func astHierFmt(title string, space int, handler func() string) string {
	return fmt.Sprintln(strings.Repeat(" ", space)+strings.Repeat("-", 20-space)) +
		fmt.Sprintln(strings.Repeat(" ", space)+title) +
		handler() +
		"\n"
}

func AstParse(parser astParser, cfg map[interface{}]interface{}) *errors.JVSAstError {
	for name := range cfg {
		if ok, keywords, tag := parser.KeywordsChecker(name.(string)); !ok {
			return errors.JVSAstParseError("", tag+"syntax error of \""+name.(string)+"\"! expect "+fmt.Sprint(keywords.Keys()))
		}
	}
	if err := parser.Parse(cfg); err != nil {
		return err
	}
	return nil
}

func CfgToAstItemRequired(cfg map[interface{}]interface{}, key string, handler func(interface{}) *errors.JVSAstError) *errors.JVSAstError {
	if item, ok := cfg[key]; ok {
		flag.Args()
		return handler(item)
	}
	return errors.JVSAstParseError("", "not define "+key+"!")
}

func CfgToAstItemOptional(cfg map[interface{}]interface{}, key string, handler func(interface{}) *errors.JVSAstError) *errors.JVSAstError {
	if item, ok := cfg[key]; ok {
		return handler(item)
	}
	return nil
}

func astLoadPlugin(pluginType JVSPluginType, pluginName string) *errors.JVSAstError {
	if err := loadPlugin(pluginType, pluginName); err != nil {
		errMsg := string(pluginType) + " " + pluginName + " is invalid! valid " + string(pluginType) + "s are [ "
		for _, k := range validPlugins(pluginType) {
			errMsg += k + " "
		}
		errMsg += "]!"
		return errors.JVSAstParseError("", errMsg+"\n"+err.Error())
	}
	return nil
}

type astItem struct {
	content *set.SetNonTS
}

func newAstItem(content interface{}) *astItem {
	inst := new(astItem)
	inst.content = set.New(set.NonThreadSafe).(*set.SetNonTS)
	if value, ok := content.(string); ok {
		inst.content.Add(value)
		return inst
	}
	if value, ok := content.([]interface{}); ok {
		for _, i := range value {
			s, ok := i.(string)
			if !ok {
				panic(fmt.Sprintf("content must be string or []interface{}, but it is %T !", content))
				return nil
			}
			inst.content.Add(s)
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
	item.content.Merge(i.content)
}

func (item *astItem) Replace(old, new string, cnt int) *astItem {
	inst := newAstItem("")
	item.content.Each(func(i interface{}) bool {
		s := strings.Replace(i.(string), old, new, cnt)
		inst.content.Add(s)
		return true
	})
	return inst
}

func (item *astItem) GetString() string {
	l := set.StringSlice(item.content)
	sort.Strings(l)
	return strings.Join(l, " ")
}

type astItems struct {
	name       string
	preAction  string
	postAction string
	option     *astItem
}

func newAstItems(name string) *astItems {
	inst := new(astItems)
	inst.name = name
	inst.option = newAstItem("")
	return inst
}

func (items *astItems) preActionName() string {
	return "pre_" + items.name + "_action"
}

func (items *astItems) postActionName() string {
	return "post_" + items.name + "_action"
}

func (items *astItems) optionName() string {
	return items.name + "_option"
}

func (items *astItems) Cat(i *astItems) {
	if i == nil {
		return
	}
	items.preAction += " " + i.preAction
	items.postAction += " " + i.postAction
	items.option.Cat(i.option)
}

func (items *astItems) Replace(old, new string, cnt int) *astItems {
	inst := newAstItems(items.name)
	inst.preAction = strings.Replace(items.preAction, old, new, cnt)
	inst.postAction = strings.Replace(items.postAction, old, new, cnt)
	inst.option = items.option.Replace(old, new, cnt)
	return inst
}

func (items *astItems) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey(items.preActionName())
	keywords.AddKey(items.optionName())
	keywords.AddKey(items.postActionName())
	if !CheckKeyWord(s, keywords) {
		return false, keywords, ""
	}
	return true, nil, ""
}

func (items *astItems) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := CfgToAstItemOptional(cfg, items.preActionName(), func(i interface{}) *errors.JVSAstError {
		if l, ok := i.([]interface{}); ok {
			for _, s := range l {
				items.preAction += s.(string) + "\n"
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToAstItemOptional(cfg, items.optionName(), func(i interface{}) *errors.JVSAstError {
		items.option.Cat(newAstItem(i))
		return nil
	}); err != nil {
		return err
	}
	if err := CfgToAstItemOptional(cfg, items.postActionName(), func(i interface{}) *errors.JVSAstError {
		if l, ok := i.([]interface{}); ok {
			for _, s := range l {
				items.postAction += s.(string) + "\n"
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (items *astItems) GetHierString(space int) string {
	return astHierFmt(items.preActionName()+":", space, func() string {
		return fmt.Sprint(strings.Repeat(" ", space) + items.preAction)
	}) + astHierFmt(items.optionName()+":", space, func() string {
		return fmt.Sprint(strings.Repeat(" ", space) + items.option.GetString())
	}) + astHierFmt(items.postActionName()+":", space, func() string {
		return fmt.Sprint(strings.Repeat(" ", space) + items.postAction)
	})
}

//Plugins
//------------------------
type astPlugin struct {
	plugin     pluginOpts
	pluginType JVSPluginType
	attr       map[interface{}]interface{}
}

func (t *astPlugin) init(pluginType JVSPluginType) {
	t.pluginType = pluginType
	t.attr = make(map[interface{}]interface{})
}

func (t *astPlugin) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("type", "attr")
	if !CheckKeyWord(s, keywords) {
		return false, keywords, "Error in " + string(t.pluginType) + ":"
	}
	return true, nil, ""
}

func (t *astPlugin) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := CfgToAstItemRequired(cfg, "type", func(item interface{}) *errors.JVSAstError {
		if t.plugin = getPlugin(t.pluginType, item.(string)); t.plugin == nil {
			if err := astLoadPlugin(t.pluginType, item.(string)); err != nil {
				return errors.JVSAstParseError(string(t.pluginType), err.Error())
			}
			t.plugin = getPlugin(t.pluginType, item.(string))
		}
		return nil
	}); err != nil {
		return errors.JVSAstParseError(string(t.pluginType), err.Msg)
	}
	if err := CfgToAstItemOptional(cfg, "attr", func(item interface{}) *errors.JVSAstError {
		t.attr = item.(map[interface{}]interface{})
		return nil
	}); err != nil {
		return errors.JVSAstParseError(string(t.pluginType), err.Msg)
	}
	//parse discoverer
	return AstParse(t.plugin, t.attr)
}

func (t *astPlugin) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(string(t.pluginType), space, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(t.plugin.Name())
	}) + astHierFmt(string(t.pluginType)+"_attr:", nextSpace, func() string {
		s := ""
		keys := make([]string, 0)
		for k := range t.attr {
			keys = append(keys, k.(string))
		}
		utils.ForeachStringKeysInOrder(keys, func(i string) {
			if v, ok := t.attr[i]; ok {
				s += fmt.Sprint(strings.Repeat(" ", nextSpace) + fmt.Sprint(i) + ": " + fmt.Sprintln(v))
			}
		})
		return s
	})
}

type astTestDiscoverer struct {
	astPlugin
}

func newAstTestDiscoverer() *astTestDiscoverer {
	inst := new(astTestDiscoverer)
	inst.init(JVSTestDiscovererPlugin)
	return inst
}

type astChecker struct {
	astPlugin
}

func newAstChecker() *astChecker {
	inst := new(astChecker)
	inst.init(JVSCheckerPlugin)
	return inst
}

//Options
//------------------------
type AstOptionAction struct {
	compileOption, simOption *astItem
}

func NewAstOptionAction() *AstOptionAction {
	inst := new(AstOptionAction)
	inst.compileOption = newAstItem("")
	inst.simOption = newAstItem("")
	return inst
}

func (a *AstOptionAction) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("compile_option")
	keywords.AddKey("sim_option")
	if !CheckKeyWord(s, keywords) {
		return false, keywords, ""
	}
	return true, nil, ""
}

func (a *AstOptionAction) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := CfgToAstItemOptional(cfg, "compile_option", func(i interface{}) *errors.JVSAstError {
		a.compileOption.Cat(newAstItem(i))
		return nil
	}); err != nil {
		return err
	}

	if err := CfgToAstItemOptional(cfg, "sim_option", func(i interface{}) *errors.JVSAstError {
		a.simOption.Cat(newAstItem(i))
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (a *AstOptionAction) GetHierString(space int) string {
	return astHierFmt("compile_option:", space, func() string {
		return fmt.Sprint(strings.Repeat(" ", space) + a.compileOption.GetString())
	}) + astHierFmt("sim_option:", space, func() string {
		return fmt.Sprint(strings.Repeat(" ", space) + a.simOption.GetString())
	})
}

type AstOption struct {
	On        *AstOptionAction
	WithValue *AstOptionAction
	Value     string
	Name      string
	usage     string
}

func newAstOption(name string) *AstOption {
	inst := new(AstOption)
	inst.Init(name)
	inst.usage = "user-defined flag"
	return inst
}

func (t *AstOption) Init(name string) {
	t.Name = name
	t.Value = "false"
}

func (t *AstOption) GetName() string {
	return t.Name
}

func (t *AstOption) Clone() JvsAstOption {
	inst := newAstOption(t.Name)
	inst.Value = t.Value
	inst.On = t.On
	inst.WithValue = t.WithValue
	return inst
}

func (t *AstOption) Set(s string) error {
	t.Value = s
	return nil
}

func (t *AstOption) String() string {
	return t.Value
}

func (t *AstOption) IsBoolFlag() bool {
	return t.On != nil
}

func (t *AstOption) Usage() string {
	return t.usage
}

func (t *AstOption) TestHandler(test *AstTestCase) {
	if t.Value == "true" && t.IsBoolFlag() {
		test.simItems.option.Cat(t.On.simOption)
		return
	}
	test.simItems.option.Cat(t.WithValue.simOption.Replace("$"+t.Name, t.Value, -1))
}

func (t *AstOption) BuildHandler(build *AstBuild) {
	if t.Value == "true" && t.IsBoolFlag() {
		build.compileItems.option.Cat(t.On.compileOption)
		return
	}
	build.compileItems.option.Cat(t.WithValue.compileOption.Replace("$"+t.Name, t.Value, -1))
}

func (t *AstOption) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("on_action", "with_value_action", "usage")
	if !CheckKeyWord(s, keywords) {
		return false, keywords, "Error in " + t.Name + ":"
	}
	return true, nil, ""
}

func (t *AstOption) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := CfgToAstItemOptional(cfg, "usage", func(item interface{}) *errors.JVSAstError {
		t.usage = item.(string)
		return nil
	}); err != nil {
		return errors.JVSAstParseError("usage of "+t.Name, err.Msg)
	}
	if err := CfgToAstItemOptional(cfg, "on_action", func(item interface{}) *errors.JVSAstError {
		t.On = NewAstOptionAction()
		return AstParse(t.On, item.(map[interface{}]interface{}))
	}); err != nil {
		return errors.JVSAstParseError("on_action of "+t.Name, err.Msg)
	}
	if err := CfgToAstItemOptional(cfg, "with_value_action", func(item interface{}) *errors.JVSAstError {
		t.WithValue = NewAstOptionAction()
		return AstParse(t.WithValue, item.(map[interface{}]interface{}))
	}); err != nil {
		return errors.JVSAstParseError("with_value_action of "+t.Name, err.Msg)
	}
	if t.WithValue == nil && t.On == nil {
		return errors.JVSAstParseError(t.Name, "on_action and with_value_action are both nil!")
	}
	//add to flagSet
	RegisterJvsAstOption(t)
	return nil
}

func (t *AstOption) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", space, func() string {
		return astHierFmt("Usage:", nextSpace, func() string {
			return fmt.Sprintln(strings.Repeat(" ", nextSpace+1) + t.Usage())
		}) + astHierFmt("On:", nextSpace, func() string {
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

func (t *astEnv) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("simulator", "work_dir", "runner")
	if !CheckKeyWord(s, keywords) {
		return false, keywords, "Error in Env:"
	}
	return true, nil, ""
}

func (t *astEnv) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	//must parse work_dir first!
	if err := CfgToAstItemOptional(cfg, "work_dir", func(item interface{}) *errors.JVSAstError {
		if err := setWorkDir(item.(string)); err != nil {
			return errors.JVSAstParseError("work_dir in Env", err.Error())
		}
		return nil
	}); err != nil {
		return errors.JVSAstParseError("work_dir of Env", err.Msg)
	}
	//use default
	if GetWorkDir() == "" {
		if err := setWorkDir(path.Join(core.GetPrjHome(), "work")); err != nil {
			return errors.JVSAstParseError("work_dir in Env", err.Error())
		}
	}

	if err := CfgToAstItemOptional(cfg, "simulator", func(item interface{}) *errors.JVSAstError {
		simulator, ok := validSimulators[item.(string)]
		if !ok {
			if err := astLoadPlugin(JVSSimulatorPlugin, item.(string)); err != nil {
				return errors.JVSAstParseError("simulator of Env", err.Error())
			}
			simulator, _ = validSimulators[item.(string)]
		}
		setSimulator(simulator)
		return nil
	}); err != nil {
		return errors.JVSAstParseError("simulator of Env", err.Msg)
	}
	//use default
	if GetSimulator() == nil {
		if err := astLoadPlugin(JVSSimulatorPlugin, "vcs"); err != nil {
			return errors.JVSAstParseError("simulator of Env", err.Error())
		}
		simulator, _ := validSimulators["vcs"]
		setSimulator(simulator)
	}
	if err := LoadBuildInOptions(GetSimulator().BuildInOptionFile()); err != nil {
		panic("Error in loading " + GetSimulator().BuildInOptionFile() + ":" + err.Error())
	}

	if err := CfgToAstItemOptional(cfg, "runner", func(item interface{}) *errors.JVSAstError {
		runner, ok := validRunners[item.(string)]
		if !ok {
			if err := astLoadPlugin(JVSRunnerPlugin, item.(string)); err != nil {
				return errors.JVSAstParseError("runner of Env", err.Error())
			}
			runner, _ = validRunners[item.(string)]
		}
		setRunner(runner)
		return nil
	}); err != nil {
		return errors.JVSAstParseError("runner of Env", err.Msg)
	}
	//use default
	if GetRunner() == nil {
		if err := astLoadPlugin(JVSRunnerPlugin, "host"); err != nil {
			return errors.JVSAstParseError("runner of Env", err.Error())
		}
		runner, _ := validRunners["host"]
		setRunner(runner)
	}

	return nil
}

func (t *astEnv) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("Simulator:", nextSpace, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(GetSimulator().Name())
	}) + astHierFmt("Runner:", nextSpace, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(GetRunner().Name())
	}) + astHierFmt("WorkDir:", nextSpace, func() string {
		return fmt.Sprint(strings.Repeat(" ", nextSpace)) +
			fmt.Sprintln(GetWorkDir())
	})
}

//------------------------

//build
//------------------------

type AstBuild struct {
	Name                        string
	compileItems, simItems      *astItems
	testDiscoverer              *astTestDiscoverer
	compileChecker, testChecker *astChecker
}

func newAstBuild(name string) *AstBuild {
	inst := new(AstBuild)
	inst.Name = name
	inst.simItems = newAstItems("sim")
	inst.compileItems = newAstItems("compile")
	return inst
}

func (t *AstBuild) GetRawSign() string {
	return strings.Replace(t.Name+t.PreCompileAction()+t.CompileOption()+t.PostCompileAction(), " ", "", -1)
}

func (t *AstBuild) PreCompileAction() string {
	return t.compileItems.preAction
}

func (t *AstBuild) CompileOption() string {
	return t.compileItems.option.GetString()
}

func (t *AstBuild) PostCompileAction() string {
	return t.compileItems.postAction
}

func (t *AstBuild) Clone() *AstBuild {
	inst := newAstBuild(t.Name)
	inst.testDiscoverer = t.testDiscoverer
	inst.compileChecker = t.compileChecker
	inst.testChecker = t.testChecker
	inst.simItems.Cat(t.simItems)
	inst.compileItems.Cat(t.compileItems)
	return inst
}

func (t *AstBuild) GetTestDiscoverer() TestDiscoverer {
	return t.testDiscoverer.plugin.(TestDiscoverer)
}

func (t *AstBuild) GetChecker() Checker {
	return GetChecker(t.compileChecker.plugin.Name())
}

func (t *AstBuild) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	if ok, compileKeywors, _ := t.compileItems.KeywordsChecker(s); !ok {
		if ok, simKeywors, _ := t.simItems.KeywordsChecker(s); !ok {
			keywords := utils.NewStringMapSet()
			keywords.AddKey("test_discoverer")
			keywords.AddKey("compile_checker")
			keywords.AddKey("test_checker")
			if CheckKeyWord(s, keywords) {
				return true, nil, ""
			}
			return false, utils.StringMapSetUnion(compileKeywors, simKeywors, keywords), "Error in build " + t.Name + ":"
		}
	}
	return true, nil, ""
}

func (t *AstBuild) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := CfgToAstItemOptional(cfg, "test_discoverer", func(item interface{}) *errors.JVSAstError {
		t.testDiscoverer = newAstTestDiscoverer()
		return AstParse(t.testDiscoverer, item.(map[interface{}]interface{}))
	}); err != nil {
		return errors.JVSAstParseError(err.Item+" of build "+t.Name, err.Msg)
	}
	//use default
	if t.testDiscoverer == nil {
		t.testDiscoverer = newAstTestDiscoverer()
		if err := AstParse(t.testDiscoverer, map[interface{}]interface{}{"type": "uvm_test"}); err != nil {
			return errors.JVSAstParseError(err.Item+" of build "+t.Name, err.Msg)
		}
	}
	if err := CfgToAstItemOptional(cfg, "test_checker", func(item interface{}) *errors.JVSAstError {
		t.testChecker = newAstChecker()
		return AstParse(t.testChecker, item.(map[interface{}]interface{}))
	}); err != nil {
		return errors.JVSAstParseError(err.Item+" of build "+t.Name, err.Msg)
	}
	//use default
	if t.testChecker == nil {
		t.testChecker = newAstChecker()
		if err := AstParse(t.testChecker, map[interface{}]interface{}{"type": "testChecker"}); err != nil {
			return errors.JVSAstParseError(err.Item+" of build "+t.Name, err.Msg)
		}
	}

	if err := CfgToAstItemOptional(cfg, "compile_checker", func(item interface{}) *errors.JVSAstError {
		t.compileChecker = newAstChecker()
		return AstParse(t.compileChecker, item.(map[interface{}]interface{}))
	}); err != nil {
		return errors.JVSAstParseError(err.Item+" of build "+t.Name, err.Msg)
	}
	//use default
	if t.compileChecker == nil {
		t.compileChecker = newAstChecker()
		if err := AstParse(t.compileChecker, map[interface{}]interface{}{"type": "compileChecker"}); err != nil {
			return errors.JVSAstParseError(err.Item+" of build "+t.Name, err.Msg)
		}
	}

	//options
	if err := t.compileItems.Parse(cfg); err != nil {
		return errors.JVSAstParseError("build "+t.Name, err.Msg)
	}
	if err := t.simItems.Parse(cfg); err != nil {
		return errors.JVSAstParseError("build "+t.Name, err.Msg)
	}
	return nil
}

func (t *AstBuild) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.Name+":", nextSpace, func() string {
		return t.compileItems.GetHierString(nextSpace) +
			t.simItems.GetHierString(nextSpace) +
			t.testDiscoverer.GetHierString(nextSpace) +
			t.testChecker.GetHierString(nextSpace) +
			t.compileChecker.GetHierString(nextSpace)
	})
}

//------------------------

//Test and Group, linkable
//------------------------

type RunTimeOpts interface {
	GetName() string
	//bottom-up search
	GetBuild() *AstBuild
	SetBuild(*AstBuild)
	//top-down search
	GetTestCases() []*AstTestCase
	ParseArgs()
}

type astTestOpts interface {
	RunTimeOpts
	SetParent(parent astTestOpts)
	//bottom-up search
	GetOptionArgs() *utils.StringMapSet
}

type astTest struct {
	Name       string
	buildName  string
	build      *AstBuild
	optionArgs *utils.StringMapSet
	args       []string
	parent     astTestOpts
}

func (t *astTest) init(name string) {
	t.Name = name
	t.args = make([]string, 0)
	t.optionArgs = utils.NewStringMapSet()
}

func (t *astTest) Copy(i *astTest) {
	t.Name = i.Name
	t.buildName = i.buildName
	t.build = i.build
	//shared
	t.optionArgs = i.optionArgs
	//shared
	t.args = t.args
	t.parent = i.parent
}

func (t *astTest) GetName() string {
	if t.parent != nil {
		return t.parent.GetName() + "__" + t.Name
	}
	return t.Name
}

func (t *astTest) SetParent(parent astTestOpts) {
	t.parent = parent
}

func (t *astTest) GetOptionArgs() *utils.StringMapSet {
	if t.parent != nil {
		return utils.StringMapSetUnion(t.parent.GetOptionArgs(), t.optionArgs)
	}
	return t.optionArgs
}

func (t *astTest) GetBuild() *AstBuild {
	if t.build != nil {
		return t.build
	}
	if t.parent != nil {
		return t.parent.GetBuild()
	}
	return nil
}

func (t *astTest) SetBuild(build *AstBuild) {
	t.build = build
}

func (t *astTest) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	keywords := utils.NewStringMapSet()
	keywords.AddKey("build", "args")
	if !CheckKeyWord(s, keywords) {
		return false, keywords, "Error in " + t.Name + ":"
	}
	return true, nil, ""
}

func (t *astTest) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := CfgToAstItemOptional(cfg, "build", func(item interface{}) *errors.JVSAstError {
		t.buildName = item.(string)
		return nil
	}); err != nil {
		return errors.JVSAstParseError(t.Name, err.Msg)
	}
	if err := CfgToAstItemOptional(cfg, "args", func(item interface{}) *errors.JVSAstError {
		for _, arg := range item.([]interface{}) {
			t.args = append(t.args, strings.Split(arg.(string), ",")...)
		}
		return nil
	}); err != nil {
		return errors.JVSAstParseError(t.Name, err.Msg)
	}
	return nil
}

//because Link is top-down, the last repeated args take effect
func (t *astTest) Link() *errors.JVSAstError {
	//link build
	//builds have been all parsed
	if t.buildName != "" {
		build := jvsAstRoot.GetBuild(t.buildName)
		if build == nil {
			return errors.JVSAstLinkError(t.Name, "build "+t.buildName+" of "+t.Name+"is undef!")
		}
		t.build = build
	}
	for _, arg := range t.args {
		//Options have been all parsed
		opt, err := GetJvsAstOption(arg)
		if err != nil {
			return errors.JVSAstLinkError("args of "+t.Name, err.Error())
		}
		t.optionArgs.Add(opt.GetName(), opt.Clone())

	}
	return nil
}

func (t *astTest) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt("parent:", nextSpace, func() string {
		if t.parent != nil {
			return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.parent.(*AstGroup).Name)
		}
		return fmt.Sprintln(strings.Repeat(" ", nextSpace) + "null")
	}) +
		astHierFmt("optionArgs:", nextSpace, func() string {
			s := ""
			for _, arg := range t.GetOptionArgs().SortedList() {
				if v, ok := arg.(*AstOption); ok {
					s += v.GetHierString(nextSpace + 1)
				} else {
					s += fmt.Sprintln(strings.Repeat(" ", nextSpace) + "buildIn Option: " + arg.(JvsAstOptionForTest).GetName())
				}
			}
			return s
		})
}

type AstTestCase struct {
	astTest
	simItems *astItems
	seeds    []int
}

func newAstTestCase(name string) *AstTestCase {
	inst := new(AstTestCase)
	inst.astTest.init(name)
	inst.simItems = newAstItems("sim")
	return inst
}

func (t *AstTestCase) PreSimAction() string {
	return t.simItems.preAction
}

func (t *AstTestCase) SimOption() string {
	return t.simItems.option.GetString()
}

func (t *AstTestCase) PostSimAction() string {
	return t.simItems.postAction
}

func (t *AstTestCase) GetChecker() Checker {
	return GetChecker(t.build.testChecker.plugin.Name())
}

func (t *AstTestCase) Clone() *AstTestCase {
	inst := new(AstTestCase)
	inst.astTest.Copy(&t.astTest)
	inst.simItems = newAstItems("sim")
	inst.simItems.Cat(t.simItems)
	if t.seeds != nil {
		inst.seeds = make([]int, len(t.seeds))
		copy(inst.seeds, t.seeds)
	}
	return inst
}

func (t *AstTestCase) ParseArgs() {
	t.build = t.GetBuild().Clone()
	//get options sim_options in order
	t.GetOptionArgs().Foreach(func(k string, v interface{}) bool {
		if a, ok := v.(JvsAstOptionForTest); ok {
			a.TestHandler(t)
		}
		if a, ok := v.(JvsAstOptionForBuild); ok {
			a.BuildHandler(t.build)
		}
		return false
	})
}

func (t *AstTestCase) GetTestCases() []*AstTestCase {
	if t.seeds == nil {
		t.seeds = make([]int, 1)
		t.seeds[0] = jvsRand.Intn(math.MaxInt32)
	}
	testcases := make([]*AstTestCase, len(t.seeds))
	for i := range testcases {
		testcases[i] = newAstTestCase(t.GetName() + "__" + strconv.Itoa(t.seeds[i]))
		//copy sim_options and set seed
		testcases[i].simItems.Cat(t.GetBuild().simItems)
		testcases[i].simItems.Cat(t.simItems)
		testcases[i].simItems.option.Cat(newAstItem(GetSimulator().SeedOption() + strconv.Itoa(t.seeds[i])))
	}
	return testcases
}

func (t *AstTestCase) Link() *errors.JVSAstError {
	if err := t.astTest.Link(); err != nil {
		return err
	}
	//set build and check test
	if !t.GetBuild().GetTestDiscoverer().IsValidTest(t.Name) {
		return errors.JVSAstLinkError(t.Name, t.Name+" is not valid test of build"+t.build.Name+"\n"+
			"valid tests:\n"+strings.Join(t.build.GetTestDiscoverer().TestList(), "\n"))
	}

	return nil
}

func (t *AstTestCase) GetHierString(space int) string {
	nextSpace := space + 1
	//not sub test
	flattenTest := ""
	if t.GetBuild() != nil {
		t.ParseArgs()
		flattenTest = astHierFmt("Flatten Tests:", nextSpace, func() string {
			s := ""
			for _, test := range t.GetTestCases() {
				s += test.GetHierString(nextSpace)
			}
			return s
		})
	}
	return astHierFmt(t.GetName()+":", space, func() string {
		return t.astTest.GetHierString(nextSpace) +
			t.simItems.GetHierString(nextSpace) +
			astHierFmt("seeds:", nextSpace, func() string {
				return strings.Repeat(" ", nextSpace) + fmt.Sprintln(t.seeds)
			}) +
			astHierFmt("Builds:", nextSpace, func() string {
				if b := t.GetBuild(); b != nil {
					return b.GetHierString(nextSpace)
				}
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + fmt.Sprint(nil))
			}) + flattenTest
	})
}

type AstGroup struct {
	astTest
	linked bool
	Tests  []*AstTestCase
	Groups map[string]*AstGroup
}

func NewAstGroup(name string) *AstGroup {
	inst := new(AstGroup)
	inst.init(name)
	return inst
}

func (t *AstGroup) Clone() *AstGroup {
	inst := new(AstGroup)
	inst.astTest.Copy(&t.astTest)
	inst.linked = t.linked
	if t.Tests != nil {
		inst.Tests = make([]*AstTestCase, len(t.Tests))
		for i, v := range t.Tests {
			inst.Tests[i] = v.Clone()
			inst.Tests[i].parent = inst
		}
	}
	if t.Groups != nil {
		inst.Groups = make(map[string]*AstGroup)
		for k, v := range t.Groups {
			inst.Groups[k] = v.Clone()
			inst.Groups[k].parent = inst
		}
	}

	return inst
}

func (t *AstGroup) ParseArgs() {
	t.build = t.GetBuild().Clone()
	//get options sim_options in order
	t.GetOptionArgs().Foreach(func(k string, v interface{}) bool {
		if a, ok := v.(JvsAstOptionForBuild); ok {
			a.BuildHandler(t.build)
		}
		return false
	})
}

func (t *AstGroup) GetTestCases() []*AstTestCase {
	testcases := make([]*AstTestCase, 0)
	testcases = append(testcases, t.Tests...)
	for _, group := range t.Groups {
		testcases = append(testcases, group.GetTestCases()...)
	}
	return testcases
}

func (t *AstGroup) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	if ok, testKeywords, _ := t.astTest.KeywordsChecker(s); !ok {
		keywords := utils.NewStringMapSet()
		keywords.AddKey("tests", "groups")
		if !CheckKeyWord(s, keywords) {
			return false, utils.StringMapSetUnion(testKeywords, keywords), "Error in group " + t.Name + ":"
		}
	}
	return true, nil, ""
}

func (t *AstGroup) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	if err := t.astTest.Parse(cfg); err != nil {
		return errors.JVSAstParseError("group "+t.Name, err.Msg)
	}

	//AstParse tests
	if err := CfgToAstItemOptional(cfg, "tests", func(item interface{}) *errors.JVSAstError {
		t.Tests = make([]*AstTestCase, 0)
		tests, ok := item.([]interface{})
		if !ok {
			return errors.JVSAstParseError("tests", fmt.Sprintf("tests must be a list of maps! but it is %T", item))
		}
		for _, test := range tests {
			_test, ok := test.(map[interface{}]interface{})
			if !ok || len(_test) > 1 {
				return errors.JVSAstParseError("tests", fmt.Sprintf("tests must be a list of maps! but it is %T", item))
			}
			for name, attr := range _test {
				testCase := newAstTestCase(name.(string))
				if v, ok := attr.(map[interface{}]interface{}); ok {
					if err := AstParse(testCase, v); err != nil {
						return err
					}
					t.Tests = append(t.Tests, testCase)
					break
				}
				if err := AstParse(testCase, make(map[interface{}]interface{})); err != nil {
					return err
				}
				t.Tests = append(t.Tests, testCase)
			}
		}
		return nil
	}); err != nil {
		return errors.JVSAstParseError(err.Item+"of group "+t.Name, err.Msg)
	}

	//AstParse groups
	if err := CfgToAstItemOptional(cfg, "groups", func(item interface{}) *errors.JVSAstError {
		groups, ok := item.([]interface{})
		if !ok {
			return errors.JVSAstParseError("groups", fmt.Sprintf("groups must be a list of string! but it is %T", item))
		}
		t.Groups = make(map[string]*AstGroup)
		for _, name := range groups {
			if _, ok := t.Groups[name.(string)]; ok {
				return errors.JVSAstParseError("group "+t.Name, "sub group "+name.(string)+" is redefined in group "+t.Name+"!")
			}
			t.Groups[name.(string)] = nil
		}
		return nil
	}); err != nil {
		return errors.JVSAstParseError("group "+t.Name, err.Msg)
	}
	return nil
}

func (t *AstGroup) Link() *errors.JVSAstError {

	if err := t.astTest.Link(); err != nil {
		return err
	}
	//dfs link groups
	for name := range t.Groups {
		group := jvsAstRoot.GetGroup(name)
		if !group.linked {
			if group == nil {
				return errors.JVSAstLinkError("group "+t.Name, "sub group "+name+"of group "+t.Name+" is undef!")
			}
			if err := group.Link(); err != nil {
				return err
			}
		}
		t.Groups[name] = group.Clone()
		t.Groups[name].SetParent(t)
	}
	//check loop
	if t.parent != nil {
		for g := t.parent.(*AstGroup); g.parent != nil; g = g.parent.(*AstGroup) {
			if g.Name == t.Name {
				return errors.JVSAstLinkError("group "+t.Name, "Loop include: group "+t.Name+" and group "+g.Name)
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
	t.linked = true
	return nil
}

func (t *AstGroup) GetHierString(space int) string {
	nextSpace := space + 1
	return astHierFmt(t.GetName()+":", space, func() string {
		return t.astTest.GetHierString(nextSpace) +
			astHierFmt("Builds:", nextSpace, func() string {
				return fmt.Sprintln(strings.Repeat(" ", nextSpace) + t.GetBuild().Name)
			}) +
			astHierFmt("Tests:", nextSpace, func() string {
				s := ""
				for i := range t.Tests {
					s += t.Tests[i].GetHierString(nextSpace + 1)
				}
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
	})
}

//------------------------

//Root
//------------------------
type astRoot struct {
	Env     *astEnv
	Options map[string]*AstOption
	Builds  map[string]*AstBuild
	Groups  map[string]*AstGroup
}

func newAstRoot() *astRoot {
	inst := new(astRoot)
	inst.Builds = make(map[string]*AstBuild)
	inst.Groups = make(map[string]*AstGroup)
	inst.Options = make(map[string]*AstOption)
	return inst
}

func (t *astRoot) GetBuild(name string) *AstBuild {
	if build, ok := t.Builds[name]; ok {
		return build
	}
	return nil
}

func (t *astRoot) GetGroup(name string) *AstGroup {
	if group, ok := t.Groups[name]; ok {
		return group
	}
	return nil
}

func (t *astRoot) KeywordsChecker(s string) (bool, *utils.StringMapSet, string) {
	return true, nil, ""
}

func (t *astRoot) Parse(cfg map[interface{}]interface{}) *errors.JVSAstError {
	//parsing Env
	if err := CfgToAstItemOptional(cfg, "env", func(item interface{}) *errors.JVSAstError {
		t.Env = new(astEnv)
		if item != nil {
			return AstParse(t.Env, item.(map[interface{}]interface{}))
		}
		return AstParse(t.Env, make(map[interface{}]interface{}))
	}); err != nil {
		return err
	}
	//use default
	if t.Env == nil {
		t.Env = new(astEnv)
		if err := AstParse(t.Env, make(map[interface{}]interface{})); err != nil {
			return err
		}
	}
	//parsing builds
	if err := CfgToAstItemRequired(cfg, "builds", func(item interface{}) *errors.JVSAstError {
		for name, build := range item.(map[interface{}]interface{}) {
			t.Builds[name.(string)] = newAstBuild(name.(string))
			if build != nil {
				if err := AstParse(t.Builds[name.(string)], build.(map[interface{}]interface{})); err != nil {
					return err
				}
			} else {
				if err := AstParse(t.Builds[name.(string)], make(map[interface{}]interface{})); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	//parsing options
	if err := CfgToAstItemOptional(cfg, "options", func(item interface{}) *errors.JVSAstError {
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
	if err := CfgToAstItemOptional(cfg, "groups", func(item interface{}) *errors.JVSAstError {
		for name, group := range item.(map[interface{}]interface{}) {
			t.Groups[name.(string)] = NewAstGroup(name.(string))
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

func (t *astRoot) Link() *errors.JVSAstError {
	//dfs link groups
	for _, group := range t.Groups {
		if !group.linked {
			if err := group.Link(); err != nil {
				return err
			}
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
