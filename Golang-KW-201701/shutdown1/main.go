package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/helper"
	"math/rand"
	"time"
)

// START1 OMIT

func quoter(quit <-chan bool) <-chan string { // HL
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
			select {
			case c <- fmt.Sprintf("%s %d", s, i):
				// do nothing
			case <-quit: // HL
				return // HL
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c
}

// STOP1 OMIT

// START2 OMIT

func main() {
	quit1 := make(chan bool) // HL
	c1 := quoter(quit1)      // HL
	quit2 := make(chan bool) // HL
	c2 := quoter(quit2)      // HL
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
	quit1 <- true // HL
	quit2 <- true // HL
}

// STOP2 OMIT
