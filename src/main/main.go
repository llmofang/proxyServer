
package main

import (
	"myproxy"
	//"bytes"
	//"fmt"
	//"io"
	"log"
	"net/http"

	"github.com/amahi/spdy"
)


func main() {
	spdy.EnableDebug()
	handler := myproxy.NewHandler()
	http.HandleFunc("/",handler.Handle)


	//handler := spdy.ProxyConnHandlerFunc(handleProxy)
	//http.Handle("/", spdy.ProxyConnections(handler))
	err := spdy.ListenAndServeTLS(":8080", "../cert/serverTLS/server.pem", "../cert/serverTLS/server.key" , nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}	
}