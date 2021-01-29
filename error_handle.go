package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

//错误处理

//Go语言中的资源清理主要使用 defer 关键字
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

//	a bad case of error handle
//	func parse(r io.Reader) (*Point, error) {
//
//		var p Point
//
//		if err := binary.Read(r, binary.BigEndian, &p.Longitude); err != nil {
//			return nil, err
//		}
//		if err := binary.Read(r, binary.BigEndian, &p.Latitude); err != nil {
//			return nil, err
//		}
//		if err := binary.Read(r, binary.BigEndian, &p.Distance); err != nil {
//			return nil, err
//		}
//		if err := binary.Read(r, binary.BigEndian, &p.ElevationGain); err != nil {
//			return nil, err
//		}
//		if err := binary.Read(r, binary.BigEndian, &p.ElevationLoss); err != nil {
//			return nil, err
//		}
//	}

//要处理上面这种 if err != nil 的代码，可以采用下面这种方式
func Parse(r io.Reader, p *Point) (*Point, error) {
	var err error
	read := func(data interface{}) {
		err = binary.Read(r, binary.BigEndian, data)
		if err != nil {
			return
		}
	}

	read(p.Longitude)
	read(p.Latitude)
	read(p.Distance)
	read(p.ElevationGain)
	read(p.ElevationLoss)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//还可以借鉴 bufio.NewScanner 进行改进
type Reader struct {
	r   io.Reader
	err error
}

func (r *Reader) read(data interface{}) {
	if r.err != nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

func newParse(r io.Reader, p *Point) (*Point, error) {
	reader := Reader{r: r}
	reader.read(p.Longitude)
	reader.read(p.Latitude)
	reader.read(p.Distance)
	reader.read(p.ElevationGain)
	reader.read(p.ElevationLoss)

	if reader.err != nil {
		return nil, reader.err
	}

	return p, nil
}

type Point struct {
	Longitude     []byte
	Latitude      []byte
	Distance      []byte
	ElevationGain []byte
	ElevationLoss []byte
}

//利用上面的这个技巧，再来看一个例子
var b = []byte{114, 111, 115, 101, 100, 117, 97, 110, 23, 70}

//这个byte数组的数据可以这样来构造
//	var buf bytes.Buffer
//	err := binary.Write(&buf, binary.BigEndian, []byte("roseduan"))
//	_ = binary.Write(&buf, binary.BigEndian, uint8(23))
//	_ = binary.Write(&buf, binary.BigEndian, uint8(70))
//	if err != nil {
//		log.Println(err)
//	}
//
//	fmt.Println(buf.Bytes())

var r = bytes.NewReader(b)

type MyPerson struct {
	Name   [8]byte
	Age    uint8
	Weight uint8
	err    error
}

func (p *MyPerson) read(data interface{}) {
	if p.err == nil {
		p.err = binary.Read(r, binary.BigEndian, data)
	}
}

func (p *MyPerson) ReadName() *MyPerson {
	p.read(&p.Name)
	return p
}

func (p *MyPerson) ReadAge() *MyPerson {
	p.read(&p.Age)
	return p
}

func (p *MyPerson) ReadWeight() *MyPerson {
	p.read(&p.Weight)
	return p
}

func (p *MyPerson) PrintPersonInfo() {
	fmt.Printf("Name=%s, Age=%d, Weight=%d\n", p.Name, p.Age, p.Weight)
}

func main() {
	file, _ := os.OpenFile("/test/fileA", os.O_RDONLY, os.ModePerm)
	defer Close(file)
	//do something with the opened file

	fmt.Println("------------------")
	p := MyPerson{}
	p.ReadName().ReadAge().ReadWeight()

	if p.err != nil {
		log.Fatal(p.err)
	}

	p.PrintPersonInfo()
}
