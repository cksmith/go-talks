package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"time"
)

// Wrap the quote in a Message type with an error field

// START1 OMIT

type Message struct {
	Id        uint64
	LastError error
}

// STOP1 OMIT

func (msg Message) Send(ctx *helper.Context, out chan<- Message) {
	select {
	case out <- msg:
	case <-ctx.Done():
	}
}

func (msg Message) IsTest() bool {
	return msg.Id == 0
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

// START2 OMIT
func subscribe(ctx *helper.Context, in <-chan Message) (<-chan Message, chan<- Message) {
	nextId := uint64(1) // HL
	out := make(chan Message)
	responseChannel := make(chan Message)          // HL
	ticker := time.NewTicker(1 * time.Millisecond) // HL
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down subscribe")
		defer ticker.Stop() // HL
		defer close(out)
		// STOP2 OMIT

		// START3 OMIT
		for {
			select {
			case _ = <-ticker.C:
				out <- Message{Id: nextId}
				nextId++
			case <-ctx.Done():
				return true
			case msg, running := <-in:
				if !running {
					fmt.Println("Test input channel closed")
					return true
				}
				msg.Send(ctx, out)
			case msg, running := <-responseChannel:
				if !running {
					fmt.Println("Response channel closed")
				}
				fmt.Println("Response received for id", msg.Id)
			}
		}
		// STOP3 OMIT
		return true
	})
	return out, responseChannel
}

func process(ctx *helper.Context, in <-chan Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down process")
		defer close(out)
		return IterateMessages(ctx, in, out, IterateOptions{}, func(msg *Message) bool {
			return true
		})
	})
	return out
}

// START4 OMIT

func respond(ctx *helper.Context, in <-chan Message, responseChannel chan<- Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down respond")
		defer close(out)
		defer close(responseChannel) // HL
		opts := IterateOptions{DontSendMessages: true, ProcessFailingMessages: true}
		return IterateMessages(ctx, in, out, opts, func(msg *Message) bool {
			msg.Send(ctx, responseChannel) // HL
			return true
		})
	})
	return out
}

// STOP4 OMIT

// START5 OMIT

func main() {
	ctx := helper.NewContext()
	testIn := make(chan Message)
	c, responseChannel := subscribe(ctx, testIn) // HL
	c = process(ctx, c)
	c = process(ctx, c)
	c = process(ctx, c)
	c = process(ctx, c)
	testOut := respond(ctx, c, responseChannel) // HL

	time.Sleep(100 * time.Millisecond)

	// Send and receive a test message
	// ...

	// STOP5 OMIT

	fmt.Println("Sending test message")
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
	case <-time.After(1 * time.Second):
		fmt.Println("Test failed. Test timed out.")
	}

	close(testIn) // Shut down normally
	ctx.Wait()
}
