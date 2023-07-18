package main

import (
	"fmt"
	"sync"
)

func showNum(numChan chan struct{}, wordChan chan struct{}) {
	defer wg.Done()
	i := 0
	for {
		i++
		select {
		case <-wordChan:
			fmt.Printf("%v", i)
			numChan <- struct{}{}
		case <-endChan:
			return
		}
	}
}

func showWord(numChan chan struct{}, wordChan chan struct{}) {
	defer wg.Done()
	endChan = make(chan bool, 1)
	i := 'A'
	word := ""
	for {
		word = string(i)
		select {
		case <-numChan:
			if word > "Z" {
				endChan <- true
				return
			}
			fmt.Printf("%v", string(i))
			i++
			wordChan <- struct{}{}
		}
	}
}

var endChan chan bool
var wg sync.WaitGroup

func main2() {

	numChan := make(chan struct{}, 1)
	wordChan := make(chan struct{}, 1)
	wordChan <- struct{}{}
	wg.Add(2)
	go showNum(numChan, wordChan)
	go showWord(numChan, wordChan)
	wg.Wait()
}
