package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CLRF      = "\r\n"
	OK        = "HTTP/1.1 200 OK" + CLRF + CLRF
	NOT_FOUND = "HTTP/1.1 404 Not Found" + CLRF + CLRF
	EMPTY     = " "
)

var Dir string

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

func (r *ResponseStatusLine) to_string() string {
	return fmt.Sprintf("%s %s %s"+CLRF, r.Version, r.Status, r.Ok)
}

func extract_statusline(Method, Path, Version string) *ReqStatusLine {
	return &ReqStatusLine{Method: Method, Path: Path, HTTP_version: Version}
}

func handle(con net.Conn) {
	defer con.Close()
	buffer := make([]byte, 1024)
	_, err := con.Read(buffer)
	if err != nil {
		fmt.Println("Error reading buffer")
		os.Exit(1)
	}
	req := bytes.Split(buffer, []byte(CLRF))
	statusL := bytes.Split(req[0], []byte(" "))
	line := extract_statusline(string(statusL[0]), string(statusL[1]), string(statusL[2]))
	resStatusLine := ResponseStatusLine{Version: "HTTP/1.1", Status: "200", Ok: "OK"}
	if strings.Trim(line.Path, " ") == "/" {
		con.Write([]byte(OK))
		return
	} else if line.Path == "/user-agent" {
		reqHeaders := &Headers{}
		if req[1] != nil {
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

			resStatusLine := ResponseStatusLine{Version: "HTTP/1.1", Status: "200", Ok: "OK"}
			headers := &Headers{}
			var head2 Header
			head1 := Header{Key: "Content-Type", val: "text/plain"}
			lenActual := len(strings.Trim(reqHeaders.header[1].val, " "))
			head2 = Header{Key: "Content-Length", val: strconv.Itoa(lenActual)}
			headers.header = append(headers.header, head1)
			headers.header = append(headers.header, head2)
			res := &Response{
				statusline: resStatusLine.to_string(),
				headers:    headers.to_string(),
				body:       strings.Trim(reqHeaders.header[1].val, " "),
			}
			con.Write([]byte(res.statusline + res.headers + res.body))
			return

		}
	} else if strings.Contains(line.Path, "/echo") {
		_, parsedPathLen := strings.SplitN(line.Path, "/", -1)[1:],
			strings.SplitN(line.Path, "/", -1)[2:]
		HEADERS := &Headers{}
		var head2 Header
		head1 := Header{Key: "Content-Type", val: "text/plain"}
		lenActual := len(strings.Join(parsedPathLen, "/"))
		head2 = Header{Key: "Content-Length", val: strconv.Itoa(lenActual)}
		HEADERS.header = append(HEADERS.header, head1)
		HEADERS.header = append(HEADERS.header, head2)
		res := &Response{
			statusline: resStatusLine.to_string(),
			headers:    HEADERS.to_string(),
			body:       strings.Join(parsedPathLen, "/"),
		}
		con.Write([]byte(res.statusline + res.headers + res.body))
		return

	} else if strings.Contains(line.Path, "/files") {
		if line.Method == "GET" {
			filename := strings.Split(line.Path, "/")[2]
			pathToFile := filepath.Join(Dir, filename)
			_, err := os.Open(pathToFile)
			if errors.Is(err, os.ErrNotExist) {
				con.Write([]byte(NOT_FOUND))
				return
			}

			HEADERS := &Headers{}
			var head2 Header
			head1 := Header{Key: "Content-Type", val: "application/octet-stream"}
			content, _ := os.ReadFile(pathToFile)
			lenActual := len(string(content))
			head2 = Header{Key: "Content-Length", val: strconv.Itoa(lenActual)}
			HEADERS.header = append(HEADERS.header, head1)
			HEADERS.header = append(HEADERS.header, head2)
			res := &Response{
				statusline: resStatusLine.to_string(),
				headers:    HEADERS.to_string(),
				body:       string(content),
			}
			con.Write([]byte(res.statusline + res.headers + res.body))
			return
		} else {
			fmt.Println(req)
		}
	} else {
		con.Write([]byte(NOT_FOUND))
		return
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.StringVar(&Dir, "directory", "", "enter the dir")
	flag.Parse()
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		r, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle(r)
	}
}
