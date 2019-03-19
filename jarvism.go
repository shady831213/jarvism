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
