package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println("Cannot read messages.txt file on root: ", err.Error())
		os.Exit(1)
	}

	for line := range getLinesChannel(f) {
		fmt.Printf("read: %s\n", line)
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
