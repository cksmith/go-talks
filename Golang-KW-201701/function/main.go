package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/helper"
	"math/rand"
	"time"
)

// START OMIT

func quoter() {
	for i := 0; ; i++ {
		s, _ := helper.GetQuote()
		fmt.Println(s, i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	quoter()
}

// STOP OMIT
