package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println("Cannot read messages.txt file on root: ", err.Error())
		os.Exit(1)
	}

	defer f.Close()

	readBuffer := make([]byte, 8)

	for {
		_, err := f.Read(readBuffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("cannot read streaming file: ", err.Error())
		}

		fmt.Printf("read: %s\n", string(readBuffer))
	}
}
