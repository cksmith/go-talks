package main

import (
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201701/quote"
	"math/rand"
	"time"
)

func quoter() {
	for i := 0; ; i++ {
		s, _ := quote.Get()
		fmt.Println(s, i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	go quoter()
	fmt.Println("Running")
	time.Sleep(5 * time.Second)
	fmt.Println("Done")
}
