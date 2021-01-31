package main

import "fmt"

//Pipeline 模式

//一个简单的示例
func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

//平方函数
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out
}

//过滤奇数函数
func odd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n % 2 == 1 {
				out <- n * n
			}
		}
		close(out)
	}()

	return out
}

//求和函数
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		sum := 0
		for n := range in {
			sum += n
		}
		out <- sum
		close(out)
	}()
	return out
}

//func EchoFunc(in <-chan int, fn func(in <- chan int, out chan int)) <-chan int {
//	out := make(chan int)
//	go fn(in, out)
//	return out
//}

//上面几个函数的多层嵌套使用，可以使用一个代理函数解决
type EchoFunc func([]int) <- chan int
type PipeFunc func(in <-chan int) <-chan int

func pipeline(nums []int, echoFunc EchoFunc, pipeFunc ...PipeFunc) <-chan int {
	ch := echoFunc(nums)
	for i := range pipeFunc {
		ch = pipeFunc[i](ch)
	}
	return ch
}

func main() {
	////简单的pipeline的使用方式
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	for n := range square(echo(nums)) {
		fmt.Println(n)
	}

	//另一种更简洁的方式
	for n := range pipeline(nums, echo, square, odd, sum) {
		fmt.Println(n)
	}
}
