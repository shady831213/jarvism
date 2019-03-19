package base

import (
	"flag"
	"fmt"
	"github.com/shady831213/jarvism/core"
	"github.com/shady831213/jarvism/core/loader"
	"log"
	"os"
	"strings"
	"sync"
)

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) error

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'jarvsim help' output.
	Short string

	// Long is the long message shown in the 'jarvsim help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool

	// Commands lists the available commands and help topics.
	// The order here is the order in which they are printed by 'jarvsim help'.
	// Note that subcommands are in general best avoided.
	Commands []*Command
}

var Jarvism = &Command{
	UsageLine: "jarvism",
	Long:      `Just A Really Very Impressive Simulation Manager.`,
	// Commands initialized in package main
}

// LongName returns the command's long name: all the words in the usage line between "jarvsim" and a flag or argument,
func (c *Command) LongName() string {
	name := c.UsageLine
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}
	if name == "jarvism" {
		return ""
	}
	return strings.TrimPrefix(name, "jarvism ")
}

// Name returns the command's short name: the last word in the usage line before a flag or argument.
func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
	c.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Run 'jarvsim help %s ' for details.\n", c.LongName())
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var atExitFuncs []func()

func AtExit(f func()) {
	atExitFuncs = append(atExitFuncs, f)
}

func Exit() {
	for _, f := range atExitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

func Fatalf(format string, args ...interface{}) {
	Errorf(format, args...)
	Exit()
}

func Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	SetExitStatus(1)
}

func ExitIfErrors() {
	if exitStatus != 0 {
		Exit()
	}
}

var exitStatus = 0
var exitMu sync.Mutex

func SetExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func (c *Command) AddCommand(subCommands ...*Command) {
	if c.Commands == nil {
		c.Commands = make([]*Command, 0)
	}
	c.Commands = append(c.Commands, subCommands...)
}

// Usage is the usage-reporting function, filled in by package main
// but here for reference by other packages.
var Usage func()

func Parse() error {
	cfg, err := core.GetCfgFile()
	if err != nil {
		return err
	}
	if err := loader.Load(cfg); err != nil {
		return err
	}
	return nil

}

func IsArg(arg string) bool {
	return arg[0] == '-'
}

func IsHelp(arg string) bool {
	return arg == "help" || arg == "h" || arg == "-help" || arg == "-h"
}
