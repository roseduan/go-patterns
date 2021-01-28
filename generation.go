package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

//一个简单的容器，支持任意数据类型
type Container []interface{}

func (c *Container) Put(val interface{}) {
	*c = append(*c, val)
}

func (c *Container) Get() interface{} {
	res := (*c)[0]
	*c = (*c)[1:]
	return res
}

//使用反射进行自动类型转换的Container
type MyContainer struct {
	s reflect.Value
}

func NewContainer(t reflect.Type, size int) *MyContainer {
	if size <= 0 {
		size = 64
	}

	return &MyContainer{
		s: reflect.MakeSlice(reflect.SliceOf(t), 0, size),
	}
}

func (c *MyContainer) MyPut(val interface{}) error {
	//类型检查
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		return errors.New(fmt.Sprintf("Put: can`t put a %T into a slice of %s",
			val, c.s.Type().Elem()))
	}

	c.s = reflect.Append(c.s, reflect.ValueOf(val))
	return nil
}

func (c *MyContainer) MyGet(res interface{}) error {
	v := reflect.ValueOf(res)
	if v.Kind() != reflect.Ptr || v.Elem().Type() != c.s.Type().Elem() {
		return errors.New(fmt.Sprintf("Get: needs *%s but got %T", c.s.Type().Elem(), v))
	}

	v.Elem().Set(c.s.Index(0))
	c.s = c.s.Slice(1, c.s.Len())
	return nil
}

func main() {
	//Container容器的使用
	fmt.Println("------------Container容器的使用------------")
	c := &Container{}
	c.Put(1)
	c.Put("roseduan")
	c.Put("good")
	c.Put(1.4542)

	fmt.Printf("%+v\n", c.Get())
	fmt.Printf("%+v\n", c.Get())
	fmt.Printf("%+v\n", c.Get())

	//取出数据之后，需要进行数据类型的转换
	e, ok := c.Get().(float64)
	if !ok {
		log.Println("the value is not float64")
	} else {
		fmt.Println(e)
	}

	fmt.Println("------------My Container容器的使用------------")
	container := NewContainer(reflect.TypeOf(1), 16)
	err := container.MyPut(112)
	if err != nil {
		log.Println(err)
	}

	_ = container.MyPut(1122)
	_ = container.MyPut(23200)

	var r1 int
	if err := container.MyGet(&r1); err != nil {
		log.Println(err)
	}

	fmt.Println("r1 = ", r1)

	var r2 int
	_ = container.MyGet(&r2)
	fmt.Println("r2 = ", r2)
}
