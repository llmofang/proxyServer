package main

import (
	"crypto/tls"
	"fmt"
	//"flag"
	"io"
	//"net"
	"net/http"
	"github.com/amahi/spdy"
)
var tlsConfig *tls.Config
type Handler struct {
	client *spdy.Client
}
func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("%v\n",r)
	fmt.Println("%v\n",r.URL)
	fmt.Println("%v\n",r.Method)
	fmt.Println("%v\n",r.Header)
	if r.Method == "CONNECT" {
		//req, err := http.NewRequest("CONNECT", "https:"+r.URL.String(),nil)
		//handle(err)
		//h.client.Do(req)
	}else{
		req, err := http.NewRequest(r.Method,r.URL.String(),nil)
		handle(err)
		for headerKey := range r.Header{
			headerVal := r.Header.Get(headerKey)
			req.Header.Set(headerKey, headerVal)
		}
		resp, err := h.client.Do(req)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer resp.Body.Close()
		for headerKey := range resp.Header {
			headerVal := resp.Header.Get(headerKey)
			w.Header().Set(headerKey, headerVal)
		}
		io.Copy(w, resp.Body)
	}
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	handler := NewHandler()
	cert, err := tls.LoadX509KeyPair("../cert/clientTLS/client.pem", "../cert/clientTLS/client.key")
	if err != nil {
		fmt.Printf("server: loadkeys: %s", err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true, NextProtos: []string{"spdy/3"}}
	conn, err := tls.Dial("tcp", "127.0.0.1:8080", tlsConfig)
	if err != nil {
		fmt.Printf("client: dial: %s", err)
	}

	client, err := spdy.NewClientConn(conn)
	handle(err)	
	if client == nil{

	}
	handler.client = client
	http.ListenAndServe(":8181", handler)
}