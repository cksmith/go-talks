package helper

import (
	"context"
	"sync"
)

// START1 OMIT

type Context struct {
	context.Context
	Wg       sync.WaitGroup
	DoneFunc context.CancelFunc
}

func NewContext() *Context {
	ctx := new(Context)
	ctx.Context, ctx.DoneFunc = context.WithCancel(context.Background())
	return ctx
}

// STOP1 OMIT

// START2 OMIT

// Get the channel that indicates when the pipeline is shutting down
func (ctx *Context) Done() <-chan struct{} {
	return ctx.Context.Done()
}

// Wait for the pipeline to shut down
func (ctx *Context) Wait() {
	ctx.Wg.Wait()
}

// Shut down the pipeline and wait for shutdown to complete
func (ctx *Context) Stop() {
	ctx.DoneFunc()
	ctx.Wait()
}

// STOP2 OMIT

// START3 OMIT

// Launch a stage of the pipeline (in a goroutine) that cancels the pipeline
// if the function returns false
func (ctx *Context) Run(f func() bool) {
	ctx.Wg.Add(1)
	go func() {
		defer ctx.Wg.Done()
		if !f() {
			ctx.DoneFunc()
		}
	}()
}

// STOP3 OMIT
