package main

import (
	"fmt"
	"time"
)

func main() {

	ch := time.Tick(1 * time.Second)
	go func() {
		for range ch {
			fmt.Println("Hello, world!")
		}
	}()

	time.Sleep(5 * time.Second)
}
