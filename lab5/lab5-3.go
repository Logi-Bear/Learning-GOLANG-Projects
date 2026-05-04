package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const MAX_VALUE = 100

type Queue struct {
	elements []int
}

func (queue *Queue) dequeue() int {
	temp := queue.elements[0]
	queue.elements = queue.elements[1:]
	return temp
}

func display(queue *Queue) {
	for i := range queue.elements {
		fmt.Printf("%d ", queue.elements[i])
	}
	fmt.Printf("\n")
}

func (queue *Queue) enqueue(value int) {
	queue.elements = append(queue.elements, value)
}

func generateRandNumber() int {
	max := big.NewInt(MAX_VALUE)

	tmp, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	return int(tmp.Int64())
}

func (queue *Queue) getLen() int {
	return len(queue.elements)
}

func isExists(queue *Queue, searchKey int) bool {
	exists := false
	for i := range queue.elements {
		if queue.elements[i] == searchKey {
			exists = true
		}
	}
	return exists
}

func main() {
	queue := Queue{}

	for i := 0; i < 10; i++ {
		queue.enqueue(generateRandNumber())
	}

	display(&queue)

	length := queue.getLen()

	fmt.Println("Queue len =", length)

	for length > 0 {
		fmt.Println("Dequeued ", queue.dequeue())
		display(&queue)
		length--
	}
}
