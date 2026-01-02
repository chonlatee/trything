package main

import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var _, _ fmt.Stringer = (*Request)(nil), (*Response)(nil)
var _, _ encoding.TextMarshaler = (*Request)(nil), (*Response)(nil)

var (
	host, path, method string
	port               int
)

type Header struct{ key, value string }

type Request struct {
	Method  string
	Path    string
	Headers []Header
	Body    string
}

type Response struct {
	StatusCode int
	Headers    []Header
	Body       string
}

func main() {

	flag.StringVar(&method, "method", "GET", "HTTP method to use")
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.IntVar(&port, "port", 8080, "port to connect to")
	flag.StringVar(&path, "path", "/", "part to request")
	flag.Parse()

	ip, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, ip)
	if err != nil {
		panic(err)
	}

	log.Printf("connected to %s (@ %s)", host, conn.RemoteAddr())

	defer conn.Close()

	var reqFields = []string{
		fmt.Sprintf("%s %s HTTP/1.1", method, path),
		"Host: " + host,
		"User-Agent: httpget",
		"",
	}

	request := strings.Join(reqFields, "\r\n") + "\r\n"
	conn.Write([]byte(request))

	log.Printf("sent request:\n%s", request)

	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		line := scanner.Bytes()
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", line); err != nil {
			log.Printf("error writing to connection: %s", err)
		}

		if scanner.Err() != nil {
			log.Printf("error reading from connectin: %s", err)
			return
		}
	}

}

func NewRequest(method, path, host, body string) (*Request, error) {
	switch {
	case method == "":
		return nil, errors.New("missing required argument: method")
	case path == "":
		return nil, errors.New("missing required argument: path")
	case !strings.HasPrefix(path, "/"):
		return nil, errors.New("path must start with /")
	case host == "":
		return nil, errors.New("missing required argument: host")
	default:
		headers := make([]Header, 2)
		headers[0] = Header{"Host", host}
		if body != "" {
			headers = append(headers, Header{"Content-Length", fmt.Sprintf("%d", len(body))})
		}

		return &Request{Method: method, Path: path, Headers: headers, Body: body}, nil

	}
}

func NewResponse(status int, body string) (*Response, error) {
	switch {
	case status < 100 || status > 599:
		return nil, errors.New("invalid status code")
	default:
		if body == "" {
			body = http.StatusText(status)
		}
		headers := []Header{
			{"Content-Length", fmt.Sprintf("%d", len(body))},
		}
		return &Response{StatusCode: status, Headers: headers, Body: body}, nil
	}
}

func (resp *Response) WithHeader(key, value string) *Response {
	resp.Headers = append(resp.Headers, Header{AsTitle(key), value})
	return resp
}

func (r *Request) WithHeader(key, value string) *Request {
	r.Headers = append(r.Headers, Header{AsTitle(key), value})
	return r
}

func AsTitle(key string) string {
	if key == "" {
		panic("empty header key")
	}

	if isTitleCase(key) {
		return key
	}

	return newTitleCase(key)
}

func newTitleCase(key string) string {
	var b strings.Builder
	b.Grow(len(key))

	for i := range key {
		if i == 0 || key[i-1] == '-' {
			b.WriteByte(upper(key[i]))
		} else {
			b.WriteByte(lower(key[i]))
		}
	}

	return b.String()
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c + 'A' - 'a'
	}

	return c
}

func isTitleCase(key string) bool {
	for i := range key {
		if i == 0 || key[i-1] == '-' {
			if key[i] >= 'a' && key[i] <= 'z' {
				return false
			}
		} else if key[i] >= 'A' && key[i] <= 'Z' {
			return false
		}
	}

	return true
}

func (r *Request) WriteTo(w io.Writer) (n int64, err error) {
	printf := func(format string, args ...any) error {
		m, err := fmt.Fprintf(w, format, args...)
		n += int64(m)
		return err
	}

	if err := printf("%s %s HTTP/1.1\r\n", r.Method, r.Path); err != nil {
		return n, err
	}

	for _, h := range r.Headers {
		if err := printf("%s: %s\r\n", h.key, h.value); err != nil {
			return n, err
		}
	}

	printf("\r\n")

	if err := printf("%s\r\n", r.Body); err != nil {
		return n, err
	}

	return n, err
}

