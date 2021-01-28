package main

import (
	"fmt"
	"reflect"
)

//切片，接口，时间，性能

//1.切片的内部结构，源码在 runtime/slice.go 中可查看
//type slice struct {
//	array unsafe.Pointer
//	len int
//	cap int
//}

//2.接口编程

//一个简单的例子
type Person struct {
	Name string
	Sexual string
	Age int
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
	fmt.Println(reflect.DeepEqual(v1, v2))	//true

	m1 := map[int]string{1: "a", 2: "b", 3: "c"}
	m2 := map[int]string{2: "b", 3: "c", 1: "a"}
	fmt.Println(reflect.DeepEqual(m1, m2))	//true

	sli1 := []int{1, 2, 3}
	sli2 := []int{1, 2, 3}
	fmt.Println(reflect.DeepEqual(sli1, sli2))	//true
}
