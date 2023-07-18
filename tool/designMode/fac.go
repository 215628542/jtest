package main

import "fmt"

type Fruit interface {
	show()
}

type Apple struct {
}

func (a *Apple) show() {
	fmt.Println("apple show")
}

type FruitFac struct {
}

func (f *FruitFac) CreateFruit(fu Fruit) {
	fu.show()
}

func main3() {

}
