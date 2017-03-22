package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"math/rand"
	"time"
)

// START1 OMIT

// Wrap the quote in a Message type with an error field
type Message struct {
	Quote     string
	LastError error
}

// STOP1 OMIT

// START2 OMIT

func (msg Message) Send(ctx *helper.Context, out chan<- Message) {
	select {
	case out <- msg:
	case <-ctx.Done():
	}
}

// STOP2 OMIT

// Options for IterateMessages below. Better than adding bool arguments.
type IterateOptions struct {
	ProcessFailingMessages bool
	DontSendMessages       bool
}

type IterateFunc func(msg *Message) bool

// START3 OMIT

func IterateMessages(ctx *helper.Context, in <-chan Message, out chan<- Message, options IterateOptions,
	f IterateFunc) bool {

	for msg := range in {
		if options.ProcessFailingMessages || msg.LastError == nil {
			if f(&msg) {
				if options.DontSendMessages {
					continue
				}
			} else {
				return false
			}
		} else {
			fmt.Println("Skipping block function due to past error")
		}
		msg.Send(ctx, out)
	}
	return true
}

// STOP3 OMIT

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
			case <-in:
				return true
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		return true
	})
	return out
}

// START4 OMIT

func printer(ctx *helper.Context, in <-chan Message) {
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down printer")
		opts := IterateOptions{DontSendMessages: true, ProcessFailingMessages: true} // HL
		return IterateMessages(ctx, in, nil, opts, func(msg *Message) bool {         // HL
			if msg.LastError == nil { // HL
				if len(msg.Quote) > 200 {
					fmt.Println("Message is too long! Aborting.")
					return false
				}
				fmt.Println(msg.Quote)
			} else { // HL
				fmt.Println("Skipping message. An error occurred:", msg.LastError) // HL
			} // HL
			return true
		}) // HL
	})
}

// STOP4 OMIT

func main() {
	ctx := helper.NewContext()
	in := make(chan Message)
	c := quoter(ctx, in)
	printer(ctx, c)

	time.Sleep(10 * time.Second)
	close(in) // Shut down normally
	ctx.Wait()
}
