package main

import (
	"fmt"
	"sync"
)

type safeCounter struct {
	mutex sync.Mutex
	value int
}

func (counter *safeCounter) increments() {
	counter.mutex.Lock()
	counter.value++
	fmt.Println(counter.value)
	counter.mutex.Unlock()

}
func (counter *safeCounter) get() {
	fmt.Println(counter.value)
}

func slowInc(c *safeCounter) {

	for i := 0; i < 10; i++ {
		c.increments()

	}

}

func fastInc(c *safeCounter) {

	for i := 0; i < 10; i++ {
		c.increments()
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	counter := safeCounter{}
	go func() {
		slowInc(&counter)
		wg.Done()
	}()
	go func() {
		fastInc(&counter)
		wg.Done()
	}()
	wg.Wait()

}
