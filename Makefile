make:
		go build -o http-server cmd/server/main.go
		go build -o cgi-server cmd/cgi/main.go
		go build -o http-client cmd/client/main.go
		go build -o gen-cert cmd/certificate/main.go
		go build -o http-server-tls cmd/server-tls/main.go
		go build -o http-server-mtls cmd/server-mtls/main.go

clean:
		rm -f ./http-server
		rm -f ./http-client
		rm -f ./gen-cert
		rm -f ./http-server-tls
		rm -f ./http-server-mtls
		rm -f ./cgi-server

