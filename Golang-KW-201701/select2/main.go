package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/helper"
	"math/rand"
	"time"
)

func quoter() <-chan string { // Returns a receive-only channel of strings
	c := make(chan string)
	go func() { // Launch the go-routine inside the function
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
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
	timeout := time.After(5 * time.Second) // HL
	for {
		select {
		case v1 := <-c1:
			fmt.Println(v1)
		case v2 := <-c2:
			fmt.Println(v2)
		case <-timeout: // HL
			return // HL
		}
	}
}

// STOP OMIT
