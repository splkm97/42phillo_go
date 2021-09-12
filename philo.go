package _2phillo_go

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutexes []*sync.Mutex
	philos []philosopher
	forks []bool
	dieChan chan int
)

type philosopher struct {
	forkCnt int
	timeCnt int
	eatCnt int
	sleepCnt int
	printTiming chan string
}

func initPhilo(philoNum int) {
	forks = make([]bool, philoNum)
	dieChan = make(chan int, philoNum)
	mutexes = make([]*sync.Mutex, 0, philoNum)
	philos = make([]philosopher, 0, philoNum)
	for i := 0 ; i < philoNum ; i ++ {
		forks[i] = false
		philos = append(philos, philosopher{0, 0, 0, 0, make(chan string)})
		mutexes = append(mutexes, &sync.Mutex{})
	}
}

func isEatable(id, nextId int) bool {
	mutexes[id].Lock()
	if !forks[id] {
		forks[id] = true
		philos[id].forkCnt ++
	}
	mutexes[id].Unlock()
	mutexes[nextId].Lock()
	if !forks[nextId] {
		forks[nextId] = true
		philos[id].forkCnt ++
	}
	mutexes[nextId].Unlock()
	switch philos[id].forkCnt {
	case 2:
		return true
	default:
		return false
	}
}

func philoTask(id, philoNum, dieTime, eatTime, sleepTime, mustEat int) {
	cntTime, sleepCnt, eatCnt := initTimeCnt(id)
	eatFlag := false
	for {
		nextId := id + 1
		if nextId == philoNum {
			nextId = 0
		}
		cntTime++
		switch eatFlag {
		case true:
			if eatCnt == eatTime {
				eatFlag = !eatFlag
				philos[id].forkCnt -= 2
				cntTime -= dieTime
				eatCnt = 0
			}
			eatCnt ++
		case false:
			if sleepCnt == sleepTime {
				if isEatable(id, nextId) {
					eatFlag = !eatFlag
				}
				sleepCnt = 0
			}
			sleepCnt ++
		}
		if dieTime == cntTime {
			dieChan <- id
			return
		}
	}
	updateTimeCnt(id, cntTime, eatCnt, sleepCnt)
}

func initTimeCnt(id int) (int, int, int) {
	cntTime := philos[id].timeCnt
	sleepCnt := philos[id].sleepCnt
	eatCnt := philos[id].eatCnt
	return cntTime, sleepCnt, eatCnt
}

func updateTimeCnt(id int, cntTime int, eatCnt int, sleepCnt int) {
	philos[id].timeCnt = cntTime
	philos[id].eatCnt = eatCnt
	philos[id].sleepCnt = sleepCnt
}

func printState(philoNum int) {
	for {
		for i := 0; i < philoNum; i++ {
			str := <-philos[i].printTiming
			fmt.Println(str)
		}
	}
}

func philoDo(philoNum, dieTime, eatTime, sleepTime, mustEat int) int {
	initPhilo(philoNum)
	for i := 0 ; i < philoNum ; i ++ {
		go philoTask(i, philoNum, dieTime, eatTime, sleepTime, mustEat)
	}
	time.Sleep(1 * time.Second)
	go printState(philoNum)
	diePid := <-dieChan
	dieMsg := <-philos[diePid].printTiming
	fmt.Println(dieMsg)
	return 0
}