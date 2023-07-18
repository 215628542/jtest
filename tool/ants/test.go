package main

import (
	"fmt"
	"github.com/spf13/cast"
	"runtime"
	"time"
)

func t(c chan int) {
	//fmt.Println(<-c)
	time.Sleep(1 * time.Second)
	select {
	case m := <-c:
		panic("协程退出" + cast.ToString(m))
		//default:
		//	panic("协程退出")
	}

}

func main1() {

	size := 5
	chanSlic := make([]chan int, size)
	// 创建多个协程，监听任务消费
	for i := 0; i < size; i++ {
		chanSlic[i] = make(chan int, 1)
		go func(i int) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			t(chanSlic[i])
			fmt.Println("testteset")
		}(i)
	}
	fmt.Println("====")
	num := runtime.NumGoroutine()
	fmt.Println(num)

	for k, cl := range chanSlic {
		fmt.Println("kkk")
		fmt.Println(k)
		cl <- k
		close(cl)
		time.Sleep(3 * time.Second)
		fmt.Println(runtime.NumGoroutine())
	}

	fmt.Println("skfjsidfj")
	time.Sleep(5 * time.Second)
	num2 := runtime.NumGoroutine()
	fmt.Println(num2)

	fmt.Println("end")

	time.Sleep(10 * time.Minute)
}
