package runtime_test

import (
	"github.com/shady831213/jarvism/core/runtime"
	"github.com/shady831213/jarvism/utils"
	"strconv"
	"testing"
	"time"
)

func TestPrinter(t *testing.T) {
	status := "(cnt:0/total:10)"
	done1 := make(chan bool, 1)
	job := make(chan bool)
	go runtime.PrintProccessing(utils.Green)("test1", &status, done1)
	go func() {
		i := 0
		for {
			<-job
			i++
			status = "(cnt:" + strconv.Itoa(i) + "/total:10)"
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(50 * time.Millisecond)
			runtime.PrintStatus(strconv.Itoa(i), utils.Green("done"))
			job <- true
		}
	}()
	time.Sleep(2 * time.Second)
	done1 <- true
	time.Sleep(time.Second)
}
