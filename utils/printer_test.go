package utils

import (
	"strconv"
	"testing"
	"time"
)

func TestPrinter(t *testing.T) {
	//fmt.Print(Blink(Green)("testing..."))
	done1 := make(chan bool, 1)
	go printProccessing("test1", done1)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			printDone(strconv.Itoa(i))
		}
	}()
	time.Sleep(10 * time.Second)
	done1 <- true
}
