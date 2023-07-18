package main

import (
	"fmt"
	"github.com/spf13/cast"
	"runtime"
	"time"
)

func tt(m int) {

	if m == 2 {
		panic("panic")
	}

	//time.Sleep(1 * time.Second)
	fmt.Println("输出" + cast.ToString(m))
}

func main2() {

	c1 := make(chan int, 1)
	c2 := make(chan int, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			select {
			case m := <-c1:
				tt(m)
			}
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			select {
			case m := <-c2:
				tt(m)
			}
		}
	}()

	c1 <- 1
	c2 <- 1
	time.Sleep(2 * time.Second)

	fmt.Println("====")
	fmt.Println(runtime.NumGoroutine())

	c1 <- 2
	c2 <- 2
	// 输入2，触发panic时，协程数减少2，说明panic会让协程退出，这样就会减少工作的协程数

	time.Sleep(2 * time.Second)
	fmt.Println("====")
	fmt.Println(runtime.NumGoroutine())

	fmt.Println("end")

	time.Sleep(10 * time.Minute)
}
