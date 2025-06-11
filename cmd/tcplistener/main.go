package main

import (
	"fmt"
	"net"
	"os"

	"github.com/nt2311-vn/HttpProtocol/internal/request"
)

var portAddr string = ":42069"

func main() {
	l, err := net.Listen("tcp", portAddr)
	if err != nil {
		fmt.Println("Cannot listen to port ", portAddr, err.Error())
		os.Exit(1)
	}

	defer l.Close()

	fmt.Println("Connection has been established on port ", portAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Cannot accept connection: ", err.Error())
			os.Exit(1)
		}

		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Cannot parse request: %v", err)
			os.Exit(1)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")

		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(request.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
