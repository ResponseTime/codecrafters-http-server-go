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
	OK        = "HTTP/1.1 200 OK" + CLRF
	NOT_FOUND = "HTTP/1.1 404 Not Found" + CLRF
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
		res += string(r.Key + ": " + r.val + CLRF)
	}
	res += CLRF
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
	reqHeaders := &Headers{}
	for _, i := range req[1:] {
		if len(string(i)) == 0 {
			break
		}
		splits := strings.Split(string(i), ":")
		reqHeaders.header = append(
			reqHeaders.header,
			Header{Key: strings.Trim(splits[0], ""), val: strings.Trim(splits[1], "")},
		)
	}
	statusL := bytes.Split(req[0], []byte(" "))
	// line := extract_statusline(string(statusL[0]), string(statusL[1]), string(statusL[2]))
	// parsedPath, parsedPathLen := strings.SplitN(line.Path, "/", -1)[1:],
	// 	strings.SplitN(line.Path, "/", -1)[2:]
	resStatusLine := ResponseStatusLine{Version: "HTTP/1.1", Status: "200", Ok: "OK"}
	echo := false
	if bytes.Contains(req[0], []byte("/user-agent")) || bytes.Contains(req[0], []byte("/ ")) {
		// todo

		fmt.Println(string(req[0]))
		resStatusLine.Status = "200"
	}
	if bytes.Contains(req[0], []byte("/echo")) {
		echo = true
	} else {
		resStatusLine.Status = "404"
	}

	HEADERS := &Headers{}
	// lenActual := strings.Join(parsedPathLen, "/")

	var head2 Header
	head1 := Header{Key: "Content-Type", val: "text/plain"}
	head2 = Header{Key: "Content-Length", val: strconv.Itoa(0)}
	var res *Response

	if resStatusLine.Status == "404" {
		HEADERS.header = append(HEADERS.header, head1)
		HEADERS.header = append(HEADERS.header, head2)
		res = &Response{
			statusline: resStatusLine.to_string(),
			headers:    HEADERS.to_string(),
			body:       "",
		}
	}
	if echo {
		lenActual := len(strings.TrimPrefix(string(statusL[1]), "/echo/"))
		head2 = Header{Key: "Content-Length", val: strconv.Itoa(lenActual)}
		HEADERS.header = append(HEADERS.header, head1)
		HEADERS.header = append(HEADERS.header, head2)
		res = &Response{
			statusline: resStatusLine.to_string(),
			headers:    HEADERS.to_string(),
			body:       strings.TrimPrefix(string(statusL[1]), "/echo/"),
		}
	} else {
		lenActual := len(reqHeaders.header[1].val)

		head2 = Header{Key: "Content-Length", val: strconv.Itoa(lenActual)}
		HEADERS.header = append(HEADERS.header, head1)
		HEADERS.header = append(HEADERS.header, head2)
		res = &Response{
			statusline: resStatusLine.to_string(),
			headers:    HEADERS.to_string(),
			body:       reqHeaders.header[1].val,
		}
	}
	// fmt.Println(res.body, len(res.body))
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
	defer r.Close()
	// r.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	handle(r)
}
