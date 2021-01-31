package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

//修饰器模式


//一个简单的修饰器函数
func decorator(fn func(s string)) func (s string) {
	return func(s string) {
		fmt.Println("started")
		fn(s)
		fmt.Println("done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

//一个计算函数运行时间的例子
type SumFunc func(int64, int64) int64

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TimedSumFunc(f SumFunc) SumFunc {
	return func(start int64, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("-------Time Elapsed (%s): %v-------\n", getFuncName(f), time.Since(t))
		}(time.Now())

		return f(start, end)
	}
}

func SumFunc1(start, end int64) int64 {
	var sum int64 = 0
	if start > end {
		start, end = end, start
	}

	for i := start; i < end; i++ {
		sum += i
	}

	return sum
}

func SumFunc2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}

	return (end - start + 1) * (end + start) / 2
}

//一个http的例子
func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("-------with server header-------")
		w.Header().Set("Server", "HelloServer v0.0.1")
		h(w, r)
	}
}

func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("-------with auth cookie-------")
		cookie := http.Cookie{Name: "auth", Value: "pass", Path: "/"}
		http.SetCookie(w, &cookie)
		h(w, r)
	}
}

func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--------with basic auth---------")
		cookie, err := r.Cookie("auth")
		if cookie == nil || cookie.Value != "pass" || err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func WithDebugLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--------with debug log---------")
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(r.Form)
		log.Println("path = ", r.URL.Path)
		log.Println("scheme = ", r.URL.Scheme)
		h(w, r)
	}
}

func ServerHello(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintln(w, "Hello World!" + r.URL.Path)
}

//上面的几个with修饰方法，不太优雅，可以进行下面的改造
type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		decorator := decors[len(decors)-1-i]
		h = decorator(h)
	}
	return h
}

//泛型的修饰器
func Decorator(decoPtr, fn interface{}) (err error) {
	decoFunc := reflect.ValueOf(decoPtr).Elem()
	targetFunc := reflect.ValueOf(fn)

	v := reflect.MakeFunc(targetFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		fmt.Println("before")
		out = targetFunc.Call(in)
		fmt.Println("after")
		return
	})

	decoFunc.Set(v)
	return
}

func main() {
	h := decorator(Hello)
	h("Hello, I am roseduan")

	fn1 := TimedSumFunc(SumFunc1)
	res1 := fn1(1, 10000000)
	fmt.Println(res1)

	fn2 := TimedSumFunc(SumFunc2)
	res2 := fn2(1, 10000000)
	fmt.Println(res2)

	fmt.Println("----------运行HelloServer-----------")
	//http.HandleFunc("/v1/hello", WithServerHeader(ServerHello))
	//http.HandleFunc("/v2/hello", WithAuthCookie(ServerHello))
	//http.HandleFunc("/v3/hello", WithBasicAuth(ServerHello))
	//http.HandleFunc("/v4/hello", WithDebugLog(ServerHello))

	http.HandleFunc("/v1/hello", Handler(ServerHello, WithAuthCookie, WithServerHeader, WithBasicAuth))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
