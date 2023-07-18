package package_a

import (
	"fmt"
	"test/tool/cycle/package_b"
	"test/tool/cycle/package_i"
)

type PackageA struct {
	B package_i.PackageBInterface
}

func (a PackageA) PrintA() {
	fmt.Println("I'm a!")
}

func (a PackageA) PrintAll() {
	a.PrintA()
	a.B.PrintB()
	package_b.PackageB{}.PrintB()
}
