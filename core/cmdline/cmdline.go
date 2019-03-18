package cmdline

import (
	"errors"
	"flag"
	"fmt"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/core/plugin"
	"github.com/shady831213/jarvism/core/runtime"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var parseOnly bool

var build string
var test string
var group string

var showBuilds bool
var showGroups bool
var showTests bool
var showPlugins string

func init() {
	options.GetJvsOptions().BoolVar(&parseOnly, "parse_only", false, "only parse cfg files and return, no build and tests run")

	options.GetJvsOptions().StringVar(&build, "build", "", "assign build name, if -test assigned, this field must be assigned")
	options.GetJvsOptions().StringVar(&test, "test", "", "assign test name, -build must be assigned, -group must not be assigned")
	options.GetJvsOptions().StringVar(&group, "group", "", "assign group name, -test must not be assigned, -build must not be assigned")

	options.GetJvsOptions().BoolVar(&showBuilds, "show_builds", false, "only show all builds,no build and tests run")
	options.GetJvsOptions().BoolVar(&showGroups, "show_groups", false, "only show all groups,no build and tests run")
	options.GetJvsOptions().BoolVar(&showTests, "show_tests", false, "only show all tests of corresponding build, -build must be assigned, no build and tests run")
	options.GetJvsOptions().StringVar(&showPlugins, "show_plugins", "", "only show all plugins, value must be in [all, reporter, simulator, checker, runner, testDiscoverer]")

	flag.Usage = func() {
		options.GetJvsOptions().Usage()
	}
}

func parse() error {
	cfg, err := core.GetCfgFile()
	if err != nil {
		return err
	}
	if err := loader.Load(cfg); err != nil {
		return err
	}
	return nil

}

func Run() error {
	err := parse()
	options.GetJvsOptions().Parse(os.Args[1:])
	if parseOnly || err != nil {
		return err
	}

	//show commands
	if showBuilds {
		fmt.Println("all builds:")
		for _, b := range loader.GetJvsAstRoot().GetAllBuilds() {
			fmt.Println("\t", b)
		}
		return nil
	}

	if showGroups {
		fmt.Println("all groups:")
		for _, g := range loader.GetJvsAstRoot().GetAllGroups() {
			fmt.Println("\t", g)
		}
		return nil
	}

	if showTests {
		if build == "" {
			return errors.New(utils.Red("-show_tests assigned but no -build assigned!"))
		}
		dis := loader.GetJvsAstRoot().GetBuild(build).GetTestDiscoverer()
		fmt.Println("all tests founded by testDiscoverer " + dis.Name() + " of build " + build + ":")
		for _, t := range dis.TestList() {
			fmt.Println("\t", t)
		}
		return nil
	}

	if showPlugins != "" {
		switch showPlugins {
		case plugin.JVSReporterPlugin, plugin.JVSTestDiscovererPlugin, plugin.JVSSimulatorPlugin, plugin.JVSCheckerPlugin, plugin.JVSRunnerPlugin:
			fmt.Println("all " + showPlugins + "s:")
			for _, r := range plugin.ValidPlugins(plugin.JVSPluginType(showPlugins)) {
				fmt.Println("\t", r)
			}
		case "all":
			fmt.Println("all plugins:")
			for _, t := range plugin.ValidPluginTypes() {
				fmt.Println("all " + t + "s:")
				for _, r := range plugin.ValidPlugins(plugin.JVSPluginType(t)) {
					fmt.Println("\t", r)
				}
			}
		default:
			return errors.New(utils.Red(flag.Lookup("show_plugins").Usage))
		}
	}

	//check args
	if test != "" && group != "" {
		return errors.New(utils.Red("-test and -group assigned together!"))
	}
	if build != "" && group != "" {
		return errors.New(utils.Red("-build and -group assigned together!"))
	}
	if test != "" && build == "" {
		return errors.New(utils.Red("-test assigned but no -build assigned!"))
	}

	//run jobs
	args := strings.Split(strings.Join(os.Args[1:], "jarvismCmdlineSep"), "jarvismCmdlineSep-")
	for i := range args {
		args[i] = "-" + strings.Replace(args[i], "jarvismCmdlineSep", " ", -1)
	}
	sc := make(chan os.Signal)
	defer close(sc)
	go func() {
		var stopChan = make(chan os.Signal)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT)
		sc <- <-stopChan
	}()
	//build only
	if test == "" && build != "" {
		return runtime.RunOnlyBuild(build, args, sc)
	}
	//run test
	if test != "" && build != "" {
		return runtime.RunTest(test, build, args, sc)
	}
	//run group
	if group != "" {
		return runtime.RunGroup(group, args, sc)
	}
	return nil
}
