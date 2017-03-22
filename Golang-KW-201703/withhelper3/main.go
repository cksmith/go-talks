package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"math/rand"
	"time"
)

func quoter(ctx *helper.Context, in <-chan struct{}) <-chan string {
	out := make(chan string)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down quoter")
		defer close(out) // Close the output channel on exit
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
			select {
			case out <- fmt.Sprintf("%s %d", s, i):
				// do nothing
			case <-ctx.Done():
				return true
			case <-in:
				return true
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		return true
	})
	return out
}

// START OMIT

func printer(ctx *helper.Context, in <-chan string) {
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down printer")
		for msg := range in {
			if len(msg) > 200 { // HL
				fmt.Println("Message is too long! Aborting.") // HL
				return false                                  // HL
			} // HL
			fmt.Println(msg)
		}
		return true
	})
}

// STOP OMIT

func main() {
	ctx := helper.NewContext()
	in := make(chan struct{})
	c := quoter(ctx, in)
	printer(ctx, c)

	time.Sleep(10 * time.Second)
	close(in)
	ctx.Wait()
}
