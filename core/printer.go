package core

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const printerPadding = 80

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

type stdout struct {
}

func (s *stdout) Write(p []byte) (int, error) {
	n := 0
	paddingCnt := 0
	padding := make([]byte, 0)
	for i, v := range p {
		//"\n"
		if v == 10 {
			if i-n < printerPadding {
				blank := make([]byte, printerPadding-i+n)
				for j := range blank {
					//" "
					blank[j] = 32
				}
				padding = append(padding, append(append(p[n:i], blank...), v)...)
				paddingCnt += len(blank)
			} else {
				padding = append(padding, p[n:i+1]...)
			}
			n = i
		}
		//last one
		if i == len(p)-1 && n != i {
			padding = append(padding, p[n:i+1]...)
		}

	}
	n, err := os.Stdout.Write([]byte(padding))
	return n - paddingCnt, err
}

func Print(s string) {
	if len(s) < printerPadding {
		fmt.Print(s + strings.Repeat(" ", printerPadding-len(s)))
		return
	}
	fmt.Print(s)
}

func Println(s string) {
	if len(s) < printerPadding {
		fmt.Println(s + strings.Repeat(" ", printerPadding-len(s)))
		return
	}
	fmt.Println(s)
}

func PrintStatus(s, status string) {
	Println(s + "..." + status)
}
