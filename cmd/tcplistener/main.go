package main

import (
	"errors"
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

	fmt.Println("Connection has been established on port ", portAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Cannot accept connection: ", err.Error())
			os.Exit(1)
		}

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
				}

				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(b[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}

			currentLineContents += parts[len(parts)-1]
		}
	}()
	return ch
}
