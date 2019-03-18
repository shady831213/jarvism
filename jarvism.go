package main

import (
	"fmt"
	"github.com/shady831213/jarvism/core/cmdline"
	"github.com/shady831213/jarvism/core/utils"
	"os"
)

func main() {
	if err := cmdline.Run(); err != nil {
		fmt.Println(utils.Red(err.Error()))
		os.Exit(2)
	}
}
