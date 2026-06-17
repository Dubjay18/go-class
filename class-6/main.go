package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	// question 1
	ch1 := make(chan int, 2)
	ch2 := make(chan int, 2)

	ch1 <- 1
	ch1 <- 2

	select {
	case v := <-ch1:
		println("received from ch1:", v)
	case v := <-ch2:
		println("received from ch2:", v)
	default:
		println("no value received")
	}

	// question 2
	// Non-blocking receive
	ch3 := make(chan int, 1)
	select {
	case v := <-ch3:
		println("received from ch3:", v)
	default:
		println("nothing ready")
	}
	// question 3
	//start a goroutine that sends a value into a channel after 50ms
	select {
	case v := <-ch3:
		println("received from ch3:", v)
	case <-time.After(10 * time.Millisecond):
		println("timeout")
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Question 4
	//WithCancel.
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopped:", ctx.Err())
				return
			default:
				fmt.Println("tick")
				time.Sleep(20 * time.Millisecond)
			}
		}
	}()

	cancel() // signal the goroutine to stop

	// Question 5
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx2.Done():
				fmt.Println("stopped2:", ctx2.Err())
				return
			default:
				fmt.Println("tick")
				time.Sleep(20 * time.Millisecond)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond) // wait for the goroutine to finish
	cancel2()                          // signal the goroutine to stop

	// question 6
	b := runtime.NumGoroutine()
	fmt.Println("number of goroutines:", b)

	ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second)
	jobs := make(chan int)
	go worker(ctx3, jobs)

	time.Sleep(100 * time.Millisecond)

	during := runtime.NumGoroutine()
	fmt.Println("number of goroutines during work:", during)

	cancel3() // signal the worker to stop

	time.Sleep(100 * time.Millisecond) // wait for the worker to finish
	after := runtime.NumGoroutine()
	fmt.Println("number of goroutines after work:", after)

	if after != b {
		fmt.Println("goroutine leak detected")
	} else {
		fmt.Println("no goroutine leak detected")
	}
}

func worker(ctx context.Context, jobs <-chan int) {
	for {
		select {
		case <-ctx.Done(): // cancelled or timed out -> leave
			fmt.Println("worker exiting:", ctx.Err())
			return
		case j, ok := <-jobs:
			if !ok {
				return // jobs channel closed
			}
			fmt.Println("processing", j)
		}
	}
}
