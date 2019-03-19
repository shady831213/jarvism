package run

import (
	"errors"
	"github.com/shady831213/jarvism/cmd/base"
	"github.com/shady831213/jarvism/core/options"
	"github.com/shady831213/jarvism/core/runtime"
	"github.com/shady831213/jarvism/core/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var CmdRunParse = &base.Command{
	UsageLine:   "jarvism run_parse",
	Short:       "only parse cfg(jarvism_cfg dir or jarvism_cfg.yaml file)",
	CustomFlags: true,
}

var CmdRunTest = &base.Command{
	UsageLine:   "jarvism run_test [build_name][test_name][args]",
	Short:       "run single test, build name must assigned",
	Flag:        *options.GetJvsOptions(),
	CustomFlags: true,
}

var CmdRunBuild = &base.Command{
	UsageLine:   "jarvism run_build [build_name][args]",
	Short:       "run single build",
	Flag:        *options.GetJvsOptions(),
	CustomFlags: true,
}

var CmdRunGroup = &base.Command{
	UsageLine:   "jarvism run_group [group_name][args]",
	Short:       "run group",
	Flag:        *options.GetJvsOptions(),
	CustomFlags: true,
}

func init() {
	CmdRunParse.Run = runRunParse
	CmdRunTest.Run = runRunTest
	CmdRunBuild.Run = runRunBuild
	CmdRunGroup.Run = runRunGroup
	base.Jarvism.AddCommand(CmdRunParse, CmdRunTest, CmdRunGroup, CmdRunBuild)
}

func formatArgs(args []string) []string {
	fArgs := strings.Split(strings.Join(args, "jarvismCmdlineSep"), "jarvismCmdlineSep-")
	for i := range fArgs {
		fArgs[i] = strings.Replace(fArgs[i], "jarvismCmdlineSep", " ", -1)
		if fArgs[i][0] != '-' {
			fArgs[i] = "-" + fArgs[i]
		}
	}
	return fArgs
}

func catSignal(sc chan os.Signal) {
	stopChan := make(chan os.Signal)
	defer close(stopChan)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT)
	sc <- <-stopChan
}

func runRunParse(cmd *base.Command, args []string) error {
	return base.Parse()
}

func runRunTest(cmd *base.Command, args []string) error {
	if len(args) < 2 || base.IsArg(args[0]) || base.IsArg(args[1]) || base.IsHelp(args[0]) {
		cmd.Flag.Usage()
		return errors.New(utils.Red("jarvism run_test must assign build_name and test_name"))
	}
	if err := base.Parse(); err != nil {
		return err
	}
	var runArgs []string
	if len(args) > 2 {
		runArgs = formatArgs(args[2:])
	}
	sc := make(chan os.Signal)
	defer close(sc)
	go catSignal(sc)
	return runtime.RunTest(args[1], args[0], runArgs, sc)
}

func runRunBuild(cmd *base.Command, args []string) error {
	if len(args) < 1 || base.IsArg(args[0]) || base.IsHelp(args[0]) {
		cmd.Flag.Usage()
		return errors.New(utils.Red("jarvism run_build must assign build_name"))
	}
	if err := base.Parse(); err != nil {
		return err
	}
	var runArgs []string
	if len(args) > 1 {
		runArgs = formatArgs(args[1:])
	}
	sc := make(chan os.Signal)
	defer close(sc)
	go catSignal(sc)
	return runtime.RunOnlyBuild(args[0], runArgs, sc)
}

func runRunGroup(cmd *base.Command, args []string) error {
	if len(args) < 1 || base.IsArg(args[0]) || base.IsHelp(args[0]) {
		cmd.Flag.Usage()
		return errors.New(utils.Red("jarvism run_group must assign group_name"))
	}
	if err := base.Parse(); err != nil {
		return err
	}
	var runArgs []string
	if len(args) > 1 {
		runArgs = formatArgs(args[1:])
	}
	sc := make(chan os.Signal)
	defer close(sc)
	go catSignal(sc)
	return runtime.RunGroup(args[0], runArgs, sc)
}
