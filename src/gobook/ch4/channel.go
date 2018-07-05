package main

import "fmt"

func Count(ch chan int) {
	fmt.Println("Counting")
	ch <- 1

}

func main() {
	chs := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count(chs[i])
	}
	var result int
	for _, ch := range chs {
		result += <-ch
	}
	fmt.Println(result)
}
