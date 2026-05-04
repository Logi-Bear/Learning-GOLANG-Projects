package main

import (
	"fmt"
	"math/rand"
	"time"
)

// The function returns a receive-only channel <-chan string
func boring(msg string) <-chan string {
	ch := make(chan string)
	go func() { // Launch the goroutine from inside the function
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s: %d", msg, i) // Send a message into the channel
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return ch // Return the channel to the caller
}

func main() {
	joe := boring("Joe")
	ann := boring("Ann")

	for i := 0; i < 5; i++ {
		// Receive messages from both channels sequentially
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}
	fmt.Println("You're both boring; I'm leaving.")
}
