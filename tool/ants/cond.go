package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main_c1() {
	c := sync.NewCond(&sync.Mutex{})
	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
			//time.Sleep(1 * time.Second)

			// 加锁更改等待条件
			c.L.Lock()
			ready++
			c.L.Unlock()

			log.Printf("运动员#%d 已准备就绪\n", i)
			// 广播唤醒所有的等待者
			c.Broadcast()
			//time.Sleep(1 * time.Second)
			//c.Signal()

		}(i)
	}

	c.L.Lock()
	fmt.Println("===========")
	for ready != 10 {
		c.Wait()
		log.Println("裁判员被唤醒一次")
	}
	//time.Sleep(2 * time.Second)
	//fmt.Println("2222222222")
	c.L.Unlock()

	//所有的运动员是否就绪
	log.Println("所有运动员都准备就绪。比赛开始，3，2，1, ......")
}
