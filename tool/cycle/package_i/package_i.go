package package_i

// 新建公共接口包（父包），将需要循环调用的函数或方法抽象为接口

type PackageAInterface interface {
	PrintA()
}

type PackageBInterface interface {
	PrintB()
}
