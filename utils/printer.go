package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func PrintProccessing(color func(str string, modifier ...interface{}) string) func(s string, done chan bool) {
	return func(s string, done chan bool) {
	forLable:
		for {
			select {
			case <-done:
				break forLable
			default:
				for _, v := range `-\|/` {
					time.Sleep(100 * time.Millisecond)
					fmt.Print(color(s + "..." + string(v) + "\r"))
				}
			}
		}
	}
}

func PrintStatus(colors, colorStatus func(str string, modifier ...interface{}) string) func(s, status string) {
	return func(s, status string) {
		fmt.Println(colors(s+"...") + colorStatus(status) + strings.Repeat(" ", 10))
	}
}

func Green(str string, modifier ...interface{}) string {
	return printer(str, 32, 0, modifier...)
}

func LightGreen(str string, modifier ...interface{}) string {
	return printer(str, 32, 1, modifier...)
}

func Cyan(str string, modifier ...interface{}) string {
	return printer(str, 36, 0, modifier...)
}

func LightCyan(str string, modifier ...interface{}) string {
	return printer(str, 36, 1, modifier...)
}

func Red(str string, modifier ...interface{}) string {
	return printer(str, 31, 0, modifier...)
}

func LightRed(str string, modifier ...interface{}) string {
	return printer(str, 31, 1, modifier...)
}

func Yellow(str string, modifier ...interface{}) string {
	return printer(str, 33, 0, modifier...)
}

func Black(str string, modifier ...interface{}) string {
	return printer(str, 30, 0, modifier...)
}

func DarkGray(str string, modifier ...interface{}) string {
	return printer(str, 30, 1, modifier...)
}

func LightGray(str string, modifier ...interface{}) string {
	return printer(str, 37, 0, modifier...)
}

func White(str string, modifier ...interface{}) string {
	return printer(str, 37, 1, modifier...)
}

func Blue(str string, modifier ...interface{}) string {
	return printer(str, 34, 0, modifier...)
}

func LightBlue(str string, modifier ...interface{}) string {
	return printer(str, 34, 1, modifier...)
}

func Purple(str string, modifier ...interface{}) string {
	return printer(str, 35, 0, modifier...)
}

func LightPurple(str string, modifier ...interface{}) string {
	return printer(str, 35, 1, modifier...)
}

func Brown(str string, modifier ...interface{}) string {
	return printer(str, 33, 0, modifier...)
}

func Blink(colorFunc func(str string, modifier ...interface{}) string) func(str string) string {
	return func(str string) string {
		return colorFunc(str, 1)
	}
}

func printer(str string, color int, weight int, extraArgs ...interface{}) string {
	var isBlink int64 = 0
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	var isUnderLine int64 = 0
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "2")
	}
	if isUnderLine > 0 {
		mo = append(mo, "4")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
}