func (resp *Response) WriteTo(w io.Writer) (n int64, err error) {
	printf := func(format string, args ...any) error {
		m, err := fmt.Fprintf(w, format, args...)
		n += int64(m)
		return err
	}

	if err := printf("HTTP/1.1 %d %s\r\n", resp.StatusCode, http.StatusText(resp.StatusCode)); err != nil {
		return n, err
	}

	for _, h := range resp.Headers {
		if err := printf("%s: %s\r\n", h.key, h.value); err != nil {
			return n, err
		}
	}

	if err := printf("\r\n%s\r\n", resp.Body); err != nil {
		return n, err
	}

	return n, nil
}

func (r *Request) String() string {
	b := new(strings.Builder)
	r.WriteTo(b)
	return b.String()
}

func (resp *Response) String() string {
	b := new(strings.Builder)
	resp.WriteTo(b)
	return b.String()
}

func (r *Request) MarshalText() ([]byte, error) {
	b := new(bytes.Buffer)
	r.WriteTo(b)
	return b.Bytes(), nil
}

func (resp *Response) MarshalText() ([]byte, error) {
	b := new(bytes.Buffer)
	resp.WriteTo(b)
	return b.Bytes(), nil
}

func ParseRequest(raw string) (r Request, err error) {
	lines := splitLines(raw)

	log.Println(lines)

	if len(lines) < 3 {
		return Request{}, fmt.Errorf("malformed request: should have at least 3 lines")
	}

	first := strings.Fields(lines[0])
	r.Method, r.Path = first[0], first[1]
	if !strings.HasPrefix(r.Path, "/") {
		return Request{}, fmt.Errorf("malformed request: path should start with /")
	}
	if !strings.Contains(first[2], "HTTP") {
		return Request{}, fmt.Errorf("malformed request: first line should contain HTTP version")
	}

	var foundHost bool
	var bodyStart int

	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			bodyStart = i + 1
			break
		}
		key, val, ok := strings.Cut(lines[i], ": ")
		if !ok {
			return Request{}, fmt.Errorf("malformed request: header %q should be of form 'key: value", lines[i])
		}

		if key == "Host" {
			foundHost = true
		}

		key = AsTitle(key)

		r.Headers = append(r.Headers, Header{key, val})
	}

	end := len(lines) - 1
	r.Body = strings.Join(lines[bodyStart:end], "\r\n")

	if !foundHost {
		return Request{}, fmt.Errorf("malformed request: missing Host header")
	}

	return r, nil
}

func ParseResponse(raw string) (resp *Response, err error) {
	lines := splitLines(raw)
	log.Println(lines)

	first := strings.SplitN(lines[0], " ", 3)

	fmt.Printf("first 0: %q\n", first[0])
	if !strings.Contains(first[0], "HTTP") {
		return nil, fmt.Errorf("malformed response: first line should contain HTTP version")
	}

	resp = new(Response)
	resp.StatusCode, err = strconv.Atoi(first[1])
	if err != nil {
		return nil, fmt.Errorf("malformed response: expected status code to be an integer, got %s", first[1])
	}

	if first[2] == "" || http.StatusText(resp.StatusCode) != first[2] {
		log.Printf("missing or incorrect status text for status code %d: expected %q, but got %q", resp.StatusCode, http.StatusText(resp.StatusCode), first[2])
	}

	var bodyStart int

	for i := 1; i < len(lines); i++ {
		log.Println(i, lines[i])
		if lines[i] == "" {
			bodyStart = i + 1
			break
		}
		key, val, ok := strings.Cut(lines[i], ": ")
		if !ok {
			return nil, fmt.Errorf("malformed response: header %q should be of form 'key: value'", lines[i])
		}
		key = AsTitle(key)
		resp.Headers = append(resp.Headers, Header{key, val})
	}
	resp.Body = strings.TrimSpace(strings.Join(lines[bodyStart:], "\r\n"))
	return resp, nil

}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}

	var lines []string
	i := 0
	for {
		j := strings.Index(s[i:], "\r\n")
		if j == -1 {
			lines = append(lines, s[i:])
			return lines
		}
		lines = append(lines, s[i:i+j])
		i += j + 2 // skip the \r\n
	}
}
