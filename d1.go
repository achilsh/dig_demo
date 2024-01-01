package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"go.uber.org/dig"
)

//注入具体类
type DemoOne struct {
	A int
	B string
	C float32
	D float64 
	E bool
}

func NewDemoONe() *DemoOne {
	done := &DemoOne{
		A:100,
		B: "b",
	}
	return done
}

///
type DemoTwo struct {
	AA *DemoOne
	BB  int32
}
func NewDemoTwo(one *DemoOne)  *DemoTwo{
	dt := &DemoTwo {
		AA: one,
		BB: 100,
	}
	return dt
}


//注入接口
type IPerson interface {
	GetName() string 
	SetName(s string)
}
//
type Teacher struct {
	Name string 
}
func(tr *Teacher)GetName() string {
	return tr.Name
}
func(tr *Teacher)SetName(s string) {
	tr.Name = s
}
func NewPerson() IPerson {
	return &Teacher{}
}

//
type University struct {
	tearcher IPerson
}
func NewUniversity(T IPerson) *University {
	u := &University{
		tearcher: T,
	}
	fmt.Println("call create university.")
	return u
}


/// dig.IN 包装 依赖的多个项。
type Address struct {
	addr string
} 
func NewAddress() *Address {
	a := &Address{
		addr: "shenzhen",
	}
	return a
}
func (a *Address)String() string {
	return a.addr
}

type Age struct {
	age int32
}
func (a *Age) String() string {
	return "age: " + strconv.Itoa(int(a.age))
}
//
func NewAge() *Age {
	a := &Age {
		age: 30,
	}
	return a
}

type ManagerIn struct {
	dig.In 
	B1 *Address 
	B2 *Age
}
type Manager struct {
	A1 *Address
	A2 *Age
}

func NewManager(m ManagerIn) *Manager {
	o := &Manager{
		A1: m.B1,
		A2: m.B2,
	}
	return o
}

////////// dig.out 包装多个创建的对象。 定义的包装类型只能用在 provide 的参数函数的返回值。而且不能是指针类型。
type Out1Data struct {
	A int32
}
type Out2Data struct {
	B float32
}

type MoreOuput struct {
	dig.Out 
	O1 *Out1Data
	O2 *Out2Data
}

func NewMoreOut() MoreOuput{
	r := MoreOuput{
		O1: new(Out1Data),
		O2: new(Out2Data),
	}

	r.O1.A = 123
	r.O2.B = float32(123.123)
	
	return r
}

//// 依赖一种类型多个变量。 使用 name来区分。
type SomeType struct {
	A int32
}


//创建返回值对象，用于包装 两个同一类型的不同全局变量。包装创建的对象。
type WrapSomeTypeOut struct {
	dig.Out
	S11 *SomeType `name:"s1"`
	S22 *SomeType `name:"s2"`
}

//创建同一类型的不同全局变量。
func NewWrapSomeType() WrapSomeTypeOut  { //返回参数值不能用 指针变量
	r := WrapSomeTypeOut{
	}
	r.S11 = &SomeType{
		A: 1000,
	}
	r.S22 = &SomeType{
		A: 20000,
	}

	return r
}


//包装同一类型不同的全局变量。包装依赖项。
type WrapSomeType struct {
	dig.In
	S1 *SomeType `name:"s1"`  //依赖同一类型不同变量1 
	S2 *SomeType `name:"s2"`  //依赖同一类型不同变量2
}

//BusiDataSomeType 依赖同一类型的多个变量
type BusiDataSomeType struct {
	a *SomeType //同一类型的不同变量
	b *SomeType //同一类型的不同变量

}

//定义 创建 BusiDataSomeType对象的函数，作为 dig.Provide()入参 被调用。 
func NewBusiDataSomeType(in WrapSomeType) *BusiDataSomeType{
	r := &BusiDataSomeType{
		a: in.S1,
		b: in.S2,
	}
	return r
}




type OptionOne struct {
	A int32
}

func NewOptionOne() *OptionOne{
	return &OptionOne{
		A: 100,
	}
}

type OptionTwo struct {
	B int32
}
func NewOptionTwo() *OptionTwo {
	return &OptionTwo{
		B: 2000,
	}
}

