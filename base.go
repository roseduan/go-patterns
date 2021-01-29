package main

import (
	"fmt"
	"reflect"
)

//切片，接口，性能

//1.切片的内部结构，源码在 runtime/slice.go 中可查看
//type slice struct {
//	array unsafe.Pointer
//	len int
//	cap int
//}

//2.接口编程

//一个简单的例子
type Person struct {
	Name   string
	Sexual string
	Age    int
}

//下面这两个方法有什么区别？
func (p *Person) Print() {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d", p.Name, p.Sexual, p.Age)
}

func PrintPerson(p *Person) {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d", p.Name, p.Sexual, p.Age)
}

//再来看另一个例子
type Country struct {
	Name string
}

type City struct {
	Name string
}

type Printable interface {
	PrintStr()
}

func (c *Country) PrintStr() {
	fmt.Println(c.Name)
}

func (c *City) PrintStr() {
	fmt.Println(c.Name)
}

//使用结构体嵌入改进上面的例子
type WithName struct {
	Name string
}

type Dog struct {
	WithName
}

type Cat struct {
	WithName
}

func (w *WithName) PrintStr() {
	fmt.Println(w.Name)
}

//再简化上面的例子
type Stringable interface {
	ToString() string
}

func PrintName(s Stringable) {
	fmt.Println(s.ToString())
}

func (c *Country) ToString() string {
	return c.Name
}

func (c *City) ToString() string {
	return c.Name
}

//3.性能方面的小建议
//如果需要把数字转换成字符串，使用 strconv.Itoa() 比 fmt.Sprintf() 要快一倍左右
//尽可能避免把String转成[]Byte ，这个转换会导致性能下降
//如果在 for-loop 里对某个 Slice 使用 append()，请先把 Slice 的容量扩充到位，这样可以避免内存重新分配以及系统自动按 2 的 N 次方幂进行扩展但又用不到的情况，从而避免浪费内存
//使用StringBuffer 或是StringBuild 来拼接字符串，性能会比使用 + 或 +=高三到四个数量级
//尽可能使用并发的 goroutine，然后使用 sync.WaitGroup 来同步分片操作
//避免在热代码中进行内存分配，这样会导致 gc 很忙
//尽可能使用 sync.Pool 来重用对象
//使用 lock-free 的操作，避免使用 mutex，尽可能使用 sync/Atomic包（关于无锁编程的相关话题，可参看《无锁队列实现》或《无锁 Hashmap 实现》）
//	reference: https://coolshell.cn/articles/8239.html  https://coolshell.cn/articles/9703.html
//使用 I/O 缓冲，I/O 是个非常非常慢的操作，使用 bufio.NewWrite() 和 bufio.NewReader() 可以带来更高的性能
//对于在 for-loop 里的固定的正则表达式，一定要使用 regexp.Compile() 编译正则表达式。性能会提升两个数量级
//考虑使用 protobuf 或 msgp 而不是 JSON，因为 JSON 的序列化和反序列化里使用了反射
//使用 Map 的时候，使用整型的 key 会比字符串的要快，因为整型比较比字符串比较要快

func main() {
	fmt.Println("--------切片实践--------")
	nums := make([]int, 5)
	nums[0] = 10
	nums[1] = 20
	nums[2] = 25

	fmt.Println("nums = ", nums)
	another := nums[:3]

	//another和nums共享同一个底层数组，因此这里的改动会影响到nums
	another[0] = 100
	fmt.Println("nums = ", nums)

	//如果新的切片发生了扩容，那么则不会影响到原来的切片了
	arr := make([]int, 3)
	arr[0] = 0
	arr[1] = 1
	arr[2] = 2

	arr2 := arr[:]
	arr2 = append(arr2, 10, 20)
	arr2[0] = 100
	fmt.Println(arr, arr2)

	//切片之间的比较可以使用 reflect.DeepEqual
	v1 := []int{}
	v2 := []int{}
	fmt.Println(reflect.DeepEqual(v1, v2)) //true

	m1 := map[int]string{1: "a", 2: "b", 3: "c"}
	m2 := map[int]string{2: "b", 3: "c", 1: "a"}
	fmt.Println(reflect.DeepEqual(m1, m2)) //true

	sli1 := []int{1, 2, 3}
	sli2 := []int{1, 2, 3}
	fmt.Println(reflect.DeepEqual(sli1, sli2)) //true

	fmt.Println("--------接口编程实践--------")
	c := Country{"China"}
	c.PrintStr()

	city := City{"Shanghai"}
	city.PrintStr()

	bobby := Dog{WithName{"Bobby"}}
	mimi := Cat{WithName{"mimi"}}
	bobby.PrintStr()
	mimi.PrintStr()

	PrintName(&c)
	PrintName(&city)
}
