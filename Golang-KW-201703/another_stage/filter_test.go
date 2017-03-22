package main

import (
	"github.com/cksmith/go-talks/Golang-KW-201703/helper"
	"testing"
)

// START OMIT

func testQuote(t *testing.T, quote string, expectedPass bool) {
	ctx := helper.NewContext()
	in := make(chan Message)
	out := filter(ctx, in)
	in <- Message{Quote: quote}
	msg := <-out
	close(in)
	if expectedPass && msg.LastError != nil {
		t.Fail()
	}
	if !expectedPass && msg.LastError == nil {
		t.Fail()
	}
}

func TestFilterPassesAppropriateQuote(t *testing.T) {
	testQuote(t, "Hi mom", true)
}

func TestFilterFailsInappropriateQuote(t *testing.T) {
	testQuote(t, "Perl is fun", false)
}

// STOP OMIT
