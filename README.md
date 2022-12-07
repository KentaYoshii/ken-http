# ken-http

For better understanding [see here](./cs1680_fp_summary.pdf)

## Overview
A webserver implemented without net/http Golang package that implements subset of HTTP 1.1. We use sockets to establish each connection with the client (which could be a web browser or anything.)

## Structure
tl;dr:
- Under cmd are all the main files.
- Under pkg we have all the helper functions
- Under resources we have all the resources we are serving on this webserver
```
nima-server
│   README.md
│   Makefile
|   http-client (after make)
|   http-server (after make)    
│
└───cmd
│   │  
│   └───client
│       │   main.go
│   |___server
|       |   main.go
|   |___cgi
|       |   main.go
│   
└───pkg
|   │  
│   └───data
│       │   data.go
│   |___handle
|       |   handle.go
|   |___info
|       |   info.go
│   └───parse
│       │   parse.go
│   |___protocol
|       |   protocol.go
└───resources
│   │  
│   └───test_dependency
│       │   ...
│   |___test_visual
|       |   ...
|   |___test_single
|       |   ...
```
## Getting started
```
git clone <link to this repo>
```

Then, place the resources that you want to serve on this server under one folder and put it in `resources/`

To start the server with the resources you just set, execute the following command from the root

```
./http-server ./resources/<your folder name>
```
