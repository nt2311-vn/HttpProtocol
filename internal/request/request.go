package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read request line %v", err)
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("not valid format %s", str)
	}

	method, host, httpVersion := parts[0], parts[1], parts[2]

	for _, char := range method {
		if !unicode.IsLetter(char) || !unicode.IsUpper(char) {
			return nil, errors.New("invalid method")
		}
	}

	httpInfo := strings.Split(httpVersion, "/")
	if len(httpInfo) != 2 {
		return nil, errors.New("invalid httpInfo format")
	}

	version := httpInfo[1]
	if version != "1.1" {
		return nil, errors.New("invalid support HTTP version")
	}

	return &RequestLine{
		Method:        method,
		HttpVersion:   version,
		RequestTarget: host,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, errors.New("cannot get request line from request")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}
