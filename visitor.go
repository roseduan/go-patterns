package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

//Visitor模式

//一个简单的visitor
type Visitor func(shape Shape)

type Shape interface {
	accept(v Visitor)
}

type Circle struct {
	Radius int
}

type Rectangle struct {
	Width, Height int
}

func (c Circle) accept(v Visitor) {
	v(c)
}

func (r Rectangle) accept(v Visitor) {
	v(r)
}

func JsonVisitor(shape Shape) {
	res, err := json.Marshal(shape)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
}

func XmlVisitor(shape Shape) {
	res, err := xml.Marshal(shape)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
}

//kubectl的实现方法

type VisitorFunc func(*Info, error) error

type InfoVisitor interface {
	Visit(VisitorFunc) error
}

type Info struct {
	Namespace string
	Name string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

type NameVisitor struct {
	visitor InfoVisitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("Name=%s, Namespace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor after call function")
		return err
	})
}

type OtherVisitor struct {
	visitor InfoVisitor
}

func (v OtherVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherVisitor before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherVisitor after call function")
		return err
	})
}

type LogVisitor struct {
	visitor InfoVisitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor after call function")
		return err
	})
}

func main() {
	c := Circle{10}
	r := Rectangle{10, 20}
	shapes := []Shape{c, r}

	for _, s := range shapes {
		s.accept(JsonVisitor)
		s.accept(XmlVisitor)
	}

	info := Info{}
	var v InfoVisitor = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherVisitor{v}

	loadFile := func(info *Info, err error) error {
		info.Name = "roseduan"
		info.Namespace = "roseduan.com"
		info.OtherThings = "I am a nice person"
		return nil
	}

	v.Visit(loadFile)
}
