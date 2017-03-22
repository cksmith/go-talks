package main

import (
	"context"
	"fmt"
)

func main() {
	// START1 OMIT
	// Obtain a context with a cancellation signal (the done function)
	ctx, done := context.WithCancel(context.Background())
	done()       // Trigger cancellation
	<-ctx.Done() // Wait until the context is cancelled
	fmt.Println("Request done")
	// STOP1 OMIT

	// START2 OMIT
	childCtx, childDone := context.WithCancel(ctx)
	<-childCtx.Done() // Parent context was already cancelled
	fmt.Println("Child request done")
	// STOP2 OMIT
	_ = childDone
}
