package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)


func main() {
	// question 1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("hi")
	}()
	wg.Wait()

	// question 2
	ch := make(chan string)
	go func ()  {
		ch <- "ping"	
	}()

	chValue := <-ch
	fmt.Println(chValue)

	// question 3
	// launch 5 goroutines
	ch2 := make(chan int, 5)
	for i := 1; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch2 <- i * i
		}(i)
	}
	wg.Wait()
	close(ch2)
	sum := 0
	for v := range ch2 {
		fmt.Println(v)
		sum += v
	}
	fmt.Println("Sum:", sum)

	ch3 := make(chan int, 5)
	close(ch3)
	// check if ok
	if _, ok := <-ch3; ok {
		fmt.Println(ok)
	}
	// question 5
	// worker pool
	workerCount := 4
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg2 sync.WaitGroup
	for w := 1; w <= workerCount; w++ {
		wg2.Add(1)
		go func(){
			defer wg2.Done()
			for j := range jobs {
				results <- j * j * j
			}
		}()
	}
	
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)
	wg2.Wait()
	// sum results
	close(results)
	sum2 := 0
	for r := range results {
		fmt.Println(r)
		sum2 += r
	}
	fmt.Println("Sum of cubes:", sum2)


	// question 6
	before := runtime.NumGoroutine()
	fmt.Println("Before:", before)

	leakExample()

	time.Sleep(100 * time.Millisecond)

	after := runtime.NumGoroutine()
	fmt.Println("After:", after)

	time.Sleep(2 * time.Second)

	later := runtime.NumGoroutine()
	fmt.Println("Later:", later)
}


func leakExample() {
	ch := make(chan int)
	go func() {
		<-ch // blocks forever
		fmt.Println("This will never print")
	}()
}