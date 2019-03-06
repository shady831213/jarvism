package utils

import (
	"strconv"
	"testing"
	"time"
)

func TestPrinter(t *testing.T) {
	done1 := make(chan bool, 1)
	go PrintProccessing(Green)("test1", done1)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(200 * time.Millisecond)
			PrintStatus(Black, Green)(strconv.Itoa(i), "done")
		}
	}()
	time.Sleep(3 * time.Second)
	done1 <- true
	time.Sleep(time.Second)
}
