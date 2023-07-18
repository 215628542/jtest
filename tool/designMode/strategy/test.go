package main

import "fmt"

type ChuXiao interface {
	Youhui()
}

type FanganA struct {
}

func (a *FanganA) Youhui() {
	fmt.Println("打8折")
}

type FanganB struct {
}

func (b *FanganB) Youhui() {
	fmt.Println("满100，减20")
}

type People struct {
	c ChuXiao
}

func (p *People) MakeYouhui(c ChuXiao) {
	p.c = c
}

func (p *People) Buy() {
	p.c.Youhui()
}

func main() {

	p := new(People)

	p.MakeYouhui(new(FanganB))
	p.Buy()

	p.MakeYouhui(new(FanganA))
	p.Buy()
}