// option 可选依赖只能在参数对象上使用.
// 有些依赖并不是必须的，因此可以通过结构体tag：optional:"true"把某些依赖标记成可选的,
// 在参数对象的可选字段上，即使没有提供可选字段的provide函数，依赖于参数对象的对象也能构建成功，比如：
type OptionIn struct {
	dig.In 
	OA *OptionOne `optional:"true"` //可以不用提供创建OptionOne的provide函数。
	OB *OptionTwo 
}

type OutOption struct {
	a *OptionOne
	b *OptionTwo
}
func NewOutOPtion(in OptionIn) *OutOption {
	o := &OutOption{
		a: in.OA,
		b: in.OB,
	}
	return o
}






///////////////////////////////////////////////////////////////////////
func Run() *dig.Container {
	x := dig.New()
	//
	x.Provide(NewDemoONe)
	x.Provide(NewDemoTwo)
	//
	e := x.Invoke(func(d *DemoTwo) {
		fmt.Println("d: ", d.AA.A, d.AA.B, d.AA.C, d.AA.D)
	})
	if e !=nil {
		fmt.Println("invokde fail, e: ", e)
	}

	//
	fmt.Println("call provide...")
	x.Provide(NewUniversity) // NewUniversity 仅仅会被调用一次;
	x.Provide(NewUniversity) // NewUniversity 仅仅会被调用一次;
	//
	x.Provide(NewPerson)  // NewUniversity,NewPerson 虽然存在依赖关系，但是在Provide调用没有必要按依赖关系来写;
	x.Provide(NewUniversity) // NewUniversity 仅仅会被调用一次;
	//

	fmt.Println("call invoke...")
	e = x.Invoke(func(u* University) {
		u.tearcher.SetName("joney")
		fmt.Println("name: ", u.tearcher.GetName())
	})
	if e != nil {
		fmt.Println("inject university fail, e: ",e )
	}

	//
	x.Provide(NewAge)
	x.Provide(NewAddress)
	x.Provide(NewManager)
	//
	e = x.Invoke(func(m *Manager) {
		fmt.Println("call invokde manager: ")
		fmt.Println(m.A1)
		fmt.Println(m.A2)
	})
	if e != nil {
		fmt.Println("call invokde manager fail, e: ", e)
	}

	//多个对象放在一个dig.out结构体内一起被创建。
	x.Provide(NewMoreOut)
	x.Invoke(func(a *Out1Data) { // 可以使用dig.out中的某个成员变量；且类型保持一致。
		fmt.Println("a in dig.out, value: ", a.A)

	})
	fmt.Println(".....------------------")
	x.Invoke(func(b *Out2Data) { // 不可使用 dig.out整个结构体作为 函数参数，只能是dig.out的某个成员变量。
		fmt.Println("b in dig.Out, value: ", b.B)
	})
	

	//告诉 dig 根据一些依赖 创建 一些对象。 对同一类型多个变量的依赖：
	x.Provide(NewWrapSomeType)
	x.Provide(NewBusiDataSomeType)
	
	
	// Invoke 用于创建 BusiDataSomeType 对象并运行函数。
	e = x.Invoke(func(w *BusiDataSomeType) {
		fmt.Println("w: ", w.a.A, w.b.A)
	})
	if e != nil {
		fmt.Println("invokde fail, e: ", e)
	}


	//option:
	x.Provide(NewOptionTwo)
	x.Provide(NewOutOPtion)
	e = x.Invoke(func(o *OutOption) {
		if o.a !=nil {
			fmt.Println("o.a: ", o.a.A)
		}
		if o.b != nil {
			fmt.Println("o.b: ", o.b.B)
		}
	})
	if e != nil {
		fmt.Println("run option fail,e: ", e)
	}







	//
	return x
}

func GenGraph(d *dig.Container) {
	buf := bytes.Buffer{}
	e :=dig.Visualize(d, &buf)
	if e !=nil {
		fmt.Println("visualize fail, e: ", e)
		return 
	}

	dotfile := filepath.Join(".", "demo.dot")
	e = os.WriteFile(dotfile, buf.Bytes(), 0644)
	if e != nil {
		fmt.Println("write file fail, e: ", e)
		return 
	}

	g := graphviz.New()
	graph, e := graphviz.ParseBytes(buf.Bytes())
	if e != nil {
		fmt.Println("graphvize fail, e: ", e)
		return 
	}
	graph.SetRankDir(cgraph.RLRank) //LRRank RLRank
	e = g.RenderFilename(graph, graphviz.SVG, "demo"+".svg")
	if e != nil {
		fmt.Println("grap fail, e: ", e)
		return 
	}
}