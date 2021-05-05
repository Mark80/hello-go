package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct {
	lock sync.Mutex
	value int
}

func (counter *SafeCounter) Inc() {
	counter.lock.Lock()
	counter.value += 1
	counter.lock.Unlock()
}

func (counter *SafeCounter) Get() int {
	counter.lock.Lock()
	defer counter.lock.Unlock()
	return counter.value
}


func main() {
	counter := SafeCounter{value: 0}

	for i := 0; i < 1000; i++ {
		go counter.Inc()
	}

	time.Sleep(time.Second)

	fmt.Println(counter.value)

}
