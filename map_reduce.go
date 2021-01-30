package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

//函数式操作 map reduce filter

func MapStrToStr(arr []string, fn func(s string) string) []string {
	var newArr []string
	for _, v := range arr {
		newArr = append(newArr, fn(v))
	}

	return newArr
}

func MapStrToInt(arr []string, fn func(s string) int) []int {
	var res []int
	for _, v := range arr {
		res = append(res, fn(v))
	}

	return res
}

//reduce函数示例
func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, v := range arr {
		sum += fn(v)
	}

	return sum
}

//filter函数示例
func Filter(arr []string, fn func(s string) bool) []string {
	var newArr []string
	for _, v := range arr {
		if fn(v) {
			newArr = append(newArr, v)
		}
	}

	return newArr
}

//业务示例
type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   float32
}

var list = []Employee{
	{"Hao", 44, 4, 8000},
	{"Bob", 34, 10, 5000},
	{"Alice", 23, 5, 9000},
	{"Jack", 26, 3, 4000},
	{"Tom", 48, 9, 7500},
	{"Marry", 29, 7, 6000},
	{"Mike", 32, 8, 4000},
}

func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) (count int) {
	for _, v := range list {
		if fn(&v) {
			count++
		}
	}
	return
}

func EmployeeFilterIn(list []Employee, fn func(e *Employee) bool) []Employee {
	var res []Employee
	for _, v := range list {
		if fn(&v) {
			res = append(res, v)
		}
	}

	return res
}

func EmployeeSumIf(list []Employee, fn func(e *Employee) int) int {
	sum := 0
	for _, v := range list {
		sum += fn(&v)
	}

	return sum
}

//泛型的map
func Map(data ,fn interface{}) []interface{} {
	vdata := reflect.ValueOf(data)
	vfn := reflect.ValueOf(fn)

	result := make([]interface{}, vdata.Len())
	for i := 0; i < vdata.Len(); i++ {
		result[i] = vfn.Call([]reflect.Value{vdata.Index(i)})[0].Interface()
	}

	return result
}

//上面的这个Map的问题是没有进行类型检查，所以这里可以手动进行检查

func Transform(slice, function interface{}) (interface{}, error) {
	return transform(slice, function, false)
}

func TransformInPlace(slice, function interface{}) (interface{}, error) {
	return transform(slice, function, true)
}

func transform(slice, function interface{}, inplace bool) (interface{}, error) {
	sliceType := reflect.ValueOf(slice)
	if sliceType.Kind() != reflect.Slice {
		return nil, errors.New("not slice type")
	}

	vfn := reflect.ValueOf(function)
	elemType := sliceType.Type().Elem()
	if !verifySignature(vfn, elemType, nil) {
		return nil, errors.New("func is not the right type")
	}

	sliceOutType := sliceType
	if !inplace {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(vfn.Type().Out(0)), sliceType.Len(), sliceType.Len())
	}

	for i := 0; i < sliceType.Len(); i++ {
		sliceOutType.Index(i).Set(vfn.Call([]reflect.Value{sliceType.Index(i)})[0])
	}

	return sliceOutType.Interface(), nil
}

func verifySignature(fn reflect.Value, types ...reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}

	//检查方法入参和出参是否符合预期
	if fn.Type().NumIn() != len(types)-1 || fn.Type().NumOut() != 1 {
		return false
	}

	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}

	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}

	return true
}

//利用上面的技巧，将reduce和filter函数都改造成泛型的
func GenericReduce(slice, pairFunc, zero interface{}) (interface{}, error) {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		return nil, errors.New("not slice type")
	}

	length := sliceInType.Len()
	if length == 0 {
		return zero, nil
	} else if length == 1 {
		return sliceInType.Index(0), nil
	}

	elemType := sliceInType.Type().Elem()
	vfn := reflect.ValueOf(pairFunc)
	if !verifySignature(vfn, elemType, elemType, elemType) {
		return nil, errors.New("func is not the right type")
	}

	ins := [2]reflect.Value{sliceInType.Index(0), sliceInType.Index(1)}
	out := vfn.Call(ins[:])[0]

	for i := 2; i < length; i++ {
		ins[0], ins[1] = out, sliceInType.Index(i)
		out = vfn.Call(ins[:])[0]
	}

	return out.Interface(), nil
}


