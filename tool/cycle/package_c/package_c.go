package package_c

import (
	"test/tool/cycle/package_a"
	"test/tool/cycle/package_b"
)

// 新建公共组合包（子包），在组合包中组合调用
type CombileAB struct {
	A *package_a.PackageA
	B *package_b.PackageB
}

func (c CombileAB) PrintAll() {
	c.A.PrintA()
	c.B.PrintBB()
}
