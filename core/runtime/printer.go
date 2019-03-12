package runtime

import (
	"fmt"
	"github.com/shady831213/jarvism/core/ast"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var runtimeLog *log.Logger

const printerPadding = 100

func setLog(logFileName string) (*os.File, error) {
	if err := os.MkdirAll(path.Join(ast.GetWorkDir(), "JarvismLog"), os.ModePerm); err != nil {
		return nil, err
	}
	logFile, err := os.Create(path.Join(ast.GetWorkDir(), "JarvismLog", logFileName))
	if err != nil {
		return logFile, err
	}
	runtimeLog = log.New(logFile, "Jarvism_", log.LstdFlags)
	return logFile, nil
}

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
		result := color(processingString+"...") + (*status) + color("done!")
		fmt.Println(result)
		runtimeLog.Println(result)
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
				padding = append(padding, blank...)
				paddingCnt += len(blank)
			}
			n = i
		}
		padding = append(padding, v)
	}
	n, err := os.Stdout.Write([]byte(padding))
	return n - paddingCnt, err
}

func paddingString(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		if len(lines[i]) < printerPadding {
			lines[i] += strings.Repeat(" ", printerPadding-len(lines[i]))
		}
	}
	return strings.Join(lines, "\n")
}

func Print(s string) {
	ps := paddingString(s)
	fmt.Print(ps)
	runtimeLog.Println(ps)
}

func Println(s string) {
	ps := paddingString(s)
	fmt.Println(ps)
	runtimeLog.Println(ps)
}

func PrintStatus(s, status string) {
	Println(s + "..." + status)
}
