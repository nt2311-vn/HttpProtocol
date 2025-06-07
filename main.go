package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var portAddr string = ":42069"

func main() {
	l, err := net.Listen("tcp", portAddr)
	if err != nil {
		fmt.Println("Cannot listen to port ", portAddr, err.Error())
		os.Exit(1)
	}

	defer l.Close()
	println("Connection has been established on port ", portAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Cannot accept connection: ", err.Error())
			os.Exit(1)
		}

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		readBuffer := make([]byte, 8)
		fileStr := ""
		for {
			n, err := f.Read(readBuffer)
			if err == io.EOF {
				fileStr = strings.TrimSuffix(fileStr, "\n")
				for part := range strings.SplitSeq(fileStr, "\n") {
					if part != "" {
						ch <- part
					}
				}
				close(ch)
				return
			}
			if err != nil {
				fmt.Println("cannot read streaming file: ", err.Error())
				close(ch)
				return
			}
			fileStr += string(readBuffer[:n])
		}
	}()
	return ch
}
