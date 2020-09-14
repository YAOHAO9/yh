package main

// import (
// 	"fmt"
// 	"runtime"
// 	"time"

// 	_ "net/http/pprof"
// )

// func main() {

// 	c := make(chan int, 1)
// 	go func() {
// 		time.Sleep(time.Second * 2)
// 		c <- 0
// 		fmt.Println("go routine finished")
// 	}()

// 	select {
// 	case <-c:
// 		fmt.Println("go routine")
// 	case <-time.After(time.Second * 1):
// 		fmt.Println("Timeout")
// 		fmt.Println("runtime.NumGoroutine()1:", runtime.NumGoroutine())
// 	}

// 	time.Sleep(time.Second * 3)
// 	fmt.Println("runtime.NumGoroutine()2:", runtime.NumGoroutine())
// }
