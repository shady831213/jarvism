package show

import (
	"errors"
	"fmt"
	"github.com/shady831213/jarvism/cmd/base"
	"github.com/shady831213/jarvism/core/loader"
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/core/plugin"
	"github.com/shady831213/jarvism/core/utils"
)

var CmdShowArgs = &base.Command{
	UsageLine: "jarvism show_args",
	Short:     "list all available arguments",
}

var CmdShowTests = &base.Command{
	UsageLine: "jarvism show_tests [build_name]",
	Short:     "list tests in corresponding build",
}

var CmdShowBuilds = &base.Command{
	UsageLine: "jarvism show_builds",
	Short:     "list all builds",
}

var CmdShowGroups = &base.Command{
	UsageLine: "jarvism show_groups",
	Short:     "list all groups",
}

var CmdShowPlugins = &base.Command{
	UsageLine: "jarvism show_plugins [reporter | simulator | runner | checker | testDiscoverer]",
	Short:     "list all plugins or reporter, simulator, runner, checker, testDiscoverer",
}

func init() {
	CmdShowArgs.Run = runShowArgs
	CmdShowTests.Run = runShowTests
	CmdShowBuilds.Run = runShowBuilds
	CmdShowGroups.Run = runShowGroups
	CmdShowPlugins.Run = runShowPlugins
	base.Jarvism.AddCommand(CmdShowArgs, CmdShowTests, CmdShowGroups, CmdShowBuilds, CmdShowPlugins)
}

func runShowArgs(cmd *base.Command, args []string) error {
	if err := base.Parse(); err != nil {
		return err
	}
	fmt.Println("all args:")
	options.GetJvsOptions().PrintDefaults()
	return nil
}

func runShowTests(cmd *base.Command, args []string) error {

	if len(args) < 1 || base.IsArg(args[0]) {
		return errors.New(utils.Red("jarvism showtests lost build name"))
	}
	if err := base.Parse(); err != nil {
		return err
	}
	dis := loader.GetJvsAstRoot().GetBuild(args[0]).GetTestDiscoverer()
	fmt.Println("all tests founded by testDiscoverer " + dis.Name() + " of build " + args[0] + ":")
	for _, t := range dis.TestList() {
		fmt.Println("\t", t)
	}
	return nil
}

func runShowBuilds(cmd *base.Command, args []string) error {
	if err := base.Parse(); err != nil {
		return err
	}
	fmt.Println("all builds:")
	for _, b := range loader.GetJvsAstRoot().GetAllBuilds() {
		fmt.Println("\t", b)
	}
	return nil
}

func runShowGroups(cmd *base.Command, args []string) error {
	if err := base.Parse(); err != nil {
		return err
	}
	fmt.Println("all groups:")
	for _, b := range loader.GetJvsAstRoot().GetAllGroups() {
		fmt.Println("\t", b)
	}
	return nil
}

func runShowPlugins(cmd *base.Command, args []string) error {
	if err := base.Parse(); err != nil {
		return err
	}
	argsLen := len(args)
	//list all
	if argsLen < 1 {
		for _, t := range plugin.ValidPluginTypes() {
			fmt.Println("all " + t + "s:")
			for _, r := range plugin.ValidPlugins(plugin.JVSPluginType(t)) {
				fmt.Println("\t", r)
			}
		}
		return nil
	}
	//first 5 args is token, others ignore
	if argsLen > 5 {
		argsLen = 5
	}
	for i := 0; i < argsLen; i++ {
		if args[i] != plugin.JVSReporterPlugin &&
			args[i] != plugin.JVSTestDiscovererPlugin &&
			args[i] != plugin.JVSSimulatorPlugin &&
			args[i] != plugin.JVSRunnerPlugin &&
			args[i] != plugin.JVSCheckerPlugin {
			return errors.New(utils.Red("invalid plugin type " + args[i]))
		}
		fmt.Println("all " + args[i] + "s:")
		for _, r := range plugin.ValidPlugins(plugin.JVSPluginType(args[i])) {
			fmt.Println("\t", r)
		}
	}
	return nil
}
