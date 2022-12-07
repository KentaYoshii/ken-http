package main

import (
	"net/http"
	"crypto/tls"
	"crypto/x509"
	"os"
	"log"
	"io"
	"fmt"
	"flag"
)

func main() {
	addr := flag.String("addr", "localhost:4000", "HTTPS server address")
	certFile := flag.String("certfile", "cert.pem", "trusted CA certificate")
	//Comment the bottom two lines out if this client is connecting to server with tls
	clientCertFile := flag.String("clientcert", "clientcert.pem", "certificate PEM for client")
 	clientKeyFile := flag.String("clientkey", "clientkey.pem", "key PEM for client")
	flag.Parse()
	//Comment the bottom line out if this client is connecting to server with mtls
	clientCert, err := tls.LoadX509KeyPair(*clientCertFile, *clientKeyFile)
	if err != nil {
		log.Fatal(err)
	}
  
	cert, err := os.ReadFile(*certFile)
	if err != nil {
	  log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
	  log.Fatalf("unable to parse cert from %s", *certFile)
	}
  
	client := &http.Client{
	  Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
		  RootCAs: certPool,
		  //Comment unless client connecting to mtls        
		  Certificates: []tls.Certificate{clientCert},
		},
	  },
	}
  
	r, err := client.Get("https://" + *addr)
	if err != nil {
	  log.Fatal(err)
	}
	defer r.Body.Close()
  
	html, err := io.ReadAll(r.Body)
	if err != nil {
	  log.Fatal(err)
	}
	fmt.Printf("%v\n", r.Status)
	fmt.Printf(string(html))
  }