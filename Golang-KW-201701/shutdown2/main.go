package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/quote"
	"math/rand"
	"time"
)

func quoter(quit <-chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			s, _ := quote.Get()
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

// START OMIT

func main() {
	quit := make(chan bool) // HL
	c1 := quoter(quit)      // HL
	c2 := quoter(quit)      // HL
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
	close(quit) // HL
}

// STOP OMIT
