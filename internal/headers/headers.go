package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	crlf           = "\r\n"
	headerSeperate = ":"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if bytes.HasPrefix(data, []byte(crlf)) {
		return len(crlf), true, nil
	}
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}

	headerLine := string(data[:idx])
	parts := strings.SplitN(headerLine, headerSeperate, 2)
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("invalid header format %s", headerLine)
	}
	rawKey := parts[0]

	if strings.TrimSpace(rawKey) == "" {
		return 0, false, fmt.Errorf("empty key string")
	}
	if strings.Contains(rawKey, " ") {
		return 0, false, fmt.Errorf("invalid header format: space inside key")
	}

	key := strings.TrimSpace(rawKey)
	value := strings.TrimSpace(parts[1])

	h[key] = value
	return idx + len(crlf), false, nil
}

func NewHeaders() Headers {
	return make(map[string]string)
}
