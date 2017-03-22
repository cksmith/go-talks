package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"math/rand"
	"sync"
	"time"
)

// START1 OMIT

func quoter(quit <-chan bool, wg *sync.WaitGroup) <-chan string {
	wg.Add(1)
	c := make(chan string)
	go func() {
		defer fmt.Println("Shutting down")
		defer wg.Done()
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

// STOP1 OMIT

// START2 OMIT

func main() {
	quit := make(chan bool)
	wg := new(sync.WaitGroup)
	c1 := quoter(quit, wg)
	c2 := quoter(quit, wg)
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
	wg.Wait()
}

// STOP2 OMIT
