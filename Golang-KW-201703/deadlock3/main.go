package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"time"
)

// Wrap the quote in a Message type with an error field
type Message struct {
	Id        uint64
	LastError error
}

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

func subscribe(ctx *helper.Context, in <-chan Message) (<-chan Message, chan<- Message) {
	nextId := uint64(1)
	out := make(chan Message)
	responseChannel := make(chan Message)
	ticker := time.NewTicker(1 * time.Nanosecond)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down subscribe")
		defer ticker.Stop()
		defer close(out)
		// START1 OMIT
		for {
			select {
			case _ = <-ticker.C:
				select {
				case out <- Message{Id: nextId}:
				case <-ctx.Done():
					close(out)
					return true
				case msg, running := <-responseChannel:
					if running {
						fmt.Println("Response received for id", msg.Id)
					}
				}
				nextId++
				// STOP1 OMIT
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

func respond(ctx *helper.Context, in <-chan Message, responseChannel chan<- Message) <-chan Message {
	out := make(chan Message)
	ctx.Run(func() bool {
		defer fmt.Println("Shutting down respond")
		defer close(out)
		defer close(responseChannel)
		opts := IterateOptions{DontSendMessages: true, ProcessFailingMessages: true}
		return IterateMessages(ctx, in, out, opts, func(msg *Message) bool {
			msg.Send(ctx, responseChannel)
			return true
		})
	})
	return out
}

func main() {
	ctx := helper.NewContext()
	testIn := make(chan Message)
	c, responseChannel := subscribe(ctx, testIn)
	c = process(ctx, c)
	c = process(ctx, c)
	c = process(ctx, c)
	c = process(ctx, c)
	testOut := respond(ctx, c, responseChannel)

	time.Sleep(1 * time.Millisecond)

	// Send and receive a test message
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
