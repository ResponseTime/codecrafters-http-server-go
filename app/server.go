package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const (
	CLRF      = "\r\n"
	OK        = "HTTP/1.1 200 OK" + CLRF + CLRF
	NOT_FOUND = "HTTP/1.1 404 Not Found" + CLRF + CLRF
)

type ReqStatusLine struct {
	Method       string
	Path         string
	HTTP_version string
}

func extract_statusline(Method, Path, Version string) *ReqStatusLine {
	return &ReqStatusLine{Method: Method, Path: Path, HTTP_version: Version}
}

func handle(con net.Conn) {
	buffer := make([]byte, 1024)
	con.Read(buffer)
	req := bytes.Split(buffer, []byte(CLRF))
	statusL := bytes.Split(req[0], []byte(" "))
	line := extract_statusline(string(statusL[0]), string(statusL[1]), string(statusL[2]))

	if line.Path == "/" {
		con.Write([]byte(OK))
	} else {
		con.Write([]byte(NOT_FOUND))
	}
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

	r.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	handle(r)
}
