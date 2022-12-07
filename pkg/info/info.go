package info

import (
	"net"
	"http/pkg/data"
)

type ServerInfo struct {
	ListenConn *net.TCPListener
	Resources  []*data.Data 
}