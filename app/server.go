package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	CLRF      = "\r\n"
	OK        = "HTTP/1.1 200 OK" + CLRF + CLRF
	NOT_FOUND = "HTTP/1.1 404 Not Found" + CLRF + CLRF
	EMPTY     = " "
)

type ReqStatusLine struct {
	Method       string
	Path         string
	HTTP_version string
}
type ResponseStatusLine struct {
	Version string
	Status  string
	Ok      string
}
type Response struct {
	statusline string
	headers    string
	body       string
}
type Header struct {
	Key string
	val string
}

type Headers struct {
	header []Header
}

func (h *Headers) to_string() string {
	res := ""
	for _, r := range h.header {
		res += string(r.Key + ": " + r.val + CLRF + CLRF)
	}
	return res
}

func (h *Header) to_string() string {
	return fmt.Sprintf("%s: %s"+CLRF, h.Key, h.val)
}

func (r *ResponseStatusLine) to_string() string {
	return fmt.Sprintf("%s %s %s"+CLRF, r.Version, r.Status, r.Ok)
}

func extract_statusline(Method, Path, Version string) *ReqStatusLine {
	return &ReqStatusLine{Method: Method, Path: Path, HTTP_version: Version}
}

func handle(con net.Conn) {
	buffer := make([]byte, 1024)
	_, err := con.Read(buffer)
	if err != nil {
		fmt.Println("Error reading buffer")
		os.Exit(1)
	}
	req := bytes.Split(buffer, []byte(CLRF))
	statusL := bytes.Split(req[0], []byte(" "))
	line := extract_statusline(string(statusL[0]), string(statusL[1]), string(statusL[2]))
	parsedPath, parsedPathLen := strings.Split(line.Path, "/"), len(
		strings.Split(line.Path, "/")[1],
	)
	resStatusLine := ResponseStatusLine{Version: "HTTP/1.1", Status: "200", Ok: "OK"}
	if parsedPath[0] != "echo" {
		// todo
		resStatusLine.Status = "404"
	}
	if len(parsedPath[0]) == 0 {
		resStatusLine.Status = "200"
	}

	HEADERS := &Headers{header: make([]Header, 2)}
	head1 := Header{Key: "Content-Type", val: "text/plain"}
	head2 := Header{Key: "Content-Length", val: strconv.Itoa(parsedPathLen)}
	HEADERS.header = append(HEADERS.header, head1)
	HEADERS.header = append(HEADERS.header, head2)
	res := &Response{
		statusline: resStatusLine.to_string(),
		headers:    HEADERS.to_string(),
		body:       parsedPath[1] + CLRF + CLRF,
	}
	con.Write([]byte(res.statusline + res.headers + res.body))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	r, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// r.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	handle(r)
}
