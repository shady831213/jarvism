/*
Just A Really Very Impressive Simulation Manager.

Usage:

	jarvism <command> [arguments]

The commands are:

	init        create a jarvism default project
	run_parse   only parse cfg(jarvism_cfg dir or jarvism_cfg.yaml file)
	run_test    run single test, build name must assigned
	run_group   run group
	run_build   run single build
	show_tests  list tests in corresponding build
	show_groups list all groups
	show_builds list all builds
	show_plugins list all plugins or reporter, simulator, runner, checker, testDiscoverer

Use "jarvsim help <command>" for more information about a command.

plugins:
all runners:
	 host
all testDiscoverers:
	 uvm_test
all simulators:
	 vcs
all checkers:
	 compileChecker
	 testChecker
all reporters:
	 junit

run options:

  -compile_args
    	compiling args pass to simulator (default false)
  -max_job int
    	limit of runtime coroutines, default is unlimited. (default -1)
  -repeat
    	run each testcase repeatly n times (default )
  -reporter
    	add reporter plugin, can apply multi times, default
  -seed
    	run testcase with specific seed
  -sim_args
    	simulation args pass to simulator (default false)
  -sim_only
    	bypass compile and only run simulation, default is false.
  -unique
    	if set jobId(timestamp) will be included in hash, then builds and testcases will have unique name and be in unique dir.default is false.
*/

package main

import (
	"fmt"
	"github.com/shady831213/jarvism/cmd"
	"github.com/shady831213/jarvism/core/utils"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, utils.Red(err.Error()))
		os.Exit(2)
	}
}
