package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	ok      = false                           // 是否做完饭
	food    = ""                              // 饭的名字
	rwMutex = &sync.RWMutex{}                 // 读写锁，用于锁定饭名的修改
	cond    = sync.NewCond(rwMutex.RLocker()) // 条件变量使用读写锁中的读锁
)

func makeFood() {
	// 做饭使用写锁（当然因为只有一个做饭协程，该锁并无实际意义）
	rwMutex.Lock()
	defer rwMutex.Unlock()
	fmt.Println("食堂开始做饭！")
	time.Sleep(1 * time.Second)
	ok = true
	food = "鱼香肉丝"
	fmt.Println("食堂做完饭了！")
	cond.Broadcast()
}

func waitToEat() {
	cond.L.Lock()
	defer cond.L.Unlock()
	for !ok {
		cond.Wait()
	}
	fmt.Println("总算吃到饭了，这顿吃的是", food)
}

func main_cond2() {
	for i := 0; i < 3; i++ {
		go waitToEat()
	}
	go makeFood()
	time.Sleep(30 * time.Second)
}
