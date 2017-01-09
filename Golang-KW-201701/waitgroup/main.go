package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/helper"
	"math/rand"
	"sync"
	"time"
)

// START1 OMIT

func quoter(quit <-chan bool, wg *sync.WaitGroup) <-chan string { // HL
	wg.Add(1) // Increment the WaitGroup counter // HL
	c := make(chan string)
	go func() {
		defer fmt.Println("Shutting down") // HL
		defer wg.Done()                    // Decrement the WaitGroup counter on exit // HL
		// STOP1 OMIT
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
			select {
			case c <- fmt.Sprintf("%s %d", s, i):
				// do nothing
			case <-quit:
				return
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c
}

// START2 OMIT

func main() {
	quit := make(chan bool)
	wg := new(sync.WaitGroup) // HL
	c1 := quoter(quit, wg)    // HL
	c2 := quoter(quit, wg)    // HL
	timeout := time.After(5 * time.Second)
ForLoop:
	for {
		select {
		case v1 := <-c1:
			fmt.Println(v1)
		case v2 := <-c2:
			fmt.Println(v2)
		case <-timeout:
			break ForLoop
		}
	}
	close(quit)
	wg.Wait() // Wait for all goroutines to complete // HL
}

// STOP2 OMIT
