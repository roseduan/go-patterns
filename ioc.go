package main

import "fmt"

//委托和控制反转


//Go语言中的结构体可以嵌入
type Widget struct {
	X, Y int
}

type Label struct {
	Widget
	Text string
}

type Button struct {
	Label
}

type ListBox struct {
	Widget
	Texts []string
	Index int
}

//方法重写
type Painter interface {
	Paint()
}

type Clicker interface {
	Click()
}

func (l *Label) Paint() {
	fmt.Printf("Label.Paint(%q)\n", l.Text)
}

func (b *Button) Paint()  {
	fmt.Printf("Button.Paint(%q)\n", b.Text)
}

func (lis *ListBox) Paint() {
	fmt.Printf("ListBox.Paint(%q)\n", lis.Texts)
}

func (b *Button) Click() {
	fmt.Printf("Button.Click(%q)\n", b.Text)
}

func (lis *ListBox) Click() {
	fmt.Printf("ListBox.Click(%q)\n", lis.Texts)
}

//控制反转的例子

//定义一个简单的Set
type IntSet struct {
	data map[int]bool
}

func NewIntSet() IntSet {
	return IntSet{make(map[int]bool)}
}

func (i *IntSet) Add(x int) {
	i.data[x] = true
}

func (i *IntSet) Delete(x int) {
	delete(i.data, x)
}

func (i *IntSet) Contains(x int) bool {
	return i.data[x]
}

//实现一个有Undo功能的IntSet
type UndoableIntSet struct {
	IntSet
	fn []func()
}

func NewUndoableIntSet() *UndoableIntSet {
	return &UndoableIntSet{NewIntSet(), nil}
}

func (u *UndoableIntSet) Add(x int) {
	if !u.Contains(x) {
		u.IntSet.data[x] = true
		u.fn = append(u.fn, func() {
			u.Delete(x)
		})
	} else {
		u.fn = append(u.fn, nil)
	}
}

func (u *UndoableIntSet) Delete(x int) {
	if u.Contains(x) {
		delete(u.IntSet.data, x)
		u.fn = append(u.fn, func() {
			u.Add(x)
		})
	} else {
		u.fn = append(u.fn, nil)
	}
}

func (u *UndoableIntSet) Undo() bool {
	if len(u.fn) == 0 {
		return false
	}

	idx := len(u.fn)-1
	if u.fn[idx] != nil {
		u.fn[idx]()
		u.fn[idx] = nil
	}

	u.fn = u.fn[:idx]
	return true
}

//依赖反转
type Undo []func()

func (u *Undo) Add(fn func()) {
	*u = append(*u, fn)
}

func (u *Undo) Execute() bool {
	if len(*u) == 0 {
		return false
	}

	idx := len(*u)-1
	if fn := (*u)[idx]; fn != nil {
		fn()
		(*u)[idx] = nil
	}

	*u = (*u)[:idx]
	return true
}

type StringSet struct {
	data map[string]bool
	undo Undo
}

func NewStringSet() *StringSet {
	return &StringSet{data: make(map[string]bool)}
}

func (s *StringSet) Add(x string) {
	if !s.Contains(x) {
		s.data[x] = true
		s.undo.Add(func() {
			s.Delete(x)
		})
	} else {
		s.undo.Add(nil)
	}
}

func (s *StringSet) Delete(x string) {
	if s.Contains(x) {
		delete(s.data, x)
		s.undo.Add(func() {
			s.Add(x)
		})
	} else {
		s.undo.Add(nil)
	}
}

func (s *StringSet) Contains(x string) bool {
	return s.data[x]
}

func main() {
	label := Label{Widget{10, 20}, "a label"}
	button1 := Button{Label{Widget{10, 70}, "OK"}}
	button2 := Button{Label{Widget{10, 20}, "Cancel"}}
	listBox := ListBox{Widget{10, 40}, []string{"AL", "AK", "AZ", "AR"}, 0}

	for _, painter := range []Painter{&label, &listBox, &button1, &button2} {
		painter.Paint()
	}

	for _, widget := range []interface{}{&label, &listBox, &button1, &button2} {
		widget.(Painter).Paint()
		if clicker, ok := widget.(Clicker); ok {
			clicker.Click()
		}
		fmt.Println() // print a empty line
	}

	fmt.Println("-----------使用UndoableIntSet-----------")
	set := NewUndoableIntSet()
	set.Add(1)
	set.Add(4)
	set.Add(3)

	set.Undo()
	fmt.Println(set.data)

	set.Undo()
	fmt.Println(set.data)

	set.Undo()
	fmt.Println(set.data)

	fmt.Println("-----------使用StringSet-----------")
	s := NewStringSet()
	s.Add("a")
	s.Add("b")
	s.Add("c")

	fmt.Println(s.data)
	s.undo.Execute()
	s.undo.Execute()
	fmt.Println(s.data)
}
