package main

import (
	"net/http"
	"net/http/cgi"
)

func cgiHandler(w http.ResponseWriter, r *http.Request) {
	handler := cgi.Handler{Path: "/home/cs1680-user/final-project-kenta/cmd/cgi/cgichild.sh"}
	handler.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", cgiHandler)
	http.ListenAndServe("localhost:8080", nil)
}