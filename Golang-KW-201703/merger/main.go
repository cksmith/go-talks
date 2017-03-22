package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"math/rand"
	"sync"
	"time"
)

// Wrap the quote in a Message type with an error field
type Message struct {
	Quote     string
	LastError error
}

func (msg Message) Send(ctx *helper.Context, out chan<- Message) {
	select {
	case out <- msg:
	case <-ctx.Done():
	}
}

func (msg Message) IsTest() bool {
	return len(msg.Quote) == 0
}

// Options for IterateMessages below. Better than adding bool arguments.
type IterateOptions struct {
	ProcessFailingMessages bool
	DontSendMessages       bool
}

type IterateFunc func(msg *Message) bool

func IterateMessages(ctx *helper.Context, in <-chan Message, out chan<- Message, options IterateOptions,
	f IterateFunc) bool {

	for msg := range in {
		if !msg.IsTest() && (options.ProcessFailingMessages || msg.LastError == nil) {
			if f(&msg) {
				if options.DontSendMessages {
					continue
				}
			} else {
				return false
			}
		} else {
			if msg.LastError != nil {
				fmt.Println("Skipping block function due to past error")
			} else {
				fmt.Println("Passing test message to next block")
			}
		}
		msg.Send(ctx, out)
	}
	return true
}

func quoter(ctx *helper.Context, in <-chan Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down quoter")
		defer close(out) // Close the output channel on exit
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
			select {
			case out <- Message{Quote: fmt.Sprintf("%s %d", s, i)}:
			case <-ctx.Done():
				return true
			case msg, running := <-in:
				if !running {
					return true
				}
				msg.Send(ctx, out)
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		return true
	})
	return out
}

func filter(ctx *helper.Context, in <-chan Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down filter")
		defer close(out)
		return IterateMessages(ctx, in, out, IterateOptions{}, func(msg *Message) bool {
			if bad, err := helper.ContainsInappropriateLanguage(msg.Quote); err == nil {
				if bad {
					msg.LastError = fmt.Errorf("Inappropriate quote")
				}
				return true
			} else {
				return false
			}
		})
	})
	return out
}

func printer(ctx *helper.Context, in <-chan Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down printer")
		defer close(out)
		opts := IterateOptions{DontSendMessages: true, ProcessFailingMessages: true}
		return IterateMessages(ctx, in, out, opts, func(msg *Message) bool {
			if msg.LastError == nil {
				if len(msg.Quote) > 500 {
					fmt.Println("Message is too long! Aborting.")
					return false
				} else {
					fmt.Println(msg.Quote)
				}
			} else {
				fmt.Println("Skipping message. An error occurred:", msg.LastError)
			}
			return true
		})
	})
	return out
}

// START1 OMIT

func merger(ctx *helper.Context, inputs ...<-chan Message) <-chan Message {
	var wg sync.WaitGroup
	out := make(chan Message)

	// Start an output goroutine for each input channel.
	output := func(in <-chan Message) {
		defer wg.Done()
		for msg := range in {
			msg.Send(ctx, out)
		}
	}
	wg.Add(len(inputs))
	for _, in := range inputs {
		go output(in)
	}

	// Close the output once all the output goroutines have exited.
	ctx.Run(func() bool {
		wg.Wait()
		close(out)
		return true
	})
	return out
}

// STOP1 OMIT

// START2 OMIT

const NumberOfFilters = 5

func main() {
	ctx := helper.NewContext()
	testIn := make(chan Message)
	c := quoter(ctx, testIn)
	filterOuts := make([]<-chan Message, NumberOfFilters) // HL
	for i := 0; i < NumberOfFilters; i++ {                // HL
		filterOuts[i] = filter(ctx, c) // HL
	} // HL
	c = merger(ctx, filterOuts...) // HL
	printer(ctx, c)

	time.Sleep(5 * time.Second)
	close(testIn) // Shut down normally
	ctx.Wait()
}

// STOP2 OMIT
