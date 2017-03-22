package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"math/rand"
	"time"
)

// START1 OMIT

func quoter(ctx *helper.Context) <-chan string {
	out := make(chan string)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down quoter")
		defer close(out) // Close the output channel on exit // HL
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
			select {
			case out <- fmt.Sprintf("%s %d", s, i):
				// do nothing
			case <-ctx.Done(): // HL
				return true
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		return true
	})
	return out
}

// STOP1 OMIT

// START2 OMIT

func printer(ctx *helper.Context, in <-chan string) {
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down printer")
		for msg := range in { // HL
			fmt.Println(msg)
		}
		return true
	})
}

// STOP2 OMIT

// START3 OMIT

func main() {
	ctx := helper.NewContext()
	c := quoter(ctx)
	printer(ctx, c)

	time.Sleep(5 * time.Second)
	ctx.Stop()
}

// STOP3 OMIT
