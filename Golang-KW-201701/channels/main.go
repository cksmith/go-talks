package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/quote"
	"math/rand"
	"time"
)

// START1 OMIT

// START2 OMIT
func quoter() <-chan string { // Returns a receive-only channel of strings
	// STOP2 OMIT
	c := make(chan string) // HL
	go func() {            // Launch the go-routine inside the function
		for i := 0; ; i++ {
			s, _ := quote.Get()
			c <- fmt.Sprintf("%s %d", s, i) // HL
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller
}

func main() {
	c := quoter()
	for i := 0; i < 5; i++ {
		fmt.Println(<-c) // HL
	}
}

// STOP1 OMIT
