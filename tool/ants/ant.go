package main

import (
	"fmt"
	"github.com/spf13/cast"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
)

var sum int32

func panicHandler(err interface{}) {
	fmt.Println(123123123)
	fmt.Fprintln(os.Stderr, err)
}

func myFunc(i interface{}) {
	defer wg.Done()
	time.Sleep(1 * time.Second)
	n := i.(int32)
	if n == 2 {
		panic("test panic" + string(sum))
	}
	atomic.AddInt32(&sum, n)

	fmt.Printf("run with %d\n", n)

}

// 使用闭包接收参数值
func demoFunc(i int) func() {
	return func() {
		defer wg.Done()
		//time.Sleep(10 * time.Millisecond)
		fmt.Println("demoFunc-" + cast.ToString(i))
	}
}

var wg sync.WaitGroup

func main1212() {
	defer ants.Release()

	runTimes := 20

	// Use the common pool.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(demoFunc(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")

	return

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
	}, ants.WithPanicHandler(panicHandler))
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		// Invoke是异步执行
		_ = p.Invoke(int32(i))
		fmt.Println("=====" + string(i))
	}

	// wait是等待ants的所有协程执行完毕
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}
