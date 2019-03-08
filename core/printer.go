package core

import (
	"fmt"
	"strings"
	"time"
)

func PrintProccessing(color func(str string, modifier ...interface{}) string) func(string, *string, chan bool) {
	return func(processingString string, status *string, done chan bool) {
	forLable:
		for {
			select {
			case <-done:
				break forLable
			default:
				for _, v := range `-\|/` {
					time.Sleep(100 * time.Millisecond)
					fmt.Print(color(processingString+"...") + (*status) + color(string(v)+"\r"))
				}
			}
		}
		fmt.Println(color(processingString+"...") + (*status) + color("done!"))
	}
}

func Print(s string) {
	fmt.Print(s + strings.Repeat(" ", 30))
}

func Println(s string) {
	fmt.Println(s + strings.Repeat(" ", 30))
}

func PrintStatus(s, status string) {
	Println(s + "..." + status)
}
