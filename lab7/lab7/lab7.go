package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	id     int
	status string
}

var (
	orderIDCounter int
	orderIDMu      sync.Mutex
)

const MAX = 20

func main() {
	ordersChannel := make(chan *Order)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < MAX; i++ {
			ordersChannel <- &Order{
				id:     i,
				status: "Received",
			}
		}
		close(ordersChannel)
	}()

	reportOrderStatus(<-ordersChannel)

	workerCount := 3
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			for order := range ordersChannel {
				processOrder(order)
			}
		}(i)
	}

	updateOrderStatus(ordersChannel)

	wg.Wait()
	fmt.Println("All Done")
}

func generateOrder() *Order {
	var orderID int

	orderIDMu.Lock()
	orderID = orderIDCounter

	time.Sleep((time.Duration(rand.Intn(100))) * time.Microsecond)

	orderIDCounter++
	orderIDMu.Unlock()

	return &Order{
		id: orderID, status: "Received",
	}
}

func generateOrderStatus() string {
	status := []string{"Received", "Processing", "Served"}[rand.Intn(3)]

	return status
}

func processOrder(order *Order) {
	time.Sleep((time.Duration(rand.Intn(500))) * time.Microsecond)
	order.status = "Served"
	fmt.Printf("Processing order %d\n", order.id)
}

func reportOrderStatus(orders []*Order) {
	time.Sleep(1 * time.Microsecond)
	fmt.Println("\n--- Order Status ---")

	for _, order := range orders {
		fmt.Printf("Order %d: %s\n", order.id, order.status)
	}
}

func updateOrderStatus(order *Order) {
	time.Sleep((time.Duration(rand.Intn(300))) * time.Microsecond)

	order.status = generateOrderStatus()

	fmt.Printf("Updating order %d status to %s\n", order.id, order.status)
}
