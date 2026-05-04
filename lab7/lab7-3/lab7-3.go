package main

import "fmt"

func sum(list []int, ch chan int) {
	sum := 0
	for _, v := range list {
		sum += v
	}
	ch <- sum // send sum to ch
}

func main() {
	list := []int{7, 2, 8, -9, 4, 0}

	ch := make(chan int)

	fmt.Println(list[:len(list)/2], list[len(list)/2:])

	go sum(list[:len(list)/2], ch)
	go sum(list[len(list)/2:], ch)
	x, y := <-ch, <-ch // receive from ch

	fmt.Println(x, y)
}
