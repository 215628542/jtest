package main

import "fmt"

// 抽象的攻击技能
type Attack interface {
	Fight()
}

// 具体的技能
type KaiQiang struct {
}

func (k *KaiQiang) Fight() {
	fmt.Println("使用了开枪技能")
}

type Hero struct {
	Name   string
	attack Attack // 攻击方式
}

func (h *Hero) Skill() {
	fmt.Println(h.Name + "使用了技能")
	h.attack.Fight() // 使用具体的战斗方式
}

// 适配器  跟其他类完全不相关，其实就是新增加需要嵌入原有的类里面
type PowerOff struct {
}

func (p *PowerOff) ShowDown() {
	fmt.Println("关机")
}

// 适配伦攻击方式，直接实现了Attack抽象接口
type Adapter struct {
	powerOff *PowerOff // 将需要替换的方法的类组合进来
}

// 适配攻击方式，直接实现了Attack抽象接口
func (a *Adapter) Fight() {
	a.powerOff.ShowDown() // 将需要替换的方法在这里引入
}

func NewAdapter(p *PowerOff) *Adapter {
	return &Adapter{powerOff: p}
}

func main() {

	gailun := Hero{Name: "盖伦", attack: new(KaiQiang)}
	gailun.Skill()

	// 改造 - 盖伦使用技能从开枪改为关机
	gailun2 := Hero{Name: "盖伦", attack: NewAdapter(new(PowerOff))}
	gailun2.Skill()

}
