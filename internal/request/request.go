package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type ParserState int

const (
	Initialized ParserState = iota
	Done
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		State: Initialized,
	}

	for req.State != Done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		bytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = Done
				break
			}
			return nil, err
		}

		readToIndex += bytesRead

		bytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[bytesParsed:])
		readToIndex -= bytesParsed
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = Done
		return n, nil

	case Done:
		return 0, fmt.Errorf("error: trying to read data in done state")

	default:
		return 0, fmt.Errorf("unknown state")
	}
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

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil
}