func GenericFilter(slice, function interface{}) interface{} {
	result, _ := filter(slice, function, false)
	return result
}

func GenericFilterInPlace(slicePtr, function interface{}) {
	in := reflect.ValueOf(slicePtr)
	if in.Kind() != reflect.Ptr {
		panic("FilterInPlace: wrong type, " +
			"not a pointer to slice")
	}
	_, n := filter(in.Elem().Interface(), function, true)
	in.Elem().SetLen(n)
}

var boolType = reflect.ValueOf(true).Type()

func filter(slice, function interface{}, inPlace bool) (interface{}, int) {

	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("filter: wrong type, not a slice")
	}

	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !verifySignature(fn, elemType, boolType) {
		panic("filter: function must be of type func(" + elemType.String() + ") bool")
	}

	var which []int
	for i := 0; i < sliceInType.Len(); i++ {
		if fn.Call([]reflect.Value{sliceInType.Index(i)})[0].Bool() {
			which = append(which, i)
		}
	}

	out := sliceInType

	if !inPlace {
		out = reflect.MakeSlice(sliceInType.Type(), len(which), len(which))
	}
	for i := range which {
		out.Index(i).Set(sliceInType.Index(which[i]))
	}

	return out.Interface(), len(which)
}

func main() {
	fmt.Println("------------使用MapStrToStr------------")
	arr := []string{"roseduan", "jack zhang", "golang", "24"}
	arr = MapStrToStr(arr, strings.ToUpper)
	fmt.Println(arr)

	fmt.Println("------------使用MapStrToInt------------")
	res := MapStrToInt(arr, func(s string) int {
		return len(s)
	})
	fmt.Println(res)

	//逻辑自定义，例如计算字符串长度之和
	arr2 := []string{"CHN Beijing", "USA New York", "UK London", "CHN Shanghai"}
	sum := Reduce(arr2, func(s string) int {
		return len(s)
	})
	fmt.Println(sum)

	c := Filter(arr2, func(s string) bool {
		return strings.HasPrefix(s, "CHN")
	})
	fmt.Println(c)

	fmt.Println("------------EmployeeCountIf------------")
	count := EmployeeCountIf(list, func(e *Employee) bool {
		return e.Salary > 5000
	})
	fmt.Println("salary more than 5000 : ", count)

	fmt.Println("------------EmployeeFilterIn------------")
	v := EmployeeFilterIn(list, func(e *Employee) bool {
		return e.Vacation > 5 && e.Vacation < 10
	})
	fmt.Println(v)

	fmt.Println("------------EmployeeSumIf------------")
	sum = EmployeeSumIf(list, func(e *Employee) int {
		return e.Age
	})
	fmt.Println(sum)

	fmt.Println("------------使用简单的泛型Map------------")

	//整数数组
	arr3 := []int{1, 2, 3, 4}
	res2 := Map(arr3, func(x int) int {
		return x * x
	})
	fmt.Println(res2)

	//字符串数组
	arr4 := []string{"Java", "Golang", "Python", "Rust"}
	res4 := Map(arr4, strings.ToUpper)
	fmt.Println(res4)

	fmt.Println("------------使用健壮版的泛型Map------------")

	//用于字符串数组
	list := []string{"1", "2", "3", "4", "5"}
	res5, err := Transform(list, func(s string) string {
		return s + s + s
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res5)

	//用于整型数组
	nums := []int{1, 2, 3, 4, 5}
	aa, _ := TransformInPlace(nums, func(n int) int {
		return n * 3
	})
	fmt.Println(nums)
	fmt.Println(aa)

	//用于结构体
	var employeeList = []Employee{
		{"Hao", 44, 4, 8000},
		{"Bob", 34, 10, 5000},
		{"Alice", 23, 5, 9000},
		{"Jack", 26, 3, 4000},
	}

	lis, _ := Transform(employeeList, func(e Employee) Employee {
		e.Salary += 1000
		e.Vacation += 1
		return e
	})
	fmt.Printf("%+v\n", lis)
}
