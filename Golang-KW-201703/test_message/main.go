package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"math/rand"
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

// START1 OMIT

func (msg Message) IsTest() bool {
	return len(msg.Quote) == 0
}

// STOP1 OMIT

// Options for IterateMessages below. Better than adding bool arguments.
type IterateOptions struct {
	ProcessFailingMessages bool
	DontSendMessages       bool
}

type IterateFunc func(msg *Message) bool

// START2 OMIT

func IterateMessages(ctx *helper.Context, in <-chan Message, out chan<- Message, options IterateOptions,
	f IterateFunc) bool {

	for msg := range in {
		if !msg.IsTest() && (options.ProcessFailingMessages || msg.LastError == nil) { // HL
			if f(&msg) {
				if options.DontSendMessages {
					continue
				}
			} else {
				return false
			}
		} else {
			if msg.LastError != nil { // HL
				fmt.Println("Skipping block function due to past error")
			} else { // HL
				fmt.Println("Passing test message to next block") // HL
			} // HL
		}
		msg.Send(ctx, out)
	}
	return true
}

// STOP2 OMIT

func quoter(ctx *helper.Context, in <-chan Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down quoter")
		defer close(out) // Close the output channel on exit
		// START3 OMIT
		for i := 0; ; i++ {
			s, _ := helper.GetQuote()
			select {
			case out <- Message{Quote: fmt.Sprintf("%s %d", s, i)}:
			case <-ctx.Done():
				return true
			case msg, running := <-in: // HL
				if !running { // HL
					return true // HL
				} // HL
				msg.Send(ctx, out) // HL
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		// STOP3 OMIT
		return true
	})
	return out
}

// START4 OMIT

func printer(ctx *helper.Context, in <-chan Message) <-chan Message {
	out := make(chan Message) // HL
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down printer")
		defer close(out) // HL
		opts := IterateOptions{DontSendMessages: true, ProcessFailingMessages: true}
		return IterateMessages(ctx, in, out, opts, func(msg *Message) bool { // HL
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
	return out // HL
}

// STOP4 OMIT

// START5 OMIT

func main() {
	ctx := helper.NewContext()
	testIn := make(chan Message) // HL
	c := quoter(ctx, testIn)     // HL
	testOut := printer(ctx, c)   // HL

	time.Sleep(5 * time.Second)

	// STOP5 OMIT

	// START6 OMIT

	// Send and receive a test message
	msg := Message{}
	msg.Send(ctx, testIn)
	select {
	case _, running := <-testOut:
		if running {
			fmt.Println("Test passed")
		} else {
			fmt.Println("Test failed. Pipeline shut down.")
		}
	case <-ctx.Done():
		fmt.Println("Test failed. Pipeline cancelled.")
	case <-time.After(10 * time.Second):
		fmt.Println("Test failed. Test timed out.")
	}

	close(testIn) // Shut down normally
	ctx.Wait()
}

// STOP6 OMIT
