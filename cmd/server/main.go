package main

import (
	"fmt"
	"os"
	"net"
	"http/pkg/protocol"
	"http/pkg/data"
	"bufio"
	"strings"
	"http/pkg/info"
	"path/filepath"
)

var sessionId = 0

func newListener(service string) *net.TCPListener {
	port_str := fmt.Sprintf(":%s", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port_str)
	protocol.CheckError(err, "tcp4")
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	protocol.CheckError(err, "tcp4")
	fmt.Printf("Server running at port %s\n", service)
	return listener;
}

func serverCLI(sClI_chan chan []string){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		curLine := strings.Trim(scanner.Text(), " ")
		token := strings.Split(curLine, " ")
		// if it is not p or q
		if token[0] != "p" && token[0] != "q" {
			fmt.Println("'p' to print out the list of resources and 'q' for a graceful exit")
		} else {
			sClI_chan <- token
		}
	}
}

func acceptLoop(nodeInfo *info.ServerInfo) {
	listener := nodeInfo.ListenConn
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		sid := sessionId
		new_client := &protocol.Client{
			Conn_Socket: conn,
			SessionId: sid,
			ServerCloseChan: make(chan bool),
		}
		fmt.Printf("[New Client] Client %d connected\n> ", new_client.SessionId)
		sessionId += 1
		go handleClient(new_client, nodeInfo)
	}
}

func handleClient(client *protocol.Client, serverInfo *info.ServerInfo) {
	req, err := protocol.ParseRequest(client.Conn_Socket)
	if err != nil {
		if err.Error() == "invalid request line" {
			return
		} else {
			protocol.CheckError(err, "parse request")
		}
	}
	
	url := req.Uri
	var page *data.Data
	//landing page
	if url == "/" {
		for _, resource := range serverInfo.Resources {
			if resource.Name == "index.html" {
				page = resource
			}
		}
		if page == nil {
			fmt.Println("no index.html file found")
			return 
		}
		bytesToSend := protocol.FormOKResponse(page)
		client.Conn_Socket.Write(bytesToSend)
	} else {
		var bytesToSend []byte
		urlName := url[1:]
		for _, resource := range serverInfo.Resources {
			if urlName == resource.Name {
				page = resource
				bytesToSend = resource.ContentBytes
				break
			}
		}
		if bytesToSend == nil {
			fmt.Println("invalid url")
			return
		}
		bytesToSend = protocol.FormOKResponse(page)
		client.Conn_Socket.Write(bytesToSend)
	}
}

func populateData(pathDir string) []*data.Data {
	dataSlice := make([]*data.Data, 0)
	err := filepath.Walk(pathDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("[Error]: error reading file %s\n", path)
            os.Exit(1)
        }
		if info.IsDir() {
			return nil
		}
		dat, err := os.ReadFile(path)
		protocol.CheckError(err, "read file")
		newData := &data.Data{
			Name:	info.Name(),
			Path:   path,
			ContentBytes: dat,
			IsDeleted: false,
		}
		dataSlice = append(dataSlice, newData)
		return nil
    })

    if err != nil {
        fmt.Println(err)
    }
	fmt.Println("Data populated")
	return dataSlice
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("[Error]: The format should be ./http-server <path to resources>")
		os.Exit(1)
	} 

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Println("[Error]: Input resource path does not exist")
		os.Exit(1)
	}

	dataSlice := populateData(os.Args[1])

	//default port for http
	listen_port := "20080" 
	listener := newListener(listen_port)

	nodeInfo := &info.ServerInfo{
		ListenConn: listener,
		Resources: dataSlice,
	}

	sCLI_chan := make(chan []string)


	go acceptLoop(nodeInfo)
	go serverCLI(sCLI_chan)

	for cmd := range sCLI_chan {
		if cmd[0] == "p" {
			fmt.Println(data.PrintResources(nodeInfo.Resources))
			fmt.Print("> ")
		} else if cmd[0] == "q" {
			fmt.Println("Server Performing a graceful exit")
			//terminateConns()
			fmt.Println("All clients disconnected")
			os.Exit(0)
		} 
	
	}
}