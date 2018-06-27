package main

import (
	"time"
	"fmt"
)

func count(ch chan int) {
	time.Sleep(2e9)
	fmt.Println("Counting")
	ch <- 1
}

func main() {
	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(1e9)
		timeout <- true
	}()

	ch := make(chan int)
	go count(ch)

	select {
	case result := <-ch:
		fmt.Println(result)
	case <-timeout:
		fmt.Println("timeout")
	}
}
