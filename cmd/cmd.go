/*
jarvism cmd interface:
Valid cmd:
	show_tests
	show_builds
	show_groups
	show_plugins

	run_parse
	run_build
	run_test
	run_group

	init
Run 'jarvsim help <command>' for details.
*/
package cmd

import (
	"flag"
	"fmt"
	"github.com/shady831213/jarvism/cmd/base"
	_ "github.com/shady831213/jarvism/cmd/init"
	_ "github.com/shady831213/jarvism/cmd/run"
	_ "github.com/shady831213/jarvism/cmd/show"
	"os"
	"strings"
)

//cmd interface entry, main() call this function
func Run() error {
	flag.Usage = base.Usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		base.Usage()
	}

	cmdName := args[0]
	if args[0] == "help" {
		help(args[1:])
		return nil
	}

BigCmdLoop:
	for bigCmd := base.Jarvism; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					printUsage(os.Stderr, bigCmd)
					base.SetExitStatus(2)
					base.Exit()
				}
				if args[0] == "help" {
					help(append(strings.Split(cmdName, " "), args[1:]...))
					return nil
				}
				cmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Flag.Usage = func() {
				cmd.Usage()
			}
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			return cmd.Run(cmd, args)
		}
		helpArg := ""
		if i := strings.LastIndex(cmdName, " "); i >= 0 {
			helpArg = " " + cmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "jarvism %s: unknown command\nRun 'jarvism help%s' for usage.\n", cmdName, helpArg)
		base.SetExitStatus(2)
		base.Exit()
	}
}

func init() {
	base.Usage = mainUsage
}

func mainUsage() {
	printUsage(os.Stderr, base.Jarvism)
	os.Exit(2)
}
