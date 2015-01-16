package main

import (
	"crypto/tls"
	"fmt"
	
	"io"
	"net"
	"net/http"
	"github.com/amahi/spdy"
)

type Handler struct {
	client *spdy.Client
	conn *tls.Conn
}
func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("%v", r.Header)
	var url string
	if r.Method == "CONNECT" {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		url ="https:"+r.URL.String()
		fmt.Println("%+v\n",r)
		fmt.Println("%v\n",url)
		fmt.Println("%v\n",r.Host)
		fmt.Println("%+v\n",r.Header)
/*
		req, err := http.NewRequest(r.Method,url,nil) 
		handle(err)
		res, err := h.client.Do(req)
		handle(err)
		if res ==nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}	

		/*
		serverConn, err := net.Dial("tcp",r.Host)
		if err != nil {
			return
		}
		fmt.Println("%+v\n",serverConn)
		*/
		fmt.Println("datil new conn")
		serverConn, err := net.Dial("tcp","127.0.0.1:8080")
		if err != nil {
			panic(err)
			return
		}
		
		conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n" +
		"Content-Type: text/html\r\n" +
		"Content-Length: 200\r\n" +
		"\r\n"));	
		go io.Copy(serverConn, conn)
		go io.Copy(conn,serverConn)
	}else{
		url = r.URL.String()
	}
/*
	req, err := http.NewRequest(r.Method,r.URL.String(),nil) 
	handle(err)
	res, err := p.client.Do(req)
	data := make([]byte, int(res.ContentLength))
	_, err = res.Body.(io.Reader).Read(data)
	fmt.Println(string(data))
	res.Body.Close()
	*/
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
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true, NextProtos: []string{"spdy/3"}}
	conn, err := tls.Dial("tcp", "127.0.0.1:8080", &config)
	if err != nil {
		fmt.Printf("client: dial: %s", err)
	}
	handler.conn = conn
	client, err := spdy.NewClientConn(conn)
	handle(err)	
	if client == nil{

	}
	handler.client = client

	http.ListenAndServe(":8181", handler)
	/*
	client, err := spdy.NewClientConn(conn)
	handle(err)	
	req, err := http.NewRequest("CONNECT", "https://baidu.com:443", nil)
	handle(err)
	res, err := client.Do(req)
	handle(err)
	data := make([]byte, int(res.ContentLength))
	_, err = res.Body.(io.Reader).Read(data)
	fmt.Println(string(data))
	res.Body.Close()
	*/
}