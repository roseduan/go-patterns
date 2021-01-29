package main

import (
	"fmt"
	"time"
)

//Functional Options

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConn  int
	Tls      *TlsConfig
}

type TlsConfig struct {
	Cert       []byte
	PrivateKey []byte
}

//由于Go语言不支持函数重载，所以需要写不同的方法来满足不同的需求
func NewDefaultServer(addr string, port int) *Server {
	return &Server{addr, port, "tcp", 10 * time.Second, 10, nil}
}

func NewTlsServer(addr string, port int, cfg *TlsConfig) *Server {
	return &Server{addr, port, "tcp", 10 * time.Second, 10, cfg}
}

func NewServerWithTimeout(addr string, port int, timeout time.Duration) *Server {
	return &Server{addr, port, "tcp", timeout, 100, nil}
}

//要想解决这种代码冗余的问题，一种常见的方式是使用配置对象

//将可更改的配置属性放到一个结构体中
type ServerConfig struct {
	Protocol string
	Timeout  time.Duration
	MaxConn  int
	Tls      TlsConfig
}

//然后server结构体和创建server的方法就可以这样写了
//type Server struct {
//	Addr     string
//	Port     int
//	cfg *ServerConfig
//}
//
//func NewServer(addr string, port int, cfg ServerConfig) *Server {
//	return &Server{addr, port, cfg}
//}

//另一种解决方式：使用builder模式
type ServerBuilder struct {
	Server
}

func (sb *ServerBuilder) New(addr string, port int) *ServerBuilder {
	sb.Server.Addr = addr
	sb.Server.Port = port
	return sb
}

func (sb *ServerBuilder) Protocol(p string) *ServerBuilder {
	sb.Server.Protocol = p
	return sb
}

func (sb *ServerBuilder) Timeout(timeout time.Duration) *ServerBuilder {
	sb.Server.Timeout = timeout
	return sb
}

func (sb *ServerBuilder) MaxConn(maxConn int) *ServerBuilder {
	sb.Server.MaxConn = maxConn
	return sb
}

func (sb *ServerBuilder) Tls(cfg *TlsConfig) *ServerBuilder {
	sb.Server.Tls = cfg
	return sb
}

func (sb *ServerBuilder) Build() Server {
	return sb.Server
}

//builder模式的使用示例
//	builder := ServerBuilder{}
//	s := builder.New("https://roseduan.com", 80).Protocol("http").Timeout(5 * time.Second).Build()
//	fmt.Println(s)

//一种更加优雅的处理方式：functional options

type Option func(*Server)

func Protocol(p string) Option {
	return func(server *Server) {
		server.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.Timeout = timeout
	}
}

func MaxConn(maxConn int) Option {
	return func(server *Server) {
		server.MaxConn = maxConn
	}
}

func Tls(tls *TlsConfig) Option {
	return func(server *Server) {
		server.Tls = tls
	}
}

//构造方法可以这样写了
func NewServer(addr string, port int, options ...Option) *Server {
	server := &Server{Addr: addr, Port: port}
	for _, opt := range options {
		opt(server)
	}

	return server
}

func main() {
	server := NewDefaultServer("https://roseduan.com", 80)
	fmt.Println(server)

	builder := ServerBuilder{}
	s := builder.New("https://roseduan.com", 80).Protocol("http").Timeout(5 * time.Second).Build()
	fmt.Println(s)

	s1 := NewServer("https://roseduan.com", 80)
	s2 := NewServer("https://roseduan.com", 80, Protocol("http"))
	s3 := NewServer("https://roseduan.com", 80, Timeout(5 * time.Second), MaxConn(10))
	fmt.Printf("%+v\n", s1)
	fmt.Printf("%+v\n", s2)
	fmt.Printf("%+v\n", s3)
}
