package _2phillo_go

import (
	"testing"
	"time"
)

func philoDieIn10Sec(philoNum, dieTime, eatTime, sleepTime, mustEat int) bool {
	isDie := make(chan bool)
	go func(isDie chan bool) {
		time.Sleep(time.Second * 10)
		isDie <- false
	}(isDie)
	philoDo(philoNum, dieTime, eatTime, sleepTime, mustEat)
	go func(isDie chan bool) {
		isDie <- true
	}(isDie)

	return true
}

func TestPhilo(t *testing.T) {
	// inf loop
	isDie := philoDieIn10Sec(4, 410, 200, 200, 0)

	if isDie {
		t.Errorf("phillosopher is die")
	}
}