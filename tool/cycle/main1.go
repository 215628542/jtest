package main

import (
	"test/tool/cycle/package_a"
	"test/tool/cycle/package_b"
)

func main() {
	a := new(package_a.PackageA)
	b := new(package_b.PackageB)
	a.B = b
	b.A = a
	a.PrintAll()
	b.PrintAll()
}
