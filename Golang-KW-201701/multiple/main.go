package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/quote"
	"math/rand"
	"time"
)

func quoter() <-chan string { // Returns a receive-only channel of strings
	c := make(chan string)
	go func() { // Launch the go-routine inside the function
		for i := 0; ; i++ {
			s, _ := quote.Get()
			c <- fmt.Sprintf("%s %d", s, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller
}

// START OMIT

func main() {
	c1 := quoter()
	c2 := quoter()
	go func() {
		for {
			fmt.Println(<-c1)
		}
	}()
	go func() {
		for {
			fmt.Println(<-c2)
		}
	}()
	time.Sleep(5 * time.Second)
}

// STOP OMIT
