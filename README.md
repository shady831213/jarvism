# jarvism[![Build Status](https://travis-ci.org/shady831213/jarvism.svg?branch=master)](https://travis-ci.org/shady831213/jarvism)
Just A Really Very Impressive Simulation Manager

# install
instatll go

[go](https://github.com/golang/go)

install jarvism
```
go get -u github.com/shady831213/jarvism

jarvism help
```

# test
```
cd [JARVISM_DIR]
./script/clean.sh
go test -v ./...
./script/clean.sh
```

# Getting start
```
$ jarvism init -prj_dir jvs_prj
$ cd jvs_prj
$ ls 
algorithms  go	jarivsm  jvs_prj  src  testcases  work
$ jarvism help init
usage: jarvism init [-prj_dir DIR][-work_dir DIR]

. $prj_dir
|--- jarvism_cfg
|------ jarvism_cfg.yaml
|------ jarvism_setup.sh(export $JVS_PRJ_HOME;export $JVS_WORK_DIR)
|--- src
|--- testcases
. $work_dir
  -prj_dir string
    	assign prj dir, default is pwd
  -work_dir string
    	assign work dir, default is $prj_dir/work
```
Enjoy!

# Config
jarvism allows you use a single yaml file ($JVS_PRJ_HOME/jarvism_cfg.yaml) or a banch of yaml files ($JVS_PRJ_HOME/jarvism_cfg/*.yaml) to config project. Refer to https://github.com/shady831213/jarvism/tree/master/core/runtime/testFiles/jarvism_cfg

## env
"env" config simulator and runner, which are all both parsable plugin, they defiend by "type" and "attr".
"type" is required. "attr" is optional, depends on plugin implementation (about plugin, see below).
If "env" is not defined, simulator and runner will be default, "vcs" and "host" respectively.
If runner or simulator not defined in "env", default will beb used.
e.g
```yaml
env: # runner is default
  simulator:
    type: "vcs" 
```


```yaml
env: # simulator is default
  runner:
    type: "host"
    attr: #attributes used by "host" runner
```

If the default simulator or runner can not meet your requirement, you can implement your own plugin, and use them in your config file.(about plugin, see below).


## builds
"builds" is a required tag.  "builds" defined one or multiple build, which used to config a specific compile flow.
e.g
```yaml
common_compile_option: &common_compile >-
  -sverilog
  -ntb_opts uvm-1.2

common_sim_option: &common_sim >-
  +UVM_VERBOSITY=UVM_LOW
  +UVM_CONFIG_DB_TRACE

builds:
  build1:
    compile_checker:
      type:
        compileChecker
      attr:
        fail:
          - .*Error.*
    compile_option:
      - *common_compile
      - -timescale=1ns/10ps
    pre_sim_action:
      - echo "pre_sim_build1"
    sim_option:
      - *common_sim
    post_sim_action:
      - echo "post_sim_build1"

  build2:
    pre_compile_action:
      - echo "pre_compile_build2"
    compile_option:
      - -debug_access+pp
      - *common_compile
    post_compile_action:
      - echo "post_compile_build2"
    sim_option:
      - *common_sim
```

As about example, "build1" and "build2" are build name. There are some attributes in a build:

+ compile_option: A list define compile options will pass to simulator compile flow

+ sim_option: A list define simulation options will pass to simulator runtime 

+ compile_checker: A parsable plugin. If it is not defined, default compile_checker "compileChecker" will be used.
		   You can add error, warning, unknown and exclued error, warning, unknown pattern through "attr".
		   If default compile_checker can not meet your requirement, you can write your own checker plugin.
+ sim_checker: Same as compile_checker, except it is use for simulation, and some pattern has been pre defined.

+ *_action: There are 4 hooks provided. You can add some cmd sequences before or after compile and simulation.

+ test_discoverer: A parsable plugin. If it is not defined, default test_discoverer "uvm_test" will be used.
		   You can define top testcases dir through attr, defualt is $JVS_PRJ_HOME/testcases.
		   If your testcases are compliance with following convention, they will be discovered automatically.
		   And you don't need add them to file list.Besides, the testcase dir name will pass to simulator through
		   +UVM_TESTNAME.
		   
		   - all testcases in testcases dir
		   
		   - testcase dir name is same as .sv file, e.g
		   
```
		   	. testcases/
			--------test1/ \\valid
			----------------test1.sv	
			--------test2/ \\invalid
			----------------test3.sv	
			--------test3/ \\invalid
			----------------test3.c	
			--------test4/ \\invalid

```
		   
		   - uvm_test name is same as .sv file and testcase dir name
		   
		   If default test_discoverer can not meet your requirement, you can write your own test_discoverer plugin.
		   
## groups
"groups" defined one or multiple group, which used to organize testcases.
e.g
```yaml
groups:
  group1:
    build: build1
    args:
      - -vh, -repeat 1
    tests:
      - test1:
          args:
            - -repeat 10
      - test2:
          args:
            - -seed 1

  group2:
    build: build2
    args:
      - -vh
      - -repeat 1
    tests:
      - test3:
          args:
            - -repeat 10
    groups:
      - group1

  group3:
    build: build2
    args:
      - -vh
      - -repeat 20
    tests:
      - test1:
    groups:
      - group2
      - group1
```

As about example, "group1", "group2" and "group3" are group name. There are some attributes in a group:

+ build: Assign a defined build name to this group. This build will be used for all tests and subgroups in this group, if they don't define their own build.

+ args: A list define pre-defined and user-defined(about user-defined options, see below) arguments. These arguments will be used for all tests and subgroups in this group, if they don't override them. For expample, in test3 of group2, -repeat value will be 10. Multiple args in one line is allowed, but they must be seperated by ",".

+ tests: A list define testcases in the group. Each test can config it's own build and args. Build and args defined more nested have higher priority.

+ groups: A list of defined sub groups.

If some testcases in the same group tree use the same build with the same compile_option and pre/post_compile_action, jarvism can detected and try to let them share the same compile database.

# Example

TBD

But if you have vcs, you can run tests in plugins/runner/host
```
cd $GOPATH/src/github.com/shady831213/jarvism/plugins/runner/host
go test -v
```

# Plugins

TBD

# godoc
[doc](https://godoc.org/github.com/shady831213/jarvism)

# usage
```
Just A Really Very Impressive Simulation Manager.

-config_top assign a config top inside $JVS_PRJ_HOME/jarvism_cfg, default is $JVS_PRJ_HOME/jarvism_cfg or $JVS_PRJ_HOME/jarvism_cfg.yaml. 
If only jarsivm_cfg.yaml in $JVS_PRJ_HOME,or config_top is not existed in $JVS_PRJ_HOME/jarvism_cfg, this argument will be ignored.

Usage:

	jarvism [-config_top config_top] <command> [arguments]

The commands are:

	init        create a jarvism default project
	run_parse   only parse cfg(jarvism_cfg dir or jarvism_cfg.yaml file)
	run_test    run single test, build name must assigned
	run_group   run group
	run_build   run single build
	show_args   list all available arguments
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

all args:
  -compile_args
    	compiling args pass to simulator (default false)
  -max_job int
    	limit of runtime coroutines, default is unlimited. (default -1)
  -quite_comp
    	quite compiling with -q, and close lint with +lint=none (default false)
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
  -test_phase
    	user-defined flag (default false)
  -unique
    	if set jobId(timestamp) will be included in hash, then builds and testcases will have unique name and be in unique dir.default is false.
  -vh
    	user-defined flag (default false)
  -wave
    	dump waveform, vaule is format[FSDB, VPD], use macro[DUMP_FSDB, DUMP_VPD] in your testbench, default is VPD (default false)

```
