package package_b

import (
	"fmt"
	"test/tool/cycle/package_i"
)

type PackageB struct {
	A package_i.PackageAInterface
}

func (b PackageB) PrintB() {
	fmt.Println("I'm b!")
}

func (b PackageB) PrintBB() {
	fmt.Println("I'm bb!")
}

func (b PackageB) PrintAll() {
	b.PrintB()
	b.A.PrintA()
}
