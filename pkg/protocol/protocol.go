package protocol

import (
	"net"
	"log"
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"http/pkg/parse"
	"http/pkg/data"
)

type request struct {
	Method string // GET, POST, etc.
	Header parse.MIMEHeader
	Body   []byte
	Uri    string // The raw URI from the request
	Proto  string // "HTTP/1.1"
}

type Client struct {
	SessionId int
	Conn_Socket *net.TCPConn
	ServerCloseChan chan bool
}

func CheckError(err error, msg string){
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

func FormOKResponse(resource *data.Data) []byte {
	header := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html;" +
		"charset=utf-8\r\n" +
		"Content-Length: %d\r\n"+
		"\r\n", len(resource.ContentBytes))
	headerBytes := []byte(header)
	bytesToSend := make([]byte, 0)
	bytesToSend = append(bytesToSend, headerBytes...)
	bytesToSend = append(bytesToSend, resource.ContentBytes...)
	return bytesToSend
}

func ParseRequest(c *net.TCPConn) (*request, error) {
	b := bufio.NewReader(c)
	tp := parse.NewReader(b) // need replace
	req := new(request)
	// Parse request line: parse "GET /index.html HTTP/1.0"
	var s string
	s, _ = tp.ReadLine() // need replace
	sp := strings.Split(s, " ")
	if len(sp) < 3 {
		return nil, errors.New("invalid request line")
	}
	req.Method, req.Uri, req.Proto = sp[0], sp[1], sp[2]
	// Parse request headers
	mimeHeader, _ := tp.ReadMIMEHeader() // need replace
	req.Header = mimeHeader
	// Parse request body
	if req.Method == "GET" || req.Method == "HEAD" {
		return req, nil
	}
	if len(req.Header["Content-Length"]) == 0 {
		return nil, errors.New("no content length")
	}
	length, err := strconv.Atoi(req.Header["Content-Length"][0])
	if err != nil {
		return nil, err
	}
	body := make([]byte, length)
	if _, err = io.ReadFull(b, body); err != nil {
		return nil, err
	}
	req.Body = body
	return req, nil
}