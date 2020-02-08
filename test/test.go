package main

import "fmt"

func main() {

	naturals := make(chan int, 1)

	squares := make(chan int, 1)

	// Counter

	go func() {

		for x := 0; x < 100; x++ {

			naturals <- x

		}

	}()

	// Squarer

	go func() {

		for {

			x := <-naturals

			squares <- x * x

		}

	}()

	// Printer (in main goroutine)

	for {

		fmt.Println(<-squares)

	}

}
