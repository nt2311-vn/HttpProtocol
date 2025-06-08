package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

const (
	crlf           = "\r\n"
	headerSeperate = ":"
)

type Headers map[string]string

var specialChars = []byte{
	'!',
	'#',
	'$',
	'%',
	'&',
	'\'',
	'*',
	'+',
	'-',
	'.',
	'^',
	'_',
	'`',
	'|',
	'~',
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
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

	if !validToken([]byte(key)) {
		return 0, false, fmt.Errorf("invalid token header found: %s", key)
	}
	h.Set(key, value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func NewHeaders() Headers {
	return make(map[string]string)
}

func validToken(token []byte) bool {
	for _, char := range token {
		if slices.Contains(specialChars, char) || unicode.IsDigit(rune(char)) ||
			unicode.IsLetter(rune(char)) {
			return true
		}
	}

	return false
}
